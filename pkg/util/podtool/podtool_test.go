package podtool

import (
	"fmt"
	"github.com/KubeOperator/kubepi/pkg/util/kubernetes"
	"testing"
)

func TestPodtool(t *testing.T) {
	config := kubernetes.Config{
		Host:  "https://172.16.10.83:8443",
		Token: "eyJhbGciOiJSUzI1NiIsImtpZCI6Ik80eGVHUVBSMzZZdDhvMVNDNmRvVHBJSkFMdUg1dzhQUWJ1LUtNX1dKdkEifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrby1hZG1pbi10b2tlbi10OXY1OCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6Ijc4ZTU3NzE4LTIxNzAtNDIzOC04OGQ1LWJhMDNiYjM4MjJlMiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprby1hZG1pbiJ9.m55qWGuVYfcLBbKVDKvTmRMNxLUsgYCjqgV-s1BN7I6Vgrlv3XdnbiRAq7GGGCg_3pLDpW_HKBwDv87kHx5cf-_lOKfKSnp4UwTNabJ8kLCcJaVYP2ov3BGYcbPoXv9zjlJ2EZkcJszNRzEkNgbLUK4WrG9o2ek8yNCuNLtpc1GY7xQB3eaIhyakCxnxMs3rOTzkVe8YDoSmBL_iabhZixWKWtTYIMBnh5pl7PNxwcTFLHfbytSGs-jiW62ybgkPFE6nOpKENym0y03FBmIXUyVV26isnuzbdeS4buvTJu8PGTtODzDqLXQQHZMqlwvlv1xRDRPTiA4LW_R2UQxlJQ",
	}
	client, err := kubernetes.NewKubernetesClient(&config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	config2 := kubernetes.NewClusterConfig(&config)

	tool := PodTool{
		K8sClient:     client,
		RestClient:    config2,
		Namespace:     "kube-system",
		PodName:       "ingress-nginx-controller-77xtt",
		ContainerName: "ingress-nginx-controller",
	}

	//res,err := tool.ExecCommand([]string{"ls","-lF","--full-time", "/etc/nginx"})
	//command := "echo '" +"test1111"+ "' >>  /tmp/files9"
	command2 := "cat /dev/null > /tmp/files9"
	commands := []string{"sh", "-c", command2}
	res, err := tool.ExecCommand(commands)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ss := string(res)
	fmt.Println(ss)
	//array := strings.Split(ss,"\n")
	//for _,line := range array {
	//	msgs := strings.Fields(line)
	//	if len(msgs) <= 2 || (msgs[0] == "total" && len(msgs) ==2)   {
	//		continue
	//	}
	//
	//
	//	fmt.Println(msgs[5],msgs[6],msgs[7])
	//}

}
