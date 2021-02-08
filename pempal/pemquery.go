package pempal

import (
	"github.com/pempal/templates"
	"log"
	"strings"

	"crypto/x509"
	"encoding/pem"
	"github.com/pempal/pemio"
	"gopkg.in/yaml.v3"
)

type PEMQuery struct {
	Query         []string
	CaseSensitive bool

	Types    []string
	Verbose  bool
	Password string
}

type QueryResult struct {
	Block      *pem.Block
	QueryMatch []string
	FilePath   string
}

func (pq PEMQuery) QueryPaths(ps []string, recursive bool) ([]*QueryResult, error) {
	var r []*QueryResult
	for _, p := range ps {
		qrs, err := pq.QueryPath(p, recursive)
		if err != nil {
			return nil, err
		}
		r = append(r, qrs...)
	}
	return r, nil
}

func (pq PEMQuery) QueryPath(p string, recursive bool) ([]*QueryResult, error) {
	ps := &pemio.PEMScanner{
		FilePath:  p,
		Recursive: recursive,
		Verbose:   pq.Verbose,
	}
	pfs, err := ps.ScanPath()
	if err != nil {
		return nil, err
	}
	var r []*QueryResult
	for _, pf := range pfs {
		qrs, err := pq.QueryPemFile(pf)
		if err != nil {
			return nil, err
		}
		if len(qrs) > 0 {
			r = append(r, qrs...)
		}
	}
	return r, nil
}

func (pq PEMQuery) QueryPemFile(pf *pemio.PEMFile) ([]*QueryResult, error) {
	bls := pq.filterBlocksByType(pf.Blocks)
	if len(bls) == 0 {
		return nil, nil
	}
	var r []*QueryResult
	for _, bl := range bls {
		if pq.Password != "" && x509.IsEncryptedPEMBlock(bl) {
			b, err := x509.DecryptPEMBlock(bl, []byte(pq.Password))
			if err != nil {
				if pq.Verbose {
					log.Printf("failed to decrypt %s  %v\n", pf.FilePath, err)
				}
			}
			bl.Bytes = b
		}
		qr, err := pq.filterBlockByQuery(bl, pf.FilePath)
		if err != nil {
			return nil, err
		}
		if qr != nil {
			r = append(r, qr)
		}
	}
	return r, nil
}

func (pq PEMQuery) filterBlocksByType(pbs []*pem.Block) []*pem.Block {
	if len(pq.Types) == 0 {
		return pbs
	}
	var pbss []*pem.Block
	for _, pb := range pbs {
		if !containsType(pb.Type, pq.Types) {
			continue
		}
		pbss = append(pbss, pb)
	}
	return pbss
}

func (pq PEMQuery) filterBlockByQuery(b *pem.Block, p string) (*QueryResult, error) {
	fp := &QueryResult{Block: b, FilePath: p}
	if len(pq.Query) == 0 {
		return fp, nil
	}
	t, err := templates.NewTemplate(b)
	if err != nil {
		return nil, err
	}
	by, err := yaml.Marshal(t)
	if err != nil {
		return nil, err
	}

	for _, q := range pq.Query {
		l := findLine(string(by), q, pq.CaseSensitive)
		if l == "" {
			continue
		}
		fp.QueryMatch = append(fp.QueryMatch, l)
	}

	if len(fp.QueryMatch) > 0 {
		return fp, nil
	}
	return nil, nil
}

func findLine(s string, q string, cs bool) string {
	if !cs {
		s = strings.ToLower(s)
		q = strings.ToLower(q)
	}
	ss := strings.Split(s, "\n")
	for _, l := range ss {
		if strings.Contains(l, q) {
			return strings.TrimSpace(l)
		}
	}
	return ""
}

func containsType(t string, ts []string) bool {
	for _, s := range ts {
		if strings.EqualFold(t, s) {
			return true
		}
	}
	return false
}
