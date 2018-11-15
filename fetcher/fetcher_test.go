package fetcher

import (
	"fmt"
	"testing"

	"github.com/vickleford/kuby/filefaker"
	"github.com/vickleford/kuby/httpclienttest"
)

func TestFetcherKnowsTheRightPlace(t *testing.T) {
	fileFaker := filefaker.New()
	expectedScheme := "https"
	expectedHost := "storage.googleapis.com"
	expectedPath := "/kubernetes-release/release/v1.10.0/bin/linux/amd64/kubectl"

	testClient, spy := httpclienttest.New("bonk")
	fetcher := New(fileFaker, testClient)
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

func TestFetcherFetches(t *testing.T) {
	expectedWrite := "beebop"

	fileFaker := filefaker.New()
	testClient, _ := httpclienttest.New(expectedWrite)

	fetcher := New(fileFaker, testClient)
	fetcher.Pull("v1.10.0")

	if actualWrite := fileFaker.Observe(); actualWrite != expectedWrite {
		t.Errorf("Expected content %s, got %s", expectedWrite, actualWrite)
	}
}

func TestWriteErrorReturnsError(t *testing.T) {
	fileFaker := filefaker.New()
	fileFaker.Err = fmt.Errorf("eeeek can't write")
	testClient, _ := httpclienttest.New("blah")

	fetcher := New(fileFaker, testClient)
	err := fetcher.Pull("v1.10.0")
	if err == nil {
		t.Error("Expected an error")
	}
}

// cross-platform... maybe... one day.
// curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.12.0/bin/darwin/amd64/kubectl
// curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.12.0/bin/linux/amd64/kubectl
