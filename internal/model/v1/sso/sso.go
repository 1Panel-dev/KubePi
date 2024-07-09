package sso

import (
	"context"
	v1 "github.com/KubeOperator/kubepi/internal/model/v1"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type Sso struct {
	v1.BaseModel     `storm:"inline"`
	v1.Metadata      `storm:"inline"`
	Enable           bool   `json:"enable"`
	Protocol         string `json:"protocol"`
	InterfaceAddress string `json:"interfaceAddress"`
	ClientId         string `json:"clientId"`
	ClientSecret     string `json:"clientSecret"`
}

type OpenID struct {
	Code         string
	Language     string
	Oauth2Config *oauth2.Config
	OidcProvider *oidc.Provider
	Ctx          context.Context
}
