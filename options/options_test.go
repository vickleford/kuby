package options

import "testing"
import "os"

func TestNoConfigGivesDefaultKubeconf(t *testing.T) {
	// this unintentionally depends on KUBECONFIG being unset.
	// fix that.
	expected := os.ExpandEnv("${HOME}/.kube/config")
	osArgsSim := []string{"kuby", "get", "nodes"}
	args := New(osArgsSim)
	if args.ConfigFile != expected {
		t.Errorf("Expected %s, got %s", expected, args.ConfigFile)
	}
}

func TestKubeconfigFlagSetsValue(t *testing.T) {
	expected := "turtlepower"
	osArgsSim := []string{"kuby", "--kubeconfig", expected, "get", "nodes"}
	args := New(osArgsSim)
	if args.ConfigFile != expected {
		t.Errorf("Expected %s, got %s", expected, args.ConfigFile)
	}
}

func TestRespectsKubectlConfigEnvVar(t *testing.T) {
	expected := "envvarconfig"
	os.Setenv("KUBECONFIG", expected)
	osArgsSim := []string{"kuby", "get", "nodes"}
	args := New(osArgsSim)
	if args.ConfigFile != expected {
		t.Errorf("Expected %s, got %s", expected, args.ConfigFile)
	}
}

func TestKubeconfigFlagTrumpsEnvVar(t *testing.T) {
	expected := "fromflag"
	os.Setenv("KUBECONFIG", "fromenvvar")
	osArgsSim := []string{"kuby", "--kubeconfig", expected, "get", "nodes"}
	args := New(osArgsSim)
	if args.ConfigFile != expected {
		t.Errorf("Expected %s, got %s", expected, args.ConfigFile)
	}
}

func TestContextFlagNotPresent(t *testing.T) {
	expected := ""
	osArgsSim := []string{"kuby", "get", "nodes"}
	args := New(osArgsSim)
	if args.Context != expected {
		t.Errorf("Expected %s, got %s", expected, args.Context)
	}
}

func TestContextFlagDetected(t *testing.T) {
	expected := "llama"
	osArgsSim := []string{"kuby", "--context=llama", "get", "nodes"}
	args := New(osArgsSim)
	if args.Context != expected {
		t.Errorf("Expected %s, got %s", expected, args.Context)
	}
}

func TestUnknownShorthandFlagPassesThrough(t *testing.T) {
	New([]string{"kuby", "-n", "thatnamespace", "get", "pods"})
}

func TestUnknownFlagPassesThrough(*testing.T) {
	New([]string{"kuby", "--namespace", "foobar"})
}
