package main

import "testing"
import "os"

func TestNoConfigGivesDefaultKubeconf(t *testing.T) {
	expected := "~/.kube/config"
	osArgsSim := []string{"kuby", "get", "nodes"}
	args := NewArgTranslator(osArgsSim)
	if args.ConfigFile != expected {
		t.Errorf("Expected %s, got %s", expected, args.ConfigFile)
	}
}

func TestKubeconfigFlagSetsValue(t *testing.T) {
	expected := "turtlepower"
	osArgsSim := []string{"kuby", "--kubeconfig", expected, "get", "nodes"}
	args := NewArgTranslator(osArgsSim)
	if args.ConfigFile != expected {
		// displayCircumstances()
		t.Errorf("Expected %s, got %s", expected, args.ConfigFile)
	}
}

func TestRespectsKubectlConfigEnvVar(t *testing.T) {
	expected := "envvarconfig"
	os.Setenv("KUBECONFIG", expected)
	osArgsSim := []string{"kuby", "get", "nodes"}
	args := NewArgTranslator(osArgsSim)
	if args.ConfigFile != expected {
		t.Errorf("Expected %s, got %s", expected, args.ConfigFile)
	}
}

func TestKubeconfigFlagTrumpsEnvVar(t *testing.T) {
	expected := "fromflag"
	os.Setenv("KUBECONFIG", "fromenvvar")
	osArgsSim := []string{"kuby", "--kubeconfig", expected, "get", "nodes"}
	args := NewArgTranslator(osArgsSim)
	if args.ConfigFile != expected {
		// displayCircumstances()
		t.Errorf("Expected %s, got %s", expected, args.ConfigFile)
	}
}

// https://groups.google.com/forum/#!topic/Golang-Nuts/1aZmhhSvwWc
