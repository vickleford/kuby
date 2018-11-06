package ctxmgr

import (
	"strings"
	"testing"
)

func TestNoContextFromCliCurrentContextFound(t *testing.T) {
	ctxmgr := New(strings.NewReader(conf_with_context))
	expectedUsername := "admin"
	expectedPassword := "wheredyougetthosepeepers"

	if user := ctxmgr.Username(); user != expectedUsername {
		t.Errorf("Expected username %s, got %s",
			expectedUsername, user)
	}

	if pass := ctxmgr.Password(); pass != expectedPassword {
		t.Errorf("Expected password %s, got %s",
			expectedPassword, pass)
	}
}

func TestNoContextFromCommandLineNoCurrentContext(t *testing.T) {
	ctxmgr := New(strings.NewReader(conf_no_context))
	expectedUsername := "admin"
	expectedPassword := "jeeperscreepers"

	if user := ctxmgr.Username(); user != expectedUsername {
		t.Errorf("Expected username %s, got %s",
			expectedUsername, user)
	}

	if pass := ctxmgr.Password(); pass != expectedPassword {
		t.Errorf("Expected password %s, got %s",
			expectedPassword, pass)
	}
}

func TestContextGivenFromCommandLineCurrentContextFound(t *testing.T) {
	ctxmgr := NewWithContext(strings.NewReader(conf_with_context), "us-east-1")
	expectedUsername := "admin"
	expectedPassword := "jeeperscreepers"

	if user := ctxmgr.Username(); user != expectedUsername {
		t.Errorf("Expected username %s, got %s",
			expectedUsername, user)
	}

	if pass := ctxmgr.Password(); pass != expectedPassword {
		t.Errorf("Expected password %s, got %s",
			expectedPassword, pass)
	}
}

func TestContextGivenFromCommandLineNoCurrentContext(t *testing.T) {
	ctxmgr := NewWithContext(strings.NewReader(conf_no_context), "us-west-2")
	expectedUsername := "admin"
	expectedPassword := "wheredyougetthosepeepers"

	if user := ctxmgr.Username(); user != expectedUsername {
		t.Errorf("Expected username %s, got %s",
			expectedUsername, user)
	}

	if pass := ctxmgr.Password(); pass != expectedPassword {
		t.Errorf("Expected password %s, got %s",
			expectedPassword, pass)
	}
}

var conf_no_context = `
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

var conf_with_context = `
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
