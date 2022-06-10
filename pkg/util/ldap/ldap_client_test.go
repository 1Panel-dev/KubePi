package ldap

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"testing"
)

func TestLdapClient(t *testing.T) {

	client := NewLdapClient("172.16.10.89", "389", "CN=zhengkun2,CN=Users,DC=ko,DC=com", "Calong@2015", false)
	err := client.Connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for i := 1; i < 1200; i++ {
		username := fmt.Sprintf("kubepi%d", i)
		email := username + "@fit2cloud.com"
		userdn := "CN=" + username + ",CN=Users,DC=ko,DC=com"
		add := ldap.NewAddRequest(userdn, nil)
		add.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "user"})
		add.Attribute("cn", []string{username})
		add.Attribute("sAMAccountName", []string{username})
		add.Attribute("mail", []string{email})
		err = client.Conn.Add(add)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
