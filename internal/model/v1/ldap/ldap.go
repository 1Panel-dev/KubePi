package ldap

import (
	"encoding/json"
	v1 "github.com/KubeOperator/kubepi/internal/model/v1"
)

type Ldap struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Address      string `json:"address"`
	Port         string `json:"port"`
	Dn           string `json:"dn"`
	Filter       string `json:"filter"`
	Mapping      string `json:"mapping"`
	TLS          bool   `json:"tls"`
	Enable       bool   `json:"enable"`
	SizeLimit    int    `json:"sizeLimit"`
	TimeLimit    int    `json:"timeLimit"`
}

func (l *Ldap) GetAttributes() ([]string, error) {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(l.Mapping), &m)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (l *Ldap) GetMappings() (map[string]string, error) {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(l.Mapping), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
