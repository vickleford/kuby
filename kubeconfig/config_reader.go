package kubeconfig

import (
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Kubeconfig struct {
	CurrentContext string `yaml:"current-context"`
	Contexts       []struct {
		Name    string
		Context struct {
			Cluster string
			User    string
		}
	}
	Users []struct {
		Name string
		User struct {
			Password string
			Username string
		}
	}
}

type ContextManager struct {
	Username string
	Password string
}

func NewContextManager(conf io.Reader) *ContextManager {
	bytes, err := ioutil.ReadAll(conf)
	if err != nil {
		// freak out
		fmt.Printf("Could not read config: %s\n", err)
	}

	kubeconfig := new(Kubeconfig)
	err = yaml.Unmarshal(bytes, kubeconfig)
	if err != nil {
		// freak out
		fmt.Printf("Could not convert to YAML: %s\n", err)
	}

	ctx := new(ContextManager)
	// if kubeconfig.CurrentContext == "" {

	// }
	fmt.Printf("Using current context '%s'", kubeconfig.CurrentContext)
	var contextUser string
	for _, c := range kubeconfig.Contexts {
		if c.Name == kubeconfig.CurrentContext {
			contextUser = c.Context.User
			fmt.Printf("Using context user '%s'\n", contextUser)
			break
		}
	}
	for _, u := range kubeconfig.Users {
		if contextUser == u.Name {
			fmt.Printf("Found context user '%s'\n", contextUser)
			ctx.Username = u.User.Username
			ctx.Password = u.User.Password
			break
		}
	}

	return ctx
}
