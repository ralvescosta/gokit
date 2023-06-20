package auth0

import (
	"context"
	"errors"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/ralvescosta/gokit/auth"
	"github.com/ralvescosta/gokit/configs"
	"golang.org/x/oauth2"
)

type (
	auth0TokenManager struct {
		cfg *configs.Auth0Configs
		*oidc.Provider
		oauth2.Config
	}
)

func NewAuth0TokenManger(cfg *configs.Auth0Configs) auth.TokenManager {
	provider, err := oidc.NewProvider(context.Background(), "https://"+cfg.Domain+"/")
	if err != nil {
		return nil
	}

	conf := oauth2.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &auth0TokenManager{
		cfg:      cfg,
		Provider: provider,
		Config:   conf,
	}
}

func (m *auth0TokenManager) Validate(ctx context.Context, token string) (*auth.Session, error) {
	oauthToken, err := m.Exchange(ctx, token)
	if err != nil {
		return nil, err
	}

	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: m.cfg.ClientId,
	}

	IDToken, err := m.Verifier(oidcConfig).Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	var claims map[string]interface{}
	if err := IDToken.Claims(&claims); err != nil {
		return nil, err
	}

	return nil, nil
}
