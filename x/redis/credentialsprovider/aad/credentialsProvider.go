package aad

import (
	"github.com/takekazu/redis-aad/x/redis/credentialsprovider"
)

type Provider struct {
	username string
}

var _ credentialsprovider.CredentialsProvider_ = &Provider{}

func New(username string) *Provider {
	return &Provider{
		username: username,
	}
}

func (p *Provider) CredentialsProvider() (username string, password string) {
	return p.username, "password"
}
