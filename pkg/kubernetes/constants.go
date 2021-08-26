package kubernetes

import (
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const (
	LabelManageKey   = "kubeoperator.io/manage"
	LabelRoleTypeKey = "kubeoperator.io/role-type"
	LabelClusterId   = "kubeoperator.io/cluster-id"
	LabelUsername    = "kubeoperator.io/username"

	RoleTypeCluster   = "cluster"
	RoleTypeNamespace = "namespace"
)

var initClusterRoles = []rbacV1.ClusterRole{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "admin-cluster",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeCluster,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-cluster",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-namespaces",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-namespaces",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-nodes",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"nodes"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-events",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"events"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-cluster-rbac",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"clusterroles", "clusterrolebindings"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-cluster-rbac",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"clusterroles", "clusterrolebindings"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-cluster-storage",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"persistentvolumes"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"storageclasses"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-cluster-storage",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"persistentvolumes"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"storageclasses"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-service-discovery",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"services", "endpoints"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingresses", "networkpolicies"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-service-discovery",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeNamespace},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"services", "endpoints"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingresses", "networkpolicies"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-config",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeNamespace},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps", "secrets", "resourcequotas", "limitranges"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"autoscaling"},
				Resources: []string{"horizontalpodautoscalers"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-config",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeNamespace},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps", "secrets", "resourcequotas", "limitranges"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"autoscaling"},
				Resources: []string{"horizontalpodautoscalers"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-storage",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"persistentvolumeclaims"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-storage",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"persistentvolumeclaims"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-workload",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "containers"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"deployments", "daemonsets", "replicasets", "statefulsets"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"batch"},
				Resources: []string{"jobs", "cronjobs"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-workload",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "containers"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"deployments", "daemonsets", "replicasets", "statefulsets"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"batch"},
				Resources: []string{"jobs", "cronjobs"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-rbac",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"roles", "rolebindings"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"serviceaccounts"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"policy"},
				Resources: []string{"podsecuritypolicies"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-rbac",
			Annotations: map[string]string{
				"builtin":    "true",
				"created-at": time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "ekko",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"roles", "rolebindings"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"serviceaccounts"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"policy"},
				Resources: []string{"podsecuritypolicies"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
}
