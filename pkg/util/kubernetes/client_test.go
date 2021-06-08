package kubernetes

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestNewKubernetesClient(t *testing.T) {
	c, err := NewKubernetesClient(&Config{
		Host:  "",
		Token: ""
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
