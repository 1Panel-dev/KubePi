package clusteraccess

import (
	"testing"

	v1Cluster "github.com/1Panel-dev/KubePi/internal/model/v1/cluster"
	"k8s.io/client-go/rest"
)

func TestApplyUserAccessUsesImpersonationWhenMemberCertificateMissing(t *testing.T) {
	cfg := &rest.Config{BearerToken: "admin-token"}
	cluster := &v1Cluster.Cluster{PrivateKey: []byte("private-key")}
	binding := &v1Cluster.Binding{UserRef: "alice"}

	ApplyUserAccessConfig(cfg, cluster, binding)

	if cfg.BearerToken != "admin-token" {
		t.Fatalf("expected imported cluster bearer token to remain for impersonation, got %q", cfg.BearerToken)
	}
	if cfg.Impersonate.UserName != "alice" {
		t.Fatalf("expected impersonation username alice, got %q", cfg.Impersonate.UserName)
	}
	if len(cfg.CertData) != 0 || len(cfg.KeyData) != 0 {
		t.Fatalf("expected no client certificate when binding has none")
	}
}

func TestApplyUserAccessUsesMemberCertificateWhenAvailable(t *testing.T) {
	cfg := &rest.Config{BearerToken: "admin-token"}
	cluster := &v1Cluster.Cluster{PrivateKey: []byte("private-key")}
	binding := &v1Cluster.Binding{UserRef: "alice", Certificate: []byte("cert")}

	ApplyUserAccessConfig(cfg, cluster, binding)

	if cfg.BearerToken != "" {
		t.Fatalf("expected bearer token cleared for certificate auth, got %q", cfg.BearerToken)
	}
	if cfg.Impersonate.UserName != "" {
		t.Fatalf("expected no impersonation when member certificate exists, got %q", cfg.Impersonate.UserName)
	}
	if string(cfg.CertData) != "cert" {
		t.Fatalf("expected member certificate to be used")
	}
	if len(cfg.KeyData) == 0 {
		t.Fatalf("expected cluster private key to be used with member certificate")
	}
}
