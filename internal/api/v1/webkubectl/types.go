package webkubectl

import "k8s.io/client-go/rest"

type Session struct {
	User    string
	config  *rest.Config
	Cluster string `json:"cluster"`
}

type SessionResponse struct {
	SessionId string `json:"sessionId"`
}
