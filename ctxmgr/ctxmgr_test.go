package ctxmgr

import (
	"strings"
	"testing"
)

func TestNoContextFromCliCurrentContextFound(t *testing.T) {
	ctxmgr, _ := New(strings.NewReader(conf_with_context))
	expectedUsername := "admin"
	expectedPassword := "wheredyougetthosepeepers"
	expectedUrl := "https://api.fredthefriendlycluster.us-west-2.example.com"

	if user := ctxmgr.Username(); user != expectedUsername {
		t.Errorf("Expected username %s, got %s",
			expectedUsername, user)
	}

	if pass := ctxmgr.Password(); pass != expectedPassword {
		t.Errorf("Expected password %s, got %s",
			expectedPassword, pass)
	}

	if url := ctxmgr.Endpoint(); url != expectedUrl {
		t.Errorf("Expected endpoint %s, got %s", expectedUrl, url)
	}
}

func TestNoContextFromCommandLineNoCurrentContext(t *testing.T) {
	ctxmgr, _ := New(strings.NewReader(conf_no_context))
	expectedUsername := "admin"
	expectedPassword := "jeeperscreepers"
	expectedUrl := "https://api.fredthefriendlycluster.us-east-1.example.com"

	if user := ctxmgr.Username(); user != expectedUsername {
		t.Errorf("Expected username %s, got %s",
			expectedUsername, user)
	}

	if pass := ctxmgr.Password(); pass != expectedPassword {
		t.Errorf("Expected password %s, got %s",
			expectedPassword, pass)
	}

	if url := ctxmgr.Endpoint(); url != expectedUrl {
		t.Errorf("Expected endpoint %s, got %s", expectedUrl, url)
	}
}

func TestContextGivenFromCommandLineCurrentContextFound(t *testing.T) {
	ctxmgr, _ := NewWithContext(strings.NewReader(conf_with_context), "us-east-1")
	expectedUsername := "admin"
	expectedPassword := "jeeperscreepers"
	expectedUrl := "https://api.fredthefriendlycluster.us-east-1.example.com"

	if user := ctxmgr.Username(); user != expectedUsername {
		t.Errorf("Expected username %s, got %s",
			expectedUsername, user)
	}

	if pass := ctxmgr.Password(); pass != expectedPassword {
		t.Errorf("Expected password %s, got %s",
			expectedPassword, pass)
	}

	if url := ctxmgr.Endpoint(); url != expectedUrl {
		t.Errorf("Expected endpoint %s, got %s", expectedUrl, url)
	}
}

func TestContextGivenFromCommandLineNoCurrentContext(t *testing.T) {
	ctxmgr, _ := NewWithContext(strings.NewReader(conf_no_context), "us-west-2")
	expectedUsername := "admin"
	expectedPassword := "wheredyougetthosepeepers"
	expectedUrl := "https://api.fredthefriendlycluster.us-west-2.example.com"

	if user := ctxmgr.Username(); user != expectedUsername {
		t.Errorf("Expected username %s, got %s",
			expectedUsername, user)
	}

	if pass := ctxmgr.Password(); pass != expectedPassword {
		t.Errorf("Expected password %s, got %s",
			expectedPassword, pass)
	}

	if url := ctxmgr.Endpoint(); url != expectedUrl {
		t.Errorf("Expected endpoint %s, got %s", expectedUrl, url)
	}
}

func TestContextNotFoundGivesError(t *testing.T) {
	_, err := NewWithContext(strings.NewReader(conf_with_context), "hangry")
	if err == nil {
		t.Error("Expected error when unable to find context in conf")
	}
}

func TestEksContextDoesNotBlowUp(t *testing.T) {
	context, err := New(strings.NewReader(conf_eks))
	if err != nil {
		t.Error(err.Error())
	}

	if expected := ""; context.Username() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, context.Username())
	}
}

const conf_no_context = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: REDACTED
    server: https://api.fredthefriendlycluster.us-east-1.example.com
  name: fredthefriendlycluster.us-east-1.example.com
- cluster:
    certificate-authority-data: REDACTED
    server: https://api.fredthefriendlycluster.us-west-2.example.com
  name: fredthefriendlycluster.us-west-2.example.com
contexts:
- context:
    cluster: fredthefriendlycluster.us-east-1.example.com
    user: fredthefriendlycluster.us-east-1.example.com
  name: us-east-1
- context:
    cluster: fredthefriendlycluster.us-west-2.example.com
    user: fredthefriendlycluster.us-west-2.example.com
  name: us-west-2
current-context: ""
kind: Config
preferences: {}
users:
- name: fredthefriendlycluster.us-east-1.example.com
  user:
    client-certificate-data: REDACTED
    client-key-data: REDACTED
    password: jeeperscreepers
    username: admin
- name: fredthefriendlycluster.us-east-1.example.com-basic-auth
  user:
    password: jeeperscreepers
    username: admin
- name: fredthefriendlycluster.us-west-2.example.com
  user:
    client-certificate-data: REDACTED
    client-key-data: REDACTED
    password: wheredyougetthosepeepers
    username: admin
- name: fredthefriendlycluster.us-west-2.example.com-basic-auth
  user:
    password: wheredyougetthosepeepers
    username: admin
`

const conf_with_context = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: REDACTED
    server: https://api.fredthefriendlycluster.us-east-1.example.com
  name: fredthefriendlycluster.us-east-1.example.com
- cluster:
    certificate-authority-data: REDACTED
    server: https://api.fredthefriendlycluster.us-west-2.example.com
  name: fredthefriendlycluster.us-west-2.example.com
contexts:
- context:
    cluster: fredthefriendlycluster.us-east-1.example.com
    user: fredthefriendlycluster.us-east-1.example.com
  name: us-east-1
- context:
    cluster: fredthefriendlycluster.us-west-2.example.com
    user: fredthefriendlycluster.us-west-2.example.com
  name: us-west-2
current-context: us-west-2
kind: Config
preferences: {}
users:
- name: fredthefriendlycluster.us-east-1.example.com
  user:
    client-certificate-data: REDACTED
    client-key-data: REDACTED
    password: jeeperscreepers
    username: admin
- name: fredthefriendlycluster.us-east-1.example.com-basic-auth
  user:
    password: jeeperscreepers
    username: admin
- name: fredthefriendlycluster.us-west-2.example.com
  user:
    client-certificate-data: REDACTED
    client-key-data: REDACTED
    password: wheredyougetthosepeepers
    username: admin
- name: fredthefriendlycluster.us-west-2.example.com-basic-auth
  user:
    password: wheredyougetthosepeepers
    username: admin
`

const conf_eks = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: REDACTED
    server: https://xxx.yl4.us-west-2.eks.amazonaws.com
  name: arn:aws:eks:us-west-2:111111111111:cluster/eksclustername
contexts:
- context:
    cluster: arn:aws:eks:us-west-2:111111111111:cluster/eksclustername
    user: arn:aws:eks:us-west-2:111111111111:cluster/eksclustername
  name: arn:aws:eks:us-west-2:111111111111:cluster/eksclustername
current-context: arn:aws:eks:us-west-2:111111111111:cluster/eksclustername
kind: Config
preferences: {}
users:
- name: arn:aws:eks:us-west-2:111111111111:cluster/eksclustername
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - token
      - -i
      - eksclustername
      command: aws-iam-authenticator
`
