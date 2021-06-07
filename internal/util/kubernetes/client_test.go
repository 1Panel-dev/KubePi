package kubernetes

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestNewKubernetesClient(t *testing.T) {
	c, err := NewKubernetesClient(&Config{
		Host:  "172.16.10.125:8443",
		Token: "eyJhbGciOiJSUzI1NiIsImtpZCI6Ii1odW5XaTdXY1ZTZmdRWVJFLTN4aUhLemttN21mOV83TG9jbmp5dzJtR1kifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrby1hZG1pbi10b2tlbi05c21kMiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjhjYTNlZjExLWU3YTMtNGU2Zi1iNzljLTI2NzMzZDY1MmQ5ZiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprby1hZG1pbiJ9.ntgCFo_UL8Hah69Nba__RwJbUGEvznUFAWX2XRXII5kexC_av5isAp3rJyv-Lqq_qGd9mJcDGm0TodMo4sP46qsMuKQ-LM6-lvCK77-y5wnGuqyXjT9C6-h7nOfP55gCaOX9mNY40l2fgjVVLcaaTklB9-QceRFgOzfjRZGwcExyL-JNepw78yMnRM_g1fV4w1VduBqniOgoM5uzACJPO-viyoZPmedCadMlnKhJnETshAe_zNZT9yZ90ujwprI9znrJLt6szLmQgK24EbcG4KF3T9PLyIBO_K0clG_U36Q14DowCLRenKZw9HOsU-AbAzwIGMjj2x8QohGDMZglRw",
	})
	if err != nil {
		fmt.Println(err)
	}
	l, err := c.CoreV1().ServiceAccounts("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(l)

}
