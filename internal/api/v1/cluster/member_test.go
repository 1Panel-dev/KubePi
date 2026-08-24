package cluster

import (
	"errors"
	"testing"

	v1Cluster "github.com/1Panel-dev/KubePi/internal/model/v1/cluster"
)

func TestCanContinueWithoutMemberCertificateForCSRTimeout(t *testing.T) {
	if !canContinueWithoutMemberCertificate(errors.New("csr approve time out ,30 sec")) {
		t.Fatalf("expected CSR timeout to allow impersonation fallback")
	}
}

func TestCanContinueWithoutMemberCertificateRejectsOtherErrors(t *testing.T) {
	if canContinueWithoutMemberCertificate(errors.New("forbidden")) {
		t.Fatalf("expected non-CSR errors to remain fatal")
	}
}

func TestShouldUseImpersonationOnlyForTokenKubeconfig(t *testing.T) {
	c := &v1Cluster.Cluster{}
	c.Spec.Authentication.Mode = "configFile"
	c.Spec.Authentication.ConfigFileContent = []byte(`apiVersion: v1
kind: Config
clusters:
- name: c
  cluster:
    server: https://example.invalid
contexts:
- name: c
  context:
    cluster: c
    user: u
current-context: c
users:
- name: u
  user:
    token: token-value
`)

	if !shouldUseImpersonationOnly(c) {
		t.Fatalf("expected token kubeconfig to skip CSR and use impersonation")
	}
}
func TestShouldUseImpersonationOnlyForExecKubeconfig(t *testing.T) {
	c := &v1Cluster.Cluster{}
	c.Spec.Authentication.Mode = "configFile"
	c.Spec.Authentication.ConfigFileContent = []byte(`apiVersion: v1
kind: Config
clusters:
- name: c
  cluster:
    server: https://example.invalid
contexts:
- name: c
  context:
    cluster: c
    user: u
current-context: c
users:
- name: u
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: aws
`)

	if !shouldUseImpersonationOnly(c) {
		t.Fatalf("expected exec kubeconfig to skip CSR and use impersonation")
	}
}

func TestShouldUseImpersonationOnlyRejectsCertificateKubeconfig(t *testing.T) {
	c := &v1Cluster.Cluster{}
	c.Spec.Authentication.Mode = "configFile"
	c.Spec.Authentication.ConfigFileContent = []byte(`apiVersion: v1
kind: Config
clusters:
- name: c
  cluster:
    server: https://example.invalid
contexts:
- name: c
  context:
    cluster: c
    user: u
current-context: c
users:
- name: u
  user:
    client-certificate-data: cert
    client-key-data: key
`)

	if shouldUseImpersonationOnly(c) {
		t.Fatalf("expected certificate kubeconfig to keep CSR certificate flow")
	}
}
