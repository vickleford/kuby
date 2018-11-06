package ctxmgr

import (
	"os"
	"testing"
)

func TestNoContextFromCliCurrentContextFound(t *testing.T) {
	config, err := os.Open("../fixtures/sample_with_context.conf")
	if err != nil {
		t.Errorf("Error opening test fixture: %s", err)
	}

	ctxmgr := New(config)
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
	config, err := os.Open("../fixtures/sample_no_context.conf")
	if err != nil {
		t.Errorf("Error opening test fixture: %s", err)
	}

	ctxmgr := New(config)
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
	config, err := os.Open("../fixtures/sample_with_context.conf")
	if err != nil {
		t.Errorf("Error opening test fixture: %s", err)
	}

	ctxmgr := NewWithContext(config, "us-east-1")
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
	config, err := os.Open("../fixtures/sample_no_context.conf")
	if err != nil {
		t.Errorf("Error opening test fixture: %s", err)
	}

	ctxmgr := NewWithContext(config, "us-west-2")
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
