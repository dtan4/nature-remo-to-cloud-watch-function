package ssm

import (
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// ClientInterface is an interface of a wrapper of SSM API client
type ClientInterface interface {
	LoadSecrets(path string) (map[string]string, error)
}

// Client is a wrapper of SSM API client
type Client struct {
	api ssmiface.SSMAPI
}

// NewClient creates new Client object with the given API client
func NewClient(api ssmiface.SSMAPI) *Client {
	return &Client{
		api: api,
	}
}
