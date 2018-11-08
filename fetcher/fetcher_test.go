package fetcher

import (
	"testing"

	"github.com/vickleford/kuby/httpclienttest"
)

type StringWriter struct {
	written []byte
	Err     error
	N       int // is this actually needed?
	Closed  bool
}

func (w *StringWriter) Write(p []byte) (n int, err error) {
	w.written = append(w.written, p...)
	n = len(p)
	w.N += n
	return n, err
}

func (w *StringWriter) Observe() string {
	return string(w.written)
}

func (w *StringWriter) Close() error {
	w.Closed = true
	return nil
}

func TestStringWriterObserve(t *testing.T) {
	expected := "are you sane?"
	w := StringWriter{}
	w.Write([]byte(expected))

	if actual := w.Observe(); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestStringWriterBytesWritten(t *testing.T) {
	expected := 5
	w := StringWriter{}
	bw, _ := w.Write([]byte("xxxxx"))

	if bw != expected {
		t.Errorf("Expected %d, got %d", expected, bw)
	}

	if w.N != expected {
		t.Errorf("Expected %d, got %d", expected, w.N)
	}
}

func TestFetcherFetches(t *testing.T) {
	fileFaker := new(StringWriter)
	// expectedBytes := []byte("bonk")
	expectedScheme := "https"
	expectedHost := "storage.googleapis.com"
	expectedPath := "/kubernetes-release/release/v1.10.0/bin/linux/amd64/kubectl"

	testClient, spy := httpclienttest.New("bonk")
	fetcher := New(fileFaker)
	fetcher.setClient(testClient) // let's experiment with this style
	// fetcher.setDest(fileFaker)
	fetcher.Pull("v1.10.0")
	if actualScheme := spy.Request.URL.Scheme; actualScheme != expectedScheme {
		t.Errorf("Got scheme %s, expected %s", actualScheme, expectedScheme)
	}
	if actualHost := spy.Request.URL.Host; actualHost != expectedHost {
		t.Errorf("Got host %s, expected %s", actualHost, expectedHost)
	}
	if actualPath := spy.Request.URL.Path; actualPath != expectedPath {
		t.Errorf("Got path %s, expected %s", actualPath, expectedPath)
	}

}

// https://storage.googleapis.com/kubernetes-release/release/v1.12.0/bin/linux/amd64/kubectl
