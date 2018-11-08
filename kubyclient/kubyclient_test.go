package kubyclient

import (
	"testing"

	"github.com/vickleford/kuby/httpclienttest"
)

func TestGetVersion(t *testing.T) {
	expected := "v1.10.0"
	testClient, _ := httpclienttest.New(niceResponse)
	client := New("https://api.k8s.example.com", "admin", "shhh", testClient)

	if actual := client.ClusterVersion(); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestBasicAuthUsage(t *testing.T) {
	expectedUsername := "admin"
	expectedPassword := "shhh"
	testClient, spy := httpclienttest.New(niceResponse)
	client := New("https://api.k8s.example.com", expectedUsername, expectedPassword, testClient)
	client.ClusterVersion()
	actualUsername, actualPassword, _ := spy.Request.BasicAuth()
	if actualUsername != expectedUsername {
		t.Errorf("Expected %s, got %s", expectedUsername, actualUsername)
	}

	if actualPassword != expectedPassword {
		t.Errorf("Expected %s, got %s", expectedPassword, actualPassword)
	}
}

func TestUserAgentSet(t *testing.T) {
	expected := "kuby"
	testClient, spy := httpclienttest.New(niceResponse)
	client := New("https://api.k8s.example.com", "admin", "shhh", testClient)
	client.ClusterVersion()
	if actual := spy.Request.Header.Get("User-Agent"); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestAcceptHeader(t *testing.T) {
	expected := "application/json"
	testClient, spy := httpclienttest.New(niceResponse)
	client := New("https://api.k8s.example.com", "admin", "shhh", testClient)
	client.ClusterVersion()
	if actual := spy.Request.Header.Get("Accept"); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

const niceResponse = `{
	"major": "1",
	"minor": "10",
	"gitVersion": "v1.10.0",
	"gitCommit": "098570796b32895c38a9a1c9286425fb1ececa18",
	"gitTreeState": "clean",
	"buildDate": "2018-08-02T17:11:51Z",
	"goVersion": "go1.9.3",
	"compiler": "gc",
	"platform": "linux/amd64"
}
`
