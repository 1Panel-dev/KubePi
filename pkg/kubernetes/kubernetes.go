package kubernetes

import (
	"context"
	"errors"
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"github.com/KubeOperator/ekko/pkg/certificate"
	v1 "k8s.io/api/authorization/v1"
	certv1 "k8s.io/api/certificates/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
)

type Interface interface {
	Ping() error
	Version() (*version.Info, error)
	Client() (*kubernetes.Clientset, error)
	HasPermission(attributes v1.ResourceAttributes) (PermissionCheckResult, error)
	CreateUser(commonName string) ([]byte, error)
}

type Kubernetes struct {
	*v1Cluster.Cluster
}

type PermissionCheckResult struct {
	Resource v1.ResourceAttributes
	Allowed  bool
}

func NewKubernetes(cluster v1Cluster.Cluster) Interface {
	return &Kubernetes{Cluster: &cluster}
}

func (k *Kubernetes) CreateUser(commonName string) ([]byte, error) {
	// 生成用户证书申请
	cert, err := certificate.CreateClientCertificateRequest(commonName, &k.PrivateKey)
	if err != nil {
		return nil, err
	}
	csr := certv1.CertificateSigningRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: "user1",
		},
		Spec: certv1.CertificateSigningRequestSpec{
			SignerName: "kubernetes.io/kube-apiserver-client",
			Request:    cert,
			Groups: []string{
				"system:authenticated",
			},
			Usages: []certv1.KeyUsage{
				"client auth",
			},
		},
	}
	client, err := k.Client()
	if err != nil {
		return nil, err
	}
	resp, err := client.CertificatesV1().CertificateSigningRequests().Create(context.TODO(), &csr, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	csr = *resp
	// 审批证书
	csr.Status.Conditions = append(csr.Status.Conditions, certv1.CertificateSigningRequestCondition{
		Reason:         "approve by ekko",
		Type:           certv1.CertificateApproved,
		LastUpdateTime: metav1.Now(),
		Status:         "True",
	})

	resp, err = client.CertificatesV1().CertificateSigningRequests().UpdateApproval(context.TODO(), csr.Name, &csr, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	// 读取证书
	return resp.Status.Certificate, nil
}

func (k *Kubernetes) HasPermission(attributes v1.ResourceAttributes) (PermissionCheckResult, error) {
	client, err := k.Client()
	if err != nil {
		return PermissionCheckResult{}, err
	}
	resp, err := client.AuthorizationV1().SelfSubjectAccessReviews().Create(context.TODO(), &v1.SelfSubjectAccessReview{
		Spec: v1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &attributes,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return PermissionCheckResult{}, err
	}
	return PermissionCheckResult{
		Resource: attributes,
		Allowed:  resp.Status.Allowed,
	}, nil

}

func (k *Kubernetes) Client() (*kubernetes.Clientset, error) {
	if k.Spec.Connect.Direction == "forward" {
		kubeConf := &rest.Config{
			Host: k.Spec.Connect.Forward.ApiServer,
		}
		if len(k.CaCertificate.CertData) > 0 {
			kubeConf.CAData = k.CaCertificate.CertData
		} else {
			kubeConf.Insecure = true
		}
		switch k.Spec.Authentication.Mode {
		case strings.ToLower("bearer"):
			kubeConf.BearerToken = k.Spec.Authentication.BearerToken
		case strings.ToLower("certificate"):
			kubeConf.TLSClientConfig.CertData = k.Spec.Authentication.Certificate.CertData
			kubeConf.TLSClientConfig.KeyData = k.Spec.Authentication.Certificate.KeyData
		}
		return kubernetes.NewForConfig(kubeConf)

	}
	return nil, errors.New("")
}

func (k *Kubernetes) Version() (*version.Info, error) {
	client, err := k.Client()
	if err != nil {
		return nil, err
	}
	return client.ServerVersion()

}

func (k *Kubernetes) Ping() error {
	client, err := k.Client()
	if err != nil {
		return err
	}
	_, err = client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	client.AuthorizationV1().SelfSubjectAccessReviews()
	if err != nil {
		return err
	}
	return nil
}
