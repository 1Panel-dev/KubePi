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
}
