package kubyclient

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type transportspy struct {
	request  *http.Request
	response *http.Response
	err      error
}

func (t *transportspy) RoundTrip(r *http.Request) (*http.Response, error) {
	t.request = r
	return t.response, t.err
}

func buildHappyTestClient() (*http.Client, *transportspy) {
	// build a nice response and inject
	fakeResponse := new(http.Response)
	fakeResponse.Body = ioutil.NopCloser(strings.NewReader(niceResponse))
	fakeTransport := new(transportspy)
	fakeTransport.response = fakeResponse

	// swap out the http.Client's Transport for a spy
	// client.client.Transport = fakeTransport
	testClient := new(http.Client)
	testClient.Transport = fakeTransport

	return testClient, fakeTransport
}

func TestGetVersion(t *testing.T) {
	expected := "v1.10.0"
	testClient, _ := buildHappyTestClient()
	client := New("https://api.k8s.example.com", "admin", "shhh", testClient)

	if actual := client.ClusterVersion(); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestBasicAuthUsage(t *testing.T) {
	expectedUsername := "admin"
	expectedPassword := "shhh"
	testClient, spy := buildHappyTestClient()
	client := New("https://api.k8s.example.com", expectedUsername, expectedPassword, testClient)
	client.ClusterVersion()
	actualUsername, actualPassword, _ := spy.request.BasicAuth()
	if actualUsername != expectedUsername {
		t.Errorf("Expected %s, got %s", expectedUsername, actualUsername)
	}

	if actualPassword != expectedPassword {
		t.Errorf("Expected %s, got %s", expectedPassword, actualPassword)
	}
}

func TestUserAgentSet(t *testing.T) {
	expected := "kuby"
	testClient, spy := buildHappyTestClient()
	client := New("https://api.k8s.example.com", "admin", "shhh", testClient)
	client.ClusterVersion()
	if actual := spy.request.Header.Get("User-Agent"); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestAcceptHeader(t *testing.T) {
	expected := "application/json"
	testClient, spy := buildHappyTestClient()
	client := New("https://api.k8s.example.com", "admin", "shhh", testClient)
	client.ClusterVersion()
	if actual := spy.request.Header.Get("Accept"); actual != expected {
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
