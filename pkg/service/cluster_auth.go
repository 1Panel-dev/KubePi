package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"github.com/mitchellh/go-homedir"
	v1 "k8s.io/api/certificates/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"path"
)

func main() {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	subj := pkix.Name{
		CommonName: "user1",
	}
	rawSubj := subj.ToRDNSequence()
	asn1Subj, _ := asn1.Marshal(rawSubj)
	csr := &x509.CertificateRequest{
		RawSubject:         asn1Subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}
	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, csr, key)
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})

	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", path.Join(home, ".kube/config"))
	if err != nil {
		log.Fatal(err)
	}
	clients, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	obj := &v1.CertificateSigningRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: "user1",
		},
		Spec: v1.CertificateSigningRequestSpec{
			SignerName: "kubernetes.io/kube-apiserver-client",
			Request:    pemBytes,
			Groups: []string{
				"system:authenticated",
			},
			Usages: []v1.KeyUsage{
				"client auth",
			},
		},
	}

	_, err = clients.CertificatesV1().CertificateSigningRequests().Create(context.TODO(), obj, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}

	obj.Status.Conditions = append(obj.Status.Conditions, v1.CertificateSigningRequestCondition{
		Reason:         "aaa",
		Type:           v1.CertificateApproved,
		LastUpdateTime: metav1.Now(),
		Status:         "True",
	})

	_, err = clients.CertificatesV1().CertificateSigningRequests().UpdateApproval(context.TODO(), obj.Name, obj, metav1.UpdateOptions{})
	if err != nil {
		log.Fatal(err)
	}

}
