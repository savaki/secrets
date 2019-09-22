package secrets

import (
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type Encoder interface {
	Encode(path string, v interface{}) error
}

type options struct {
	secretsManager secretsmanageriface.SecretsManagerAPI
	ssm            ssmiface.SSMAPI
}

type Option func(o *options)

func WithParameterStore(api ssmiface.SSMAPI) Option {
	return func(o *options) {
		o.ssm = api
	}
}

func WithSecretsManager(api secretsmanageriface.SecretsManagerAPI) Option {
	return func(o *options) {
		o.secretsManager = api
	}
}
