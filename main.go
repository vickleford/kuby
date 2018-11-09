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
		fmt.Fprintf(os.Stderr, "Error opening %s: %s\n",
			opts.ConfigFile, err)
		os.Exit(1)
	}

	var context ctxmgr.ContextManager
	if opts.Context == "" {
		context = ctxmgr.New(kubeconfig)
	} else {
		context = ctxmgr.NewWithContext(kubeconfig, opts.Context)
	}

	httpclient := new(http.Client)

	transportConfig := http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	insecureclient := &http.Client{Transport: &transportConfig}

	clusterclient := kubyclient.New(context, insecureclient)
	version := clusterclient.ClusterVersion()

	kubectlpath := os.ExpandEnv(fmt.Sprintf("${HOME}/.kubytest/kubectl-%s", version))
	if _, err := os.Stat(kubectlpath); os.IsNotExist(err) {
		fmt.Printf("Downloading %s to %s...\n", version, kubectlpath)
		// need to make it check for already existing files and not download again.
		downloadpath, err := os.OpenFile(kubectlpath, os.O_CREATE|os.O_WRONLY, 0755)
		// handle directory does not exist
		// handle correct expansion of ~
		// download with progress indicator https://golangcode.com/download-a-file-with-progress/
		// maybe this stuff can go into fetcher.DownloadManager
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file %s: %s\n",
				"placeholder", err)
			os.Exit(1)
		}

		downloader := fetcher.New(downloadpath, httpclient)
		downloader.Pull(version)
		downloadpath.Close()
	}

	cmd := kubectl.New(kubectlpath)
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1) // should exit what kubectl exits, not always 1.
	}
}
