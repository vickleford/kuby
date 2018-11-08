package httpclienttest

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// TransportSpy spies on http client requests and records them for testing.
// This will get moved out to a whole new project later.
type TransportSpy struct {
	Request  *http.Request
	response *http.Response
	Err      error
}

func (t *TransportSpy) RoundTrip(r *http.Request) (*http.Response, error) {
	t.Request = r
	return t.response, t.Err
}

func New(responseBody string) (*http.Client, *TransportSpy) {
	// build a nice response and inject
	fakeResponse := new(http.Response)
	fakeResponse.Body = ioutil.NopCloser(strings.NewReader(responseBody))
	fakeTransport := new(TransportSpy)
	fakeTransport.response = fakeResponse

	// swap out the http.Client's Transport for a spy
	// client.client.Transport = fakeTransport
	testClient := new(http.Client)
	testClient.Transport = fakeTransport

	return testClient, fakeTransport
}
