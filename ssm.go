package secrets

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"golang.org/x/xerrors"
)

type data map[string]interface{}

func (d data) set(key string, value interface{}) {
	d[key] = value
}

type ParameterStore struct {
	api ssmiface.SSMAPI
}

func (p *ParameterStore) fetch(path string) (map[string]interface{}, error) {
	var token *string
	data := data{}
	for {
		input := ssm.GetParametersByPathInput{
			NextToken:      token,
			Path:           aws.String(path),
			Recursive:      aws.Bool(true),
			WithDecryption: aws.Bool(true),
		}
		output, err := p.api.GetParametersByPath(&input)
		if err != nil {
			return nil, xerrors.Errorf("unable to fetch parameters with path, %v: %w", path, err)
		}

		for _, param := range output.Parameters {
			var (
				key   = *param.Name
				value = *param.Value
				rel   = strings.TrimLeft(key[len(path):], "/")
			)
			data.set(rel, value)
		}

		token = output.NextToken
		if token == nil {
			break
		}
	}

	return data, nil
}

func (p *ParameterStore) Decode(path string, v interface{}) error {
	raw, err := p.fetch(path)
	if err != nil {
		return xerrors.Errorf("unable to decode content: %w", err)
	}

	data, err := json.Marshal(raw)
	if err != nil {
		return xerrors.Errorf("unable to marshal data: %w", err)
	}

	return json.Unmarshal(data, v)
}

func NewParameterStore(opts ...Option) (*ParameterStore, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	api := o.ssm
	if api == nil {
		s, err := session.NewSession(aws.NewConfig())
		if err != nil {
			return nil, xerrors.Errorf("unable to create new *secrets.ParameterStore")
		}
		api = ssm.New(s)
	}

	return &ParameterStore{
		api: api,
	}, nil
}
