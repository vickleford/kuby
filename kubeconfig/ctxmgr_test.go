package kubeconfig

import (
	"os"
	"testing"
)

func TestContextDetectionWithContext(t *testing.T) {
	config, err := os.Open("../fixtures/sample_with_context.conf")
	if err != nil {
		t.Errorf("Error opening test fixture: %s", err)
	}

	// should i pass in the context?
	// no, i think it should figure out the context from the file.
	// one trick to remember is that i need to think about how
	// a context override from command-line trumps this.
	// maybe a diff 'constructor'? so maybe New(config io.Reader)
	// and NewWithContext(config io.reader, ctx string)
	ctxmgr := NewContextManager(config)
	expectedUsername := "admin"
	expectedPassword := "wheredyougetthosepeepers"

	if ctxmgr.Username != expectedUsername {
		t.Errorf("Expected username %s, got %s",
			expectedUsername, ctxmgr.Username)
	}

	if ctxmgr.Password != expectedPassword {
		t.Errorf("Expected password %s, got %s",
			expectedPassword, ctxmgr.Password)
	}
}
