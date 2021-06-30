package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
)

func GeneratePrivateKey() ([]byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	return derStream, nil

}

func CreateClientCertificateRequest(userName string, key []byte, org ...string) ([]byte, error) {
	privateKey, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}
	subj := pkix.Name{
		CommonName: userName,
	}
	for i := range org {
		subj.Organization = append(subj.Organization, org[i])
	}
	rawSubj := subj.ToRDNSequence()
	asn1Subj, err := asn1.Marshal(rawSubj)
	if err != nil {
		return nil, err
	}
	csr := &x509.CertificateRequest{
		RawSubject:         asn1Subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csr, privateKey)
	if err != nil {
		return nil, err
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
	return pemBytes, nil
}

func ParseX509Certificate(certPem []byte) (*x509.Certificate, error) {
	b, _ := pem.Decode(certPem)
	if b == nil {
		return nil, errors.New("can not decode pem block")
	}
	c, err := x509.ParseCertificate(b.Bytes)
	if err != nil {
		return nil, err
	}
	return c, nil
}
