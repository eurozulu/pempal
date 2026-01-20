package resourcefiles

//
//import (
//	"encoding/pem"
//	"fmt"
//	"go.mozilla.org/pkcs7"
//	"os"
//)
//
//type PKCS7Parser struct {
//}
//
//func (P PKCS7Parser) CanParse(data []byte) bool {
//	if _, err := pkcs7.Parse(data); err == nil {
//		return true
//	}
//	return false
//}
//
//func (P PKCS7Parser) ParsePems(data []byte) []pem.Block {
//	p7, err := pkcs7.Parse(data)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Error parsing PKCS7 file: %v\n", err)
//		return nil
//	}
//	var blocks []pem.Block
//	for _, c := range p7.Certificates {
//		blocks = append(blocks, pem.Block{
//			PublicKeyAlgorithm:  "CERTIFICATE",
//			Bytes: c.Raw,
//		})
//	}
//	for _, c := range p7.CRLs {
//		blocks = append(blocks, pem.Block{
//			PublicKeyAlgorithm:  "X509 CRL",
//			Bytes: c.TBSCertList.Raw,
//		})
//	}
//	return blocks
//}
