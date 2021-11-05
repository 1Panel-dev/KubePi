package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	v1Cluster "github.com/KubeOperator/kubepi/internal/model/v1/cluster"
	"github.com/KubeOperator/kubepi/pkg/certificate"
	"github.com/KubeOperator/kubepi/pkg/collectons"
	v1 "k8s.io/api/authorization/v1"
	certv1 "k8s.io/api/certificates/v1"
	certv1beta1 "k8s.io/api/certificates/v1beta1"
	rbacV1 "k8s.io/api/rbac/v1"
	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Interface interface {
	Ping() error
	Version() (*version.Info, error)
	VersionMinor() (int, error)
	Config() (*rest.Config, error)
	Client() (*kubernetes.Clientset, error)
	HasPermission(attributes v1.ResourceAttributes) (PermissionCheckResult, error)
	CreateCommonUser(commonName string) ([]byte, error)
	CreateDefaultClusterRoles() error
	GetUserNamespaceNames(username string, options ...interface{}) ([]string, error)
	CanVisitAllNamespace(username string) (bool, error)
	IsNamespacedResource(resourceName string) (bool, error)
	CleanManagedClusterRole() error
	CleanManagedClusterRoleBinding(username string) error
	CleanManagedRoleBinding(username string) error
	CleanAllRBACResource() error
	CreateOrUpdateClusterRoleBinding(clusterRoleName string, username string, builtIn bool) error
	CreateOrUpdateRolebinding(namespace string, clusterRoleName string, username string, builtIn bool) error
	CreateAppMarketCRD() error
}

type Kubernetes struct {
	*v1Cluster.Cluster
}

func (k *Kubernetes) VersionMinor() (int, error) {
	v, err := k.Version()
	if err != nil {
		return 0, err
	}
	reg := regexp.MustCompile("[^0-9]")
	minor, err := strconv.Atoi(reg.ReplaceAllString(v.Minor, ""))
	return minor, err
}

func (k *Kubernetes) CreateOrUpdateClusterRoleBinding(clusterRoleName string, username string, builtIn bool) error {
	client, err := k.Client()
	if err != nil {
		return err
	}
	name := fmt.Sprintf("%s:%s:%s", username, clusterRoleName, k.UUID)
	labels := map[string]string{
		LabelManageKey: "kubepi",
		LabelClusterId: k.UUID,
		LabelUsername:  username,
	}
	annotations := map[string]string{
		"built-in":   strconv.FormatBool(builtIn),
		"created-at": time.Now().Format("2006-01-02 15:04:05"),
	}
	item := rbacV1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
		},
		Subjects: []rbacV1.Subject{
			{
				Kind: "User",
				Name: username,
			},
		},
		RoleRef: rbacV1.RoleRef{
			Kind: "ClusterRole",
			Name: clusterRoleName,
		},
	}
	baseItem, err := client.RbacV1().ClusterRoleBindings().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}
	if baseItem != nil && baseItem.Name != "" {
		_, err := client.RbacV1().ClusterRoleBindings().Create(context.TODO(), &item, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		_, err := client.RbacV1().ClusterRoleBindings().Update(context.TODO(), &item, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Kubernetes) CreateOrUpdateRolebinding(namespace string, clusterRoleName string, username string, builtIn bool) error {
	client, err := k.Client()
	if err != nil {
		return err
	}
	labels := map[string]string{
		LabelManageKey: "kubepi",
		LabelClusterId: k.UUID,
		LabelUsername:  username,
	}
	annotations := map[string]string{
		"built-in":   strconv.FormatBool(builtIn),
		"created-at": time.Now().Format("2006-01-02 15:04:05"),
	}
	name := fmt.Sprintf("%s:%s:%s:%s", namespace, username, clusterRoleName, k.UUID)
	item := rbacV1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
			Namespace:   namespace,
		},
		Subjects: []rbacV1.Subject{
			{
				Kind: "User",
				Name: username,
			},
		},
		RoleRef: rbacV1.RoleRef{
			Kind: "ClusterRole",
			Name: clusterRoleName,
		},
	}
	baseItem, err := client.RbacV1().RoleBindings(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}
	if baseItem != nil && baseItem.Name != "" {
		_, err := client.RbacV1().RoleBindings(namespace).Create(context.TODO(), &item, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		_, err := client.RbacV1().RoleBindings(namespace).Update(context.TODO(), &item, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Kubernetes) CleanManagedClusterRole() error {
	client, err := k.Client()
	if err != nil {
		return err
	}

	labels := []string{
		fmt.Sprintf("%s=%s", LabelManageKey, "kubepi"),
		fmt.Sprintf("%s=%s", LabelClusterId, k.UUID),
	}
	return client.RbacV1().ClusterRoles().DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: strings.Join(labels, ","),
	})
}

func (k *Kubernetes) CleanManagedClusterRoleBinding(username string) error {
	client, err := k.Client()
	if err != nil {
		return err
	}
	labels := []string{
		fmt.Sprintf("%s=%s", LabelManageKey, "kubepi"),
		fmt.Sprintf("%s=%s", LabelClusterId, k.UUID),
	}
	if username != "" {
		labels = append(labels, fmt.Sprintf("%s=%s", LabelUsername, username))
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
	labels := []string{
		fmt.Sprintf("%s=%s", LabelManageKey, "kubepi"),
		fmt.Sprintf("%s=%s", LabelClusterId, k.UUID),
	}
	if username != "" {
		labels = append(labels, fmt.Sprintf("%s=%s", LabelUsername, username))
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
	labels := []string{
		fmt.Sprintf("%s=%s", LabelManageKey, "kubepi"),
		fmt.Sprintf("%s=%s", LabelClusterId, k.UUID),
		fmt.Sprintf("%s=%s", LabelUsername, username),
	}
	clusterrolebindings, err := client.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{
		LabelSelector: strings.Join(labels, ","),
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
func (k *Kubernetes) GetUserNamespaceNames(username string, options ...interface{}) ([]string, error) {
	client, err := k.Client()
	if err != nil {
		return nil, err
	}
	all := false
	if len(options) > 0 && options[0].(bool) {
		all = true
	} else {
		all, err = k.CanVisitAllNamespace(username)
		if err != nil {
			return nil, err
		}
	}

	namespaceSet := collectons.NewStringSet()
	if all {
		ns, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for i := range ns.Items {
			if ns.Items[i].Status.Phase == "Active" {
				namespaceSet.Add(ns.Items[i].Name)
			}
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

	result := namespaceSet.ToSlice()
	sort.Strings(result)
	return result, nil
}

func (k *Kubernetes) IsNamespacedResource(resourceName string) (bool, error) {
	if resourceName == "events" {
		return false, nil
	}
	client, err := k.Client()
	if err != nil {
		return false, err
	}
	apiList, err := client.ServerPreferredNamespacedResources()
	if err != nil && len(apiList) == 0 {
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

func NewKubernetes(cluster *v1Cluster.Cluster) Interface {
	return &Kubernetes{Cluster: cluster}
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
		// 如果已存在，尝试更新
		if instance == nil || instance.Name == "" {
			_, err = client.RbacV1().ClusterRoles().Create(context.TODO(), &initClusterRoles[i], metav1.CreateOptions{})
			if err != nil {
				return err
			}
		} else {
			_, err = client.RbacV1().ClusterRoles().Update(context.TODO(), &initClusterRoles[i], metav1.UpdateOptions{})
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

	client, err := k.Client()
	if err != nil {
		return nil, err
	}
	v, err := client.ServerVersion()
	if err != nil {
		return nil, err
	}
	reg := regexp.MustCompile("[^0-9]")
	minor, err := strconv.Atoi(reg.ReplaceAllString(v.Minor, ""))
	if err != nil {
		return nil, err
	}
	var data []byte
	if minor > 18 {
		csr := certv1.CertificateSigningRequest{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-%s-%d", commonName, "kubepi", time.Now().Unix()),
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
		createResp, err := client.CertificatesV1().CertificateSigningRequests().Create(context.TODO(), &csr, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
		// 审批证书
		createResp.Status.Conditions = append(createResp.Status.Conditions, certv1.CertificateSigningRequestCondition{
			Reason:         "Approved by KubePi",
			Type:           certv1.CertificateApproved,
			LastUpdateTime: metav1.Now(),
			Status:         "True",
		})

		updateResp, err := client.CertificatesV1().CertificateSigningRequests().UpdateApproval(context.TODO(), createResp.Name, createResp, metav1.UpdateOptions{})
		if err != nil {
			return nil, err
		}

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
	} else {
		name := "kubernetes.io/kube-apiserver-client"
		csr := certv1beta1.CertificateSigningRequest{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-%s-%d", commonName, "kubepi", time.Now().Unix()),
			},
			Spec: certv1beta1.CertificateSigningRequestSpec{
				SignerName: &name,
				Request:    cert,
				Groups: []string{
					"system:authenticated",
				},
				Usages: []certv1beta1.KeyUsage{
					"client auth",
				},
			},
		}
		createResp, err := client.CertificatesV1beta1().CertificateSigningRequests().Create(context.TODO(), &csr, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
		// 审批证书
		createResp.Status.Conditions = append(createResp.Status.Conditions, certv1beta1.CertificateSigningRequestCondition{
			Reason:         "Approved by KubePi",
			Type:           certv1beta1.CertificateApproved,
			LastUpdateTime: metav1.Now(),
			Status:         "True",
		})

		updateResp, err := client.CertificatesV1beta1().CertificateSigningRequests().UpdateApproval(context.TODO(), createResp, metav1.UpdateOptions{})
		if err != nil {
			return nil, err
		}

		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			getResp, err := client.CertificatesV1beta1().CertificateSigningRequests().Get(context.TODO(), updateResp.Name, metav1.GetOptions{})
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
func (k *Kubernetes) Config() (*rest.Config, error) {
	if k.Spec.Local {
		return rest.InClusterConfig()
	}
	if k.Spec.Connect.Direction == "forward" {
		kubeConf := &rest.Config{
			Host: k.Spec.Connect.Forward.ApiServer,
		}
		if len(k.CaCertificate.CertData) > 0 {
			kubeConf.CAData = k.CaCertificate.CertData
		} else {
			kubeConf.Insecure = true
		}
		switch strings.ToLower(k.Spec.Authentication.Mode) {
		case "bearer":
			kubeConf.BearerToken = k.Spec.Authentication.BearerToken
		case "certificate":
			kubeConf.TLSClientConfig.CertData = k.Spec.Authentication.Certificate.CertData
			kubeConf.TLSClientConfig.KeyData = k.Spec.Authentication.Certificate.KeyData
		case "configfile":
			cfg, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*clientcmdapi.Config, error) {
				return clientcmd.Load(k.Spec.Authentication.ConfigFileContent)
			})
			if err != nil {
				return nil, err
			}
			kubeConf = cfg
		}
		return kubeConf, nil
	}
	return nil, nil
}

func (k *Kubernetes) Client() (*kubernetes.Clientset, error) {
	cfg, err := k.Config()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(cfg)
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

func (k *Kubernetes) CreateAppMarketCRD() error {
	cfg, err := k.Config()
	if err != nil {
		return err
	}
	client, err := apiextension.NewForConfig(cfg)
	if err != nil {
		return err
	}
	client.ApiextensionsV1().CustomResourceDefinitions().Delete(context.TODO(), "appmarkets.kubepi.org", metav1.DeleteOptions{})
	crd := &apiextensionv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: "appmarkets.kubepi.org"},
		Spec: apiextensionv1.CustomResourceDefinitionSpec{
			Group: "kubepi.org",
			Scope: apiextensionv1.ClusterScoped,
			Names: apiextensionv1.CustomResourceDefinitionNames{
				Plural:     "appmarkets",
				Kind:       "Appmarket",
				Singular:   "appmarket",
				ShortNames: []string{"am"},
			},
			Versions: []apiextensionv1.CustomResourceDefinitionVersion{{
				Name:    "v1",
				Served:  true,
				Storage: true,
				Schema: &apiextensionv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionv1.JSONSchemaProps{
						Type: "object",
						Properties: map[string]apiextensionv1.JSONSchemaProps{
							"name": {
								Type: "string",
							},
							"url": {
								Type: "string",
							},
							"auth": {
								Type: "string",
							},
						},
					},
				},
			}},
		},
	}
	_, err = client.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), crd, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
