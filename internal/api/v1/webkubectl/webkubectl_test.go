package webkubectl

import (
	"testing"

	"k8s.io/client-go/rest"
)

func TestToCmdConfigPreservesImpersonation(t *testing.T) {
	sess := &Session{
		Cluster: "eks",
		User:    "alice",
		config: &rest.Config{
			Host:        "https://example.invalid",
			BearerToken: "admin-token",
			Impersonate: rest.ImpersonationConfig{UserName: "alice", Groups: []string{"developers"}},
		},
	}

	cfg := toCmdConfig(sess)
	auth := cfg.AuthInfos["alice"]
	if auth == nil {
		t.Fatalf("expected auth info for alice")
	}
	if auth.Token != "admin-token" {
		t.Fatalf("expected bearer token to be preserved")
	}
	if auth.Impersonate != "alice" {
		t.Fatalf("expected impersonation username alice, got %q", auth.Impersonate)
	}
	if len(auth.ImpersonateGroups) != 1 || auth.ImpersonateGroups[0] != "developers" {
		t.Fatalf("expected impersonation groups to be preserved")
	}
}
