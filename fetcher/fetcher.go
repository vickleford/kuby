package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Fetcher struct {
	client *http.Client
	dest   io.WriteCloser
}

func (f *Fetcher) setClient(c *http.Client) {
	f.client = c
}

// func (f *Fetcher) setDest(d io.WriteCloser) {
// 	f.dest = d
// }

func (f *Fetcher) Pull(version string) {
	path := fmt.Sprintf(
		"https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/kubectl",
		version)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %s\n", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing request: %s\n", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(f.dest, resp.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to destination: %s\n", err)
	}
}

func New(dest io.WriteCloser) *Fetcher {
	client := new(http.Client)
	f := new(Fetcher)
	f.client = client
	f.dest = dest
	return f
}
