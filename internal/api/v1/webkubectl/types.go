package webkubectl

import (
	"time"

	"k8s.io/client-go/rest"
)

type Session struct {
	User      string
	config    *rest.Config
	Cluster   string `json:"cluster"`
	ExpiresAt time.Time
}

type SessionResponse struct {
	Token string `json:"token"`
}
