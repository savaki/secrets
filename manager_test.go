package secrets

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func TestManager_Decode(t *testing.T) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.SkipNow()
	}
	if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		t.SkipNow()
	}
	if os.Getenv("AWS_REGION") == "" {
		t.SkipNow()
	}
	if os.Getenv("SECRET_PATH") == "" {
		t.SkipNow()
	}

	var (
		s   = session.Must(session.NewSession(aws.NewConfig()))
		api = secretsmanager.New(s)
	)

	decoder, err := NewManager(WithSecretsManager(api))
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	m := map[string]interface{}{}
	err = decoder.Decode(os.Getenv("SECRET_PATH"), &m)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if got := len(m); got == 0 {
		t.Fatalf("got 0; want > 0")
	}
}
