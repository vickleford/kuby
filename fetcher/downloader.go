package fetcher

import (
	"fmt"
	"net/http"
	"os"
)

// check for already existing files and not download again
// handle directory does not exist
// handle correct expansion of ~
// download with progress indicator https://golangcode.com/download-a-file-with-progress/
// make sure the download doesn't leave around an unusable kubectl on error downloading

type DownloadManager interface {
	Download(version string) (string, error)
}

type downloadManager struct {
	client *http.Client
}

func (dm *downloadManager) Download(version string) (string, error) {
	destdir := os.ExpandEnv("${HOME}/.kuby")
	if _, err := os.Stat(destdir); os.IsNotExist(err) {
		if err = os.MkdirAll(destdir, 0755); err != nil {
			return "", err
		}
	}

	kubectlpath := fmt.Sprintf("%s/kubectl-%s", destdir, version)
	if _, err := os.Stat(kubectlpath); os.IsNotExist(err) {
		fmt.Printf("Downloading %s to %s...\n", version, kubectlpath)

		downloadpath, err := os.OpenFile(kubectlpath, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			return "", err
		}
		defer downloadpath.Close()

		downloader := New(downloadpath, dm.client)
		err = downloader.Pull(version)
		if err != nil {
			defer os.Remove(kubectlpath)
			return "", err
		}
	}

	return kubectlpath, nil
}

func NewManager(client *http.Client) DownloadManager {
	dm := new(downloadManager)
	dm.client = client
	return dm
}
