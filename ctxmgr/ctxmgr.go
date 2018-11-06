package ctxmgr

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

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

type ContextManager interface {
	Username() string
	Password() string
}

type contextManager struct {
	kubeconfig *Kubeconfig
	config     io.Reader
	username   string
	password   string
}

func (c *contextManager) Username() string {
	return c.username
}

func (c *contextManager) Password() string {
	return c.password
}

func (c *contextManager) loadConfig() {
	bytes, err := ioutil.ReadAll(c.config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read config: %s\n", err)
	}

	kubeconfig := new(Kubeconfig)
	err = yaml.Unmarshal(bytes, kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not convert to YAML: %s\n", err)
	}

	c.kubeconfig = kubeconfig
}

func (c *contextManager) loadUser(user string) {
	for _, u := range c.kubeconfig.Users {
		if user == u.Name {
			c.username = u.User.Username
			c.password = u.User.Password
			break
		}
	}
}

func (c *contextManager) parse() {
	c.loadConfig()

	var contextUser string
	if c.kubeconfig.CurrentContext == "" {
		contextUser = c.kubeconfig.Contexts[0].Context.User
	} else {
		for _, context := range c.kubeconfig.Contexts {
			if context.Name == c.kubeconfig.CurrentContext {
				contextUser = context.Context.User
				break
			}
		}
	}

	c.loadUser(contextUser)
}

func (c *contextManager) parseContext(context string) {
	c.loadConfig()

	var contextUser string
	for _, c := range c.kubeconfig.Contexts {
		if c.Name == context {
			contextUser = c.Context.User
			break
		}
	}

	c.loadUser(contextUser)
}

func New(conf io.Reader) ContextManager {
	ctx := new(contextManager)
	ctx.config = conf
	ctx.parse()
	return ctx
}

func NewWithContext(conf io.Reader, context string) ContextManager {
	ctx := new(contextManager)
	ctx.config = conf
	ctx.parseContext(context)
	return ctx
}
