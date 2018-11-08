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
	Clusters       []struct {
		Name    string
		Cluster struct {
			Server string
		}
	}
	Contexts []struct {
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
	Endpoint() string
}

type contextManager struct {
	kubeconfig *Kubeconfig
	config     io.Reader
	username   string
	password   string
	cluster    string
	endpoint   string
}

func (c *contextManager) Username() string {
	return c.username
}

func (c *contextManager) Password() string {
	return c.password
}

func (c *contextManager) Endpoint() string {
	return c.endpoint
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

func (c *contextManager) loadEndpoint() {
	for _, cluster := range c.kubeconfig.Clusters {
		if cluster.Name == c.cluster {
			c.endpoint = cluster.Cluster.Server
		}
	}
}

func (c *contextManager) parse() {
	var contextUser string
	if c.kubeconfig.CurrentContext == "" {
		contextUser = c.kubeconfig.Contexts[0].Context.User
		c.cluster = c.kubeconfig.Contexts[0].Context.Cluster
		c.loadUser(contextUser)
		c.loadEndpoint()
	} else {
		c.parseContext(c.kubeconfig.CurrentContext)
	}
}

func (c *contextManager) parseContext(wantedContext string) {
	var contextUser string
	for _, context := range c.kubeconfig.Contexts {
		if context.Name == wantedContext {
			contextUser = context.Context.User
			c.cluster = context.Context.Cluster
			break
		}
	}

	c.loadUser(contextUser)
	c.loadEndpoint()
}

func New(conf io.Reader) ContextManager {
	ctx := new(contextManager)
	ctx.config = conf
	ctx.loadConfig()
	ctx.parse()
	return ctx
}

func NewWithContext(conf io.Reader, context string) ContextManager {
	ctx := new(contextManager)
	ctx.config = conf
	ctx.loadConfig()
	ctx.parseContext(context)
	return ctx
}
