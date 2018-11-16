package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/vickleford/kuby/ctxmgr"
	"github.com/vickleford/kuby/fetcher"
	"github.com/vickleford/kuby/kubectl"
	"github.com/vickleford/kuby/kubyclient"
	"github.com/vickleford/kuby/options"
)

func main() {
	opts := options.New(os.Args)

	kubeconfig, err := os.Open(opts.ConfigFile)
	defer kubeconfig.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	var context ctxmgr.ContextManager
	if opts.Context == "" {
		context, err = ctxmgr.New(kubeconfig)
	} else {
		context, err = ctxmgr.NewWithContext(kubeconfig, opts.Context)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %s\n", err)
		os.Exit(1)
	}

	httpclient := new(http.Client)

	transportConfig := http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	insecureclient := &http.Client{Transport: &transportConfig}

	clusterclient := kubyclient.New(context, insecureclient)
	version, err := clusterclient.ClusterVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	destdir := os.ExpandEnv("${HOME}/.kuby")
	if _, err := os.Stat(destdir); os.IsNotExist(err) {
		if err = os.MkdirAll(destdir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	kubectlpath := fmt.Sprintf("%s/kubectl-%s", destdir, version)
	if _, err := os.Stat(kubectlpath); os.IsNotExist(err) {
		fmt.Printf("Downloading %s to %s...\n", version, kubectlpath)

		downloadpath, err := os.OpenFile(kubectlpath, os.O_CREATE|os.O_WRONLY, 0755)
		// check for already existing files and not download again
		// handle directory does not exist
		// handle correct expansion of ~
		// download with progress indicator https://golangcode.com/download-a-file-with-progress/
		// maybe this stuff can go into fetcher.DownloadManager
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file %s: %s\n",
				kubectlpath, err)
			os.Exit(1)
		}

		downloader := fetcher.New(downloadpath, httpclient)
		err = downloader.Pull(version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing kubectl to location %s: %s\n",
				kubectlpath, err)
			downloadpath.Close()
			os.Remove(kubectlpath)
			os.Exit(1)
		}
		downloadpath.Close()
	}

	cmd := kubectl.New(kubectlpath)
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1) // should exit what kubectl exits, not always 1.
	}
}
