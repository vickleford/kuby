package kubyclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ClusterVersion struct {
	GitVersion string
}

type KubyClient interface {
	ClusterVersion() string
}

type kubyclient struct {
	clusterurl string
	username   string
	password   string
	client     *http.Client
}

func (k *kubyclient) ClusterVersion() string {
	path := fmt.Sprintf("%s/version", k.clusterurl)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %s\n", err)
	}
	req.SetBasicAuth(k.username, k.password)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "kuby")

	resp, err := k.client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing request: %s\n", err)
	}
	defer resp.Body.Close()

	var v ClusterVersion
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting response to JSON: %s\n", err)
	}

	return v.GitVersion
}

func New(clusterurl, username, password string, c *http.Client) KubyClient {
	kubyclient := new(kubyclient)
	kubyclient.clusterurl = clusterurl
	kubyclient.username = username
	kubyclient.password = password
	// kubyclient.client = new(http.Client)
	kubyclient.client = c
	return kubyclient
}
