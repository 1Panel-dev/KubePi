package ldap

import (
	"errors"
	"fmt"
	"github.com/go-ldap/ldap"
)

type Ldap struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Conn     *ldap.Conn
}

func NewLdapClient(address, port, username, password string) *Ldap {
	return &Ldap{
		Username: username,
		Address:  address,
		Password: password,
		Port:     port,
	}
}

func (l *Ldap) Connect() error {
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", l.Address, l.Port))
	if err != nil {
		return err
	}
	if err := conn.Bind(l.Username, l.Password); err != nil {
		return err
	}
	l.Conn = conn
	return nil
}

func (l *Ldap) Search(dn, filter string, attributes []string) ([]*ldap.Entry, error) {
	searchRequest := ldap.NewSearchRequest(dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil)
	sr, err := l.Conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	if len(sr.Entries) == 0 {
		return nil, errors.New("user is not found")
	}
	defer l.Conn.Close()
	return sr.Entries, err
}

func (l *Ldap) Login(dn, filter, password string) error {
	searchRequest := ldap.NewSearchRequest(dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"dn", "cn", "uid"},
		nil)
	sr, err := l.Conn.Search(searchRequest)
	if err != nil {
		return err
	}
	if len(sr.Entries) != 1 {
		return errors.New("user is not found")
	}
	userdn := sr.Entries[0].DN
	err = l.Conn.Bind(userdn, password)
	if err != nil {
		return err
	}
	defer l.Conn.Close()
	return nil
}
