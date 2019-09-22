package secrets

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"golang.org/x/xerrors"
)

type Manager struct {
	api secretsmanageriface.SecretsManagerAPI
}

func (m *Manager) Decode(path string, v interface{}) error {
	input := secretsmanager.GetSecretValueInput{
		SecretId: aws.String(path),
	}
	output, err := m.api.GetSecretValue(&input)
	if err != nil {
		return xerrors.Errorf("unable to retrieve secret, %v: %w", path, err)
	}

	data := output.SecretBinary
	if len(data) == 0 {
		data = []byte(*output.SecretString)
	}

	switch value := v.(type) {
	case *string:
		*value = string(data)

	case *[]byte:
		*value = data

	default:
		if err := json.Unmarshal(data, v); err != nil {
			return xerrors.Errorf("unable to decode secret, %v: %w", path, err)
		}
	}

	return nil
}

func NewManager(opts ...Option) (*Manager, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	api := o.secretsManager
	if api == nil {
		s, err := session.NewSession(aws.NewConfig())
		if err != nil {
			return nil, xerrors.Errorf("unable to create new *secrets.Manager")
		}
		api = secretsmanager.New(s)
	}

	return &Manager{
		api: api,
	}, nil
}
