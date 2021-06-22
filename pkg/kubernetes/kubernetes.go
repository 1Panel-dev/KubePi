package kubernetes

import (
	"context"
	"errors"
	"fmt"
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"github.com/KubeOperator/ekko/pkg/certificate"
	v1 "k8s.io/api/authorization/v1"
	certv1 "k8s.io/api/certificates/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
	"time"
)

type Interface interface {
	Ping() error
	Version() (*version.Info, error)
	Client() (*kubernetes.Clientset, error)
	HasPermission(attributes v1.ResourceAttributes) (PermissionCheckResult, error)
	CreateCommonUser(commonName string) ([]byte, error)
	CreateDefaultClusterRoles() error
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

func (k *Kubernetes) CreateDefaultClusterRoles() error {
	defaultRoles := []rbacV1.ClusterRole{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ekko-admin",
				Annotations: map[string]string{
					"ekko-i18n":  "cluster_administrator",
					"created-by": "system",
					"created-at": time.Now().Format("2006-01-02 15:04:05"),
				},
				Labels: map[string]string{
					"manage": "ekko",
				},
			},
			Rules: []rbacV1.PolicyRule{
				{
					APIGroups: []string{"*"},
					Resources: []string{"*"},
					Verbs:     []string{"*"},
				},
				{
					NonResourceURLs: []string{"*"},
					Verbs:           []string{"*"},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ekko-viewer",
				Annotations: map[string]string{
					"ekko-i18n":  "cluster_viewer",
					"created-by": "system",
					"created-at": time.Now().Format("2006-01-02 15:04:05"),
				},
				Labels: map[string]string{
					"manage": "ekko",
				},
			},
			Rules: []rbacV1.PolicyRule{
				{
					APIGroups: []string{"*"},
					Resources: []string{"*"},
					Verbs:     []string{"list", "get"},
				},
				{
					NonResourceURLs: []string{"*"},
					Verbs:           []string{"list", "get"},
				},
			},
		},
	}
	client, err := k.Client()
	if err != nil {
		return err
	}
	for i := range defaultRoles {
		instance, err := client.RbacV1().ClusterRoles().Get(context.TODO(), defaultRoles[i].Name, metav1.GetOptions{})
		if err != nil {
			if !strings.Contains(strings.ToLower(err.Error()), "not found") {
				return err
			}
		}
		if instance == nil {
			_, err = client.RbacV1().ClusterRoles().Create(context.TODO(), &defaultRoles[i], metav1.CreateOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (k *Kubernetes) CreateCommonUser(commonName string) ([]byte, error) {
	// 生成用户证书申请
	cert, err := certificate.CreateClientCertificateRequest(commonName, k.PrivateKey)
	if err != nil {
		return nil, err
	}
	csr := certv1.CertificateSigningRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s-%d", commonName, "ekko", time.Now().Unix()),
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
	createResp, err := client.CertificatesV1().CertificateSigningRequests().Create(context.TODO(), &csr, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	// 审批证书
	createResp.Status.Conditions = append(createResp.Status.Conditions, certv1.CertificateSigningRequestCondition{
		Reason:         "Approved by Ekko",
		Type:           certv1.CertificateApproved,
		LastUpdateTime: metav1.Now(),
		Status:         "True",
	})

	updateResp, err := client.CertificatesV1().CertificateSigningRequests().UpdateApproval(context.TODO(), createResp.Name, createResp, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	var data []byte
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		getResp, err := client.CertificatesV1().CertificateSigningRequests().Get(context.TODO(), updateResp.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if getResp.Status.Certificate != nil {
			data = getResp.Status.Certificate
			break
		}
	}
	if data == nil {
		return nil, errors.New("csr approve time out")
	}

	return data, nil
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
