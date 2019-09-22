package secrets

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func TestParameterStore_Decode(t *testing.T) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.SkipNow()
	}
	if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		t.SkipNow()
	}
	if os.Getenv("AWS_REGION") == "" {
		t.SkipNow()
	}
	if os.Getenv("SSM_PATH") == "" {
		t.SkipNow()
	}

	var (
		s   = session.Must(session.NewSession(aws.NewConfig()))
		api = ssm.New(s)
	)

	decoder, err := NewParameterStore(WithParameterStore(api))
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	m := map[string]interface{}{}
	err = decoder.Decode(os.Getenv("SSM_PATH"), &m)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if got := len(m); got == 0 {
		t.Fatalf("got 0; want > 0")
	}
}
