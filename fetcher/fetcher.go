package fetcher

import (
	"fmt"
	"io"
	"net/http"
)

type Fetcher interface {
	Pull(version string) error
}

type fetcher struct {
	client *http.Client
	dest   io.WriteCloser
}

func (f *fetcher) Pull(version string) error {
	path := fmt.Sprintf(
		"https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/kubectl",
		version)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f.dest, resp.Body)

	if err != nil {
		return err
	}

	return nil
}

func New(dest io.WriteCloser, client *http.Client) Fetcher {
	f := new(fetcher)
	f.client = client
	f.dest = dest
	return f
}
