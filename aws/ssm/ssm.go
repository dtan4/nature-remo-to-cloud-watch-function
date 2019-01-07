package ssm

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/pkg/errors"
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

// LoadSecrets retrieves decrypted secrets from Parameter Store
func (c *Client) LoadSecrets(path string) (map[string]string, error) {
	resp, err := c.api.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		Recursive:      aws.Bool(false),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return map[string]string{}, errors.Wrap(err, "cannot retrieve secrets from Parameter Store")
	}

	secrets := map[string]string{}

	for _, param := range resp.Parameters {
		name := aws.StringValue(param.Name)
		name = strings.TrimPrefix(name, path)
		name = strings.TrimPrefix(name, "/")

		secrets[name] = aws.StringValue(param.Value)
	}

	return secrets, nil
}
