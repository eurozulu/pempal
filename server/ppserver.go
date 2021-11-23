package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type PPServer struct {
	HostCertPath string
	KeyPath      string

	Port int
}

func (s PPServer) Start(ctx context.Context) <-chan error {
	ch := make(chan error)
	go func() {
		r := mux.NewRouter()
		ph := &ppHandler{}
		r.HandleFunc("/certificate", ph.serveGetCertificate).Methods("GET")
		r.HandleFunc("/certificate/request", ph.serveRequestCertificate).Methods("POST")
		r.HandleFunc("/keytracker", ph.serveGetKey).Methods("GET")
		r.HandleFunc("/keytracker/request", ph.serveRequestKey).Methods("POST")

		caCert, _ := ioutil.ReadFile("ca.crt")
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			ClientCAs:  caCertPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
		tlsConfig.BuildNameToCertificate()

		server := &http.Server{
			Addr:      ":9443",
			TLSConfig: tlsConfig,
		}

		server.ListenAndServeTLS("server.crt", "server.key")

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", s.Port),
			Handler: r,
		}
		chSrv := make(chan error)
		go func() {
			defer close(chSrv)
			err := srv.ListenAndServeTLS(s.HostCertPath, s.KeyPath)
			if err != nil {
				chSrv <- err
			}
		}()

		for {
			select {
			case <-ctx.Done():
				if err := srv.Close(); err != nil {
					ch <- err
				}
			case err, ok := <-chSrv:
				if !ok {
					return
				}
				ch <- err
				return
			}
		}
	}()
	return ch
}
