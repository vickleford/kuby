package options

import "testing"
import "os"

func setup(osArgumentsSimulator []string) *ArgTranslator {
	args := New(osArgumentsSimulator)
	args.Add("kubeconfig", "${HOME}/.kube/config", "KUBECONFIG")
	args.Add("context", "", "")
	args.Parse()
	return args
}

func TestNoConfigGivesDefaultKubeconf(t *testing.T) {
	os.Unsetenv("KUBECONFIG")
	expected := os.ExpandEnv("${HOME}/.kube/config")
	args := setup([]string{"kuby", "get", "nodes"})
	if actual := args.Get("kubeconfig"); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestKubeconfigFlagSetsValue(t *testing.T) {
	expected := "turtlepower"
	args := setup([]string{"kuby", "--kubeconfig", expected, "get", "nodes"})
	if actual := args.Get("kubeconfig"); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestRespectsKubectlConfigEnvVar(t *testing.T) {
	expected := "envvarconfig"
	os.Setenv("KUBECONFIG", expected)
	args := setup([]string{"kuby", "get", "nodes"})
	if actual := args.Get("kubeconfig"); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestKubeconfigFlagTrumpsEnvVar(t *testing.T) {
	expected := "fromflag"
	os.Setenv("KUBECONFIG", "fromenvvar")
	args := setup([]string{"kuby", "--kubeconfig", expected, "get", "nodes"})
	if actual := args.Get("kubeconfig"); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestContextFlagNotPresent(t *testing.T) {
	expected := ""
	args := setup([]string{"kuby", "get", "nodes"})
	if actual := args.Get("context"); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestContextFlagDetected(t *testing.T) {
	expected := "llama"
	args := setup([]string{"kuby", "--context=llama", "get", "nodes"})
	if actual := args.Get("context"); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestLateContextFlagDetected(t *testing.T) {
	expected := "llama"
	args := setup([]string{"kuby", "get", "ingress", "--all-namespaces", "--context=llama"})
	if actual := args.Get("context"); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestContextWorksWithoutEquals(t *testing.T) {
	expected := "llama"
	args := setup([]string{"kuby", "get", "ingress", "--context", "llama"})
	if actual := args.Get("context"); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestUnknownShorthandFlagPassesThrough(t *testing.T) {
	New([]string{"kuby", "-n", "thatnamespace", "get", "pods"})
}

func TestUnknownFlagPassesThrough(*testing.T) {
	New([]string{"kuby", "--namespace", "foobar"})
}
