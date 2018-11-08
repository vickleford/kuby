package fetcher

import (
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
	fetcher := New(fileFaker)
	fetcher.setClient(testClient) // let's experiment with this style
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
	fileFaker := filefaker.New()
	expectedWrite := "beebop"
	testClient, _ := httpclienttest.New("beebop")

	fetcher := New(fileFaker)
	fetcher.setClient(testClient)
	fetcher.Pull("v1.10.0")

	if actualWrite := fileFaker.Observe(); actualWrite != expectedWrite {
		t.Errorf("Expected content %s, got %s", expectedWrite, actualWrite)
	}
}

// cross-platform... maybe... one day.
// curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.12.0/bin/darwin/amd64/kubectl
// curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.12.0/bin/linux/amd64/kubectl
