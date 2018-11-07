package fetcher

import "testing"

func TestFetcherFetches(t *testing.T) {
	fetcher := New("v1.10.0") // and a file-like object too...
	fetcher.pull()

}
