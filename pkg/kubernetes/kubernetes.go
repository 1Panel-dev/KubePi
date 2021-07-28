package kubernetes

import (
	"context"
	"errors"
	"fmt"
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"github.com/KubeOperator/ekko/pkg/certificate"
	"github.com/KubeOperator/ekko/pkg/collectons"
	v1 "k8s.io/api/authorization/v1"
	certv1 "k8s.io/api/certificates/v1"
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
	GetUserNamespaceNames(username string) ([]string, error)
	CanVisitAllNamespace(username string) (bool, error)
	IsNamespacedResource(resourceName string) (bool, error)
	CleanManagedClusterRole() error
	CleanManagedClusterRoleBinding(username string) error
	CleanManagedRoleBinding(username string) error
	CleanAllRBACResource() error
}

type Kubernetes struct {
	*v1Cluster.Cluster
}

func (k *Kubernetes) CleanManagedClusterRole() error {
	client, err := k.Client()
	if err != nil {
		return err
	}
	return client.RbacV1().ClusterRoles().DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", LabelManageKey, "ekko"),
	})
}

func (k *Kubernetes) CleanManagedClusterRoleBinding(username string) error {
	client, err := k.Client()
	if err != nil {
		return err
	}
	labels := []string{fmt.Sprintf("%s=%s", LabelManageKey, "ekko")}
	if username != "" {
		labels = append(labels, fmt.Sprintf("%s=%s", "user-name", username))
	}
	return client.RbacV1().ClusterRoleBindings().DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: strings.Join(labels, ","),
	})
}

func (k *Kubernetes) CleanManagedRoleBinding(username string) error {
	client, err := k.Client()
	if err != nil {
		return err
	}
	nss, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	labels := []string{fmt.Sprintf("%s=%s", LabelManageKey, "ekko")}
	if username != "" {
		labels = append(labels, fmt.Sprintf("%s=%s", "user-name", username))
	}
	for i := range nss.Items {
		if err := client.RbacV1().RoleBindings(nss.Items[i].Name).DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{
			LabelSelector: strings.Join(labels, ","),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (k *Kubernetes) CleanAllRBACResource() error {
	if err := k.CleanManagedClusterRole(); err != nil {
		return err
	}
	if err := k.CleanManagedClusterRoleBinding(""); err != nil {
		return err
	}

	if err := k.CleanManagedRoleBinding(""); err != nil {
		return err
	}
	return nil
}

func (k *Kubernetes) CanVisitAllNamespace(username string) (bool, error) {
	client, err := k.Client()
	if err != nil {
		return false, err
	}
	roleSet := collectons.NewStringSet()
	clusterrolebindings, err := client.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("user-name=%s", username),
	})
	if err != nil {
		return false, err
	}
	for i := range clusterrolebindings.Items {
		roleSet.Add(clusterrolebindings.Items[i].RoleRef.Name)
	}
	for _, roleName := range roleSet.ToSlice() {
		role, err := client.RbacV1().ClusterRoles().Get(context.TODO(), roleName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for i := range role.Rules {
			if collectons.IndexOfStringSlice(role.Rules[i].APIGroups, "*") != -1 && collectons.IndexOfStringSlice(role.Rules[i].Resources, "*") != -1 {
				return true, nil
			}
		}
	}
	return false, nil
}
func (k *Kubernetes) GetUserNamespaceNames(username string) ([]string, error) {
	client, err := k.Client()
	if err != nil {
		return nil, err
	}
	all, err := k.CanVisitAllNamespace(username)
	if err != nil {
		return nil, err
	}
	namespaceSet := collectons.NewStringSet()
	if all {
		ns, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for i := range ns.Items {
			namespaceSet.Add(ns.Items[i].Name)
		}
	} else {
		rbs, err := client.RbacV1().RoleBindings("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for i := range rbs.Items {
			for j := range rbs.Items[i].Subjects {
				if rbs.Items[i].Subjects[j].Kind == "User" && rbs.Items[i].Subjects[j].Name == username {
					namespaceSet.Add(rbs.Items[i].Namespace)
				}
			}
		}
	}

	return namespaceSet.ToSlice(), nil
}

func (k *Kubernetes) IsNamespacedResource(resourceName string) (bool, error) {
	client, err := k.Client()
	if err != nil {
		return false, err
	}
	apiList, err := client.ServerPreferredNamespacedResources()
	if err != nil {
		return false, err
	}
	for i := range apiList {
		for j := range apiList[i].APIResources {
			if apiList[i].APIResources[j].Name == resourceName {
				return true, nil
			}
		}
	}
	return false, nil
}

type PermissionCheckResult struct {
	Resource v1.ResourceAttributes
	Allowed  bool
}

func NewKubernetes(cluster v1Cluster.Cluster) Interface {
	return &Kubernetes{Cluster: &cluster}
}

func (k *Kubernetes) CreateDefaultClusterRoles() error {

	client, err := k.Client()
	if err != nil {
		return err
	}
	for i := range initClusterRoles {
		instance, err := client.RbacV1().ClusterRoles().Get(context.TODO(), initClusterRoles[i].Name, metav1.GetOptions{})
		if err != nil {
			if !strings.Contains(strings.ToLower(err.Error()), "not found") {
				return err
			}
		}
		if instance == nil || instance.Name == "" {
			_, err = client.RbacV1().ClusterRoles().Create(context.TODO(), &initClusterRoles[i], metav1.CreateOptions{})
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
