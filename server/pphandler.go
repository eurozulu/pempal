package server

import (
	"net/http"
)

type ppHandler struct {
}

func (ph ppHandler) serveGetCertificate(writer http.ResponseWriter, request *http.Request) {

}
func (ph ppHandler) serveRequestCertificate(writer http.ResponseWriter, request *http.Request) {

}
func (ph ppHandler) serveGetKey(writer http.ResponseWriter, request *http.Request) {

}
func (ph ppHandler) serveRequestKey(writer http.ResponseWriter, request *http.Request) {

}

func (ph ppHandler) handelError(err error, w http.ResponseWriter) bool {
	if err == nil {
		return true
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return false
}
