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

	dlmanager := fetcher.NewManager(httpclient)
	kubectlpath, err := dlmanager.Download(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading kubectl: %s\n", err)
		os.Exit(1)
	}

	cmd := kubectl.New(kubectlpath)
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1) // should exit what kubectl exits, not always 1.
	}
}
