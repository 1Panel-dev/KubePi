package imagerepo

import (
	"fmt"
	"testing"
)

func TestNexus(t *testing.T) {
	client := NewClient(Config{
		Type: "Nexus",
		EndPoint: "http://172.16.10.185:8081",
		Credential: Credential{
			Username: "zhengkun",
			Password: "Calong@2015",
		},
	})
	_,err := client.ListRepos()
	if err != nil {
		fmt.Println(err)
	}
	images,err := client.ListImages("docker-hub-proxy")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(images)

}

func TestHarbor(t *testing.T)  {
	client := NewClient(Config{
		Type: "Harbor",
		EndPoint: "http://172.16.10.188",
		Credential: Credential{
			Username: "zhengkun",
			Password: "Calong@2015",
		},
	})
	repos,err := client.ListRepos()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(repos)
	images,err := client.ListImages("kubeapp")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(images)
}


