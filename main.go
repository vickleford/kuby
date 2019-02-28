package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/vickleford/kuby/ctxmgr"
	"github.com/vickleford/kuby/fetcher"
	"github.com/vickleford/kuby/kubyclient"
	"github.com/vickleford/kuby/options"
)

func main() {
	opts := options.New(os.Args)
	opts.Add("kubeconfig", "${HOME}/.kube/config", "KUBECONFIG")
	opts.Add("context", "", "")
	opts.Parse()

	kubeconfig, err := os.Open(opts.Get("kubeconfig"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer kubeconfig.Close()

	var context ctxmgr.ContextManager
	if opts.Get("context") == "" {
		context, err = ctxmgr.New(kubeconfig)
	} else {
		context, err = ctxmgr.NewWithContext(kubeconfig, opts.Get("context"))
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %s\n", err)
		os.Exit(1)
	}

	httpclient := new(http.Client)

	transportConfig := http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
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

	err = syscall.Exec(kubectlpath, os.Args, os.Environ())
	if err != nil {
		fmt.Fprintf(os.Stderr, "kubectl error: %s\n", err)
	}
}
