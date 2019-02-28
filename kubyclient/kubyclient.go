package kubyclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/vickleford/kuby/ctxmgr"
)

type ClusterVersion struct {
	GitVersion string
}

type KubyClient interface {
	ClusterVersion() (string, error)
}

type kubyclient struct {
	context ctxmgr.ContextManager
	client  *http.Client
}

func (k *kubyclient) ClusterVersion() (string, error) {
	path := fmt.Sprintf("%s/version", k.context.Endpoint())

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", err
	}
	if k.context.Username() != "" {
		req.SetBasicAuth(k.context.Username(), k.context.Password())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "kuby")

	resp, err := k.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var v ClusterVersion
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)

	re := regexp.MustCompile(`v[0-9]\.[0-9]+\.[0-9]+`)
	if filteredVersion := re.FindStringSubmatch(v.GitVersion); len(filteredVersion) == 1 {
		v.GitVersion = filteredVersion[0]
	} else {
		return "", errors.New("Could not filter version correctly")
	}

	if err != nil {
		return "", err
	}

	return v.GitVersion, nil
}

func New(ctx ctxmgr.ContextManager, c *http.Client) KubyClient {
	kubyclient := new(kubyclient)
	kubyclient.context = ctx
	kubyclient.client = c
	return kubyclient
}
