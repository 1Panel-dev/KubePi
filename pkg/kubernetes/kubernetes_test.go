package kubernetes

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"testing"
)

var cert = "-----BEGIN CERTIFICATE-----\nMIIDADCCAeigAwIBAgIBAjANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwptaW5p\na3ViZUNBMB4XDTIxMDMwODA4MjIyNloXDTIyMDMwOTA4MjIyNlowMTEXMBUGA1UE\nChMOc3lzdGVtOm1hc3RlcnMxFjAUBgNVBAMTDW1pbmlrdWJlLXVzZXIwggEiMA0G\nCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCmp2mjHtThDERXsKjGbTBEr/Do7xGy\nN7nxzyruZykzWERUAVSX6QZgg73ROUdsXFR9Xtg4d7UwoV2xrE0iXOPSnvKjnMtU\n5/buyqGNPXQmrUwfxVmBgABxZEIGRnWU8UXMGKPA3WbkG4P2CKEPz390fNxaD+iC\ntGvjZ9Vvtkt7Tzh/FHUt98+rLZHnD+66EX5I1ketm/QcsNMH3jIAoSIcxoFBuySS\nKWzuLMGb2cuyLJE1rxAoUI79wLnlR7AptJMjCRthZiOnpoGYfqYcIFiJ+bj69OHg\nMwvQRJW2wA9YmNjqXhMQ9u5DP6ynQaSMwc2WduHN12jCdLBiMKS//4PtAgMBAAGj\nPzA9MA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUH\nAwIwDAYDVR0TAQH/BAIwADANBgkqhkiG9w0BAQsFAAOCAQEAGS6Is5lsXcRwrLEx\nNJszecVEGwWJpqf9I0ZkCGkBNyZe5vCWERAPyhWHXDCLHv2n5sUfPyaBMXoWM6LE\nfkOEj7877HgEREjf35Zqd8cYCD//17/g8VQVqEaCPOmWwJgQpTj2v+v3srQoPbFW\nMVrzZSiYvACwzkvKrUby1kltyKMD9jNA0tL3sj1Dqi04wDPS7YQZaV3zGUZCmCe5\nOu8rD2J53DLBs5rz4op/XQZz3h4xuyYEs681T9itS6ZLRgZRpmuPoS6vIxk9TqLT\nhAjA7pmPzOJaBkzQchIMmiC5MlpCrBozIGDD0arpoC02Qscvde9ZyTa9GOQioQ10\nMbRkCA==\n-----END CERTIFICATE-----"
var key = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEApqdpox7U4QxEV7Coxm0wRK/w6O8Rsje58c8q7mcpM1hEVAFU\nl+kGYIO90TlHbFxUfV7YOHe1MKFdsaxNIlzj0p7yo5zLVOf27sqhjT10Jq1MH8VZ\ngYAAcWRCBkZ1lPFFzBijwN1m5BuD9gihD89/dHzcWg/ogrRr42fVb7ZLe084fxR1\nLffPqy2R5w/uuhF+SNZHrZv0HLDTB94yAKEiHMaBQbskkils7izBm9nLsiyRNa8Q\nKFCO/cC55UewKbSTIwkbYWYjp6aBmH6mHCBYifm4+vTh4DML0ESVtsAPWJjY6l4T\nEPbuQz+sp0GkjMHNlnbhzddownSwYjCkv/+D7QIDAQABAoIBAA36ctDc0CxENCNK\nzQ0/sVHBlCpliw1wwSb4Ini2rG0DFVAagHbxc7h6tFwtOsFrCScu4mHyIH+AuXQi\nqKGeOvm6nU195Ewt3LdwxZYsNmbcGEt96SEElITuTN9r34brqdgRpQKTT4MIj8v+\nM0w5Mk6Z/n2LYgw8h+QeHHfvSWuBOdlNCMqx2XI3gGc4wm3vub3t6pGz+05QPbCD\nd/drVyMVt3yIA2HHZawKZ5WxEk2aEhAg9wqNuVcztxEmnnCTtaxtkg9S6FxeQtuK\nq2WvznCUgBCvYQpzNAOolBhjz2zlyKu7/I4F6+ZwpVPwhfvgayYYLWtkFaNPgLXe\nm9V4w90CgYEAwWcqJgI2QA4LT7XbCEXBkO0mvvWX/pO0kqvHaqVSaSNe+y/fseI0\nq3+A97sZ2MDXnpMl+f0z7q8HZecEHp1NhetbCgHTpgBD7tFaJJCEa1PnFzSGCyFL\n1r5NTHLGB01W16bH0+WMmcZM1mtsx/bctLx/K/BNKmTEpD9Gn0FP97sCgYEA3Jfl\nvbI5DiTMEtewm9C/hqOPeFf18N1Azba1YHC3P4tRVFIOHwJZp1UNsq5gSu4lyYDo\n99Odub0YIkMEZBYPHTtBKZhX90i5XvQm7jmiiKWQ8tirBWrsrljdoxopE5tudCFT\noxBJ94QMakFVcoZqJhmpTscRw3JzpN8XlOH+VHcCgYBu1eaLvbzFXMcSuU97IC7c\nFWydBzZCCPf1DkjMT045PrISFc+Gq/IvTnTkg+8+DtYC5KVg7MC0Ss5ckdYEjXV+\nB/E2fPGEMqa72HJmfgPFVmIbJFilTEGgIZM++o+OY74e/E+MmgLHpaMnRo0i09CM\nK3JeBerTHsiqsDCS1+UyPQKBgQC1t39u0/kCOLfPsdRvlvefTu9qAHO+NlUi4Sbq\nyg96jiayImI1kzcNjBgboF/8ec+w/btsI+vjTO0rlC9yz2Ul/GEChde5AjSKDvBf\nACVvEYylMG05qkpMmTIDIRLDbx//FFEUm9+CwUmE4kska6vXtP3uwjhU29x97bU8\nVSqwowKBgAPjfYKh2xawWw+tPOJ1INalYmjBY/eB0ZE/YT4hD/DAO8Y1+tpM9HX5\nJNtYwOrG/cWMahFAFW82SNZ0d9xairhZ+cOGZ4BQVUviMiY547GQAVwmj+nwqpDT\ngq+oQc1yMo69JpMcpQIKNdXB7HPFLHFdSVgY/cfEbuZhhSW2tHeZ\n-----END RSA PRIVATE KEY-----"
var ca = "-----BEGIN CERTIFICATE-----\nMIIC5zCCAc+gAwIBAgIBATANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwptaW5p\na3ViZUNBMB4XDTIxMDIxODA2MjYzMFoXDTMxMDIxNzA2MjYzMFowFTETMBEGA1UE\nAxMKbWluaWt1YmVDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMXs\neb3dt6ENhYL61IeVV0vYSKobqjL+TvAPEnXPySv6BjT8ejxz8mM6qP/SVrTJnC0H\nGOtiQgRDCrRlsfzce2DsMlMR2XzgQlao0bF5h0pJdkSdvWnJT0rlnca9Fjt6jVB6\nV+WiceFM5U/LiO8fmdp9R1frefxWfYTaT/80XUo0ZVEgZUBu1QFgAfjexal0CZEi\nBFzRMddW6YEUSrHd46X79pXBz3jXZ228q8XRfkAkyuMZaCAa+poIhVbaItHRgJbd\nUtXi24G/Imt7KddDY8Q3MrlgPeFPcBRy5re4trVCJAK+U9+2Oagbe0eEH7wjB2d7\n9cBMKrUKpRpCP/lMsj8CAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgKkMB0GA1UdJQQW\nMBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3\nDQEBCwUAA4IBAQCH6yn1WbP9zMks6vlxxiqvF3Vg6xMBB+5iZ3mDrBkEYntGpgrC\nBnGnevMWDXPEZsJRUb8rEJChPpk/zB49Kc+FSWJJXTPMvtCqRr4HJoaxxyYlfjfV\nXQsScaDR2LJy9P57VbUeAOaUqN1BiJfK895C8RN1RPgR3tW92eFz5OM/oOF58Km0\nQ6enpBEeNkMpSS0gpd31fs5aL1xqn2zxIVS+cQDJK/4ct18+XSPasCXKNZmyuldJ\n/dZk2DjdR3rqwqp5QbRYBvlKlF+xdvxgM3UwWPfu/bLjo/jdv2L1ThaIPCLuUqQg\n0iBw+d0wj1Zdb6ipJ7QffjCoHtMXMI5VLDwI\n-----END CERTIFICATE-----"

func TestNewKubernetes(t *testing.T) {
	//cs, err := clientcmd.BuildConfigFromFlags("", "/Users/shenchenyang/.kube/config")
	//if err != nil {
	//	t.Error(err)
	//}

	cs := rest.Config{
		Host: "https://127.0.0.1:51404",
		TLSClientConfig: rest.TLSClientConfig{
			CertData: []byte(cert),
			KeyData:  []byte(key),
			Insecure: true,
		},
	}

	cls, err := kubernetes.NewForConfig(&cs)
	if err != nil {
		t.Error(err)
	}

	list, err := cls.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(len(list.Items))
}
