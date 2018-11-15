package kubyclient

import (
	"testing"

	"github.com/vickleford/kuby/httpclienttest"
)

type mockContextManager struct{}

func (m mockContextManager) Username() string {
	return "admin"
}

func (m mockContextManager) Password() string {
	return "shhh"
}

func (m mockContextManager) Endpoint() string {
	return "https://api.k8s.example.com"
}

var mockCtxMgr mockContextManager

func TestGetVersion(t *testing.T) {
	expected := "v1.10.0"
	testClient, _ := httpclienttest.New(niceResponse)
	client := New(mockCtxMgr, testClient)

	if actual, _ := client.ClusterVersion(); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestBasicAuthUsage(t *testing.T) {
	expectedUsername := "admin"
	expectedPassword := "shhh"
	testClient, spy := httpclienttest.New(niceResponse)
	client := New(mockCtxMgr, testClient)
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
	client := New(mockCtxMgr, testClient)
	client.ClusterVersion()
	if actual := spy.Request.Header.Get("User-Agent"); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestAcceptHeader(t *testing.T) {
	expected := "application/json"
	testClient, spy := httpclienttest.New(niceResponse)
	client := New(mockCtxMgr, testClient)
	client.ClusterVersion()
	if actual := spy.Request.Header.Get("Accept"); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestBadJsonReturnsError(t *testing.T) {
	const badResponse = "401 Unauthorized"
	testClient, _ := httpclienttest.New(badResponse)
	client := New(mockCtxMgr, testClient)
	_, err := client.ClusterVersion()
	if err == nil {
		t.Error("Expected error")
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
