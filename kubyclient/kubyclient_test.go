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

type eksContextManager struct{}

func (m eksContextManager) Username() string {
	return ""
}

func (m eksContextManager) Password() string {
	return ""
}

func (m eksContextManager) Endpoint() string {
	return "https://example.eks.amazonaws.com"
}

var mockCtxMgr mockContextManager
var ctxMgrForEks eksContextManager

func TestGetVersion(t *testing.T) {
	expected := "v1.10.0"
	testClient, _ := httpclienttest.New(niceResponse)
	client := New(mockCtxMgr, testClient)

	if actual, _ := client.ClusterVersion(); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestGetVersionWithEksResponse(t *testing.T) {
	expected := "v1.11.5"
	testClient, _ := httpclienttest.New(eksResponse)
	client := New(ctxMgrForEks, testClient)

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

func TestEksConfigDoesNotAuthenticate(t *testing.T) {
	testClient, spy := httpclienttest.New(eksResponse)
	client := New(ctxMgrForEks, testClient)
	client.ClusterVersion()
	_, _, k := spy.Request.BasicAuth()
	if k != false {
		t.Error("Expected no basic auth")
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

func TestEmptyVersionIsError(t *testing.T) {
	const emptyResponse = ""
	testClient, _ := httpclienttest.New(emptyResponse)
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

const eksResponse = `{
	"major": "1",
	"minor": "11+",
	"gitVersion": "v1.11.5-eks-6bad6d",
	"gitCommit": "6bad6d9c768dc0864dab48a11653aa53b5a47043",
	"gitTreeState": "clean",
	"buildDate": "2018-12-06T23:13:14Z",
	"goVersion": "go1.10.3",
	"compiler": "gc",
	"platform": "linux/amd64"
}
`
