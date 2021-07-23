package kubernetes

import (
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const (
	LabelManageKey   = "kubeoperator.io/manage"
	LabelRoleTypeKey = "kubeoperator.io/role-type"

	RoleTypeCluster   = "cluster"
	RoleTypeNamespace = "namespace"
)

var initClusterRoles = []rbacV1.ClusterRole{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "Admin Cluster",
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
			Name: "View Cluster",
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
			Name: "Manage Namespaces",
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
			Name: "View Namespaces",
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
			Name: "View Nodes",
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
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "View Nodes",
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
			Name: "View Events",
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
			Name: "Manage Network",
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
			Name: "View Network",
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
			Name: "Manage Config",
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
				Resources: []string{"configmaps,secrets,resourcequotas,limitranges"},
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
			Name: "View Config",
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
				Resources: []string{"configmaps,secrets,resourcequotas,limitranges"},
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
			Name: "Manage Storage",
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
				Resources: []string{"persistentvolumes", "persistentvolumeclaims"},
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
			Name: "View Storage",
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
				Resources: []string{"persistentvolumes", "persistentvolumeclaims"},
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
			Name: "Manage Workload",
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
			Name: "View Workload",
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
			Name: "Mange RBAC",
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
				Resources: []string{"clusterroles", "clusterrolebindings", "roles", "rolebindings"},
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
			Name: "Mange RBAC",
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
				Resources: []string{"clusterroles", "clusterrolebindings", "roles", "rolebindings"},
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
