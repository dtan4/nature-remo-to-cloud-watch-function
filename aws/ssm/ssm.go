package ssm

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/pkg/errors"
)

// ClientInterface is an interface of a wrapper of SSM API client
type ClientInterface interface {
	LoadSecret(name string) (string, error)
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

// LoadSecret retrieves decrypted secret from Parameter Store
func (c *Client) LoadSecret(name string) (string, error) {
	resp, err := c.api.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", errors.Wrap(err, "cannot retrieve secret from Parameter Store")
	}

	return aws.StringValue(resp.Parameter.Value), nil
}
