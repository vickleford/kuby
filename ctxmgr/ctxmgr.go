package ctxmgr

import (
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

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
	err        error
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
		c.err = err
	}

	kubeconfig := new(Kubeconfig)
	err = yaml.Unmarshal(bytes, kubeconfig)
	if err != nil {
		c.err = err
	}

	c.kubeconfig = kubeconfig
}

func (c *contextManager) loadUser(user string) {
	for _, u := range c.kubeconfig.Users {
		if user == u.Name {
			c.username = u.User.Username
			c.password = u.User.Password
			return
		}
	}
}

func (c *contextManager) loadEndpoint() {
	for _, cluster := range c.kubeconfig.Clusters {
		if cluster.Name == c.cluster {
			c.endpoint = cluster.Cluster.Server
			return
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
	var foundContext bool
	for _, context := range c.kubeconfig.Contexts {
		if context.Name == wantedContext {
			foundContext = true
			contextUser = context.Context.User
			c.cluster = context.Context.Cluster
			break
		}
	}

	if !foundContext {
		c.err = newErrConfigParse(fmt.Sprintf("Context '%s' not found",
			wantedContext))
	}

	c.loadUser(contextUser)
	c.loadEndpoint()
}

func New(conf io.Reader) (ContextManager, error) {
	ctx := new(contextManager)
	ctx.config = conf
	ctx.loadConfig()
	if ctx.err == nil {
		ctx.parse()
	}
	return ctx, ctx.err
}

func NewWithContext(conf io.Reader, context string) (ContextManager, error) {
	ctx := new(contextManager)
	ctx.config = conf
	ctx.loadConfig()
	if ctx.err == nil {
		ctx.parseContext(context)
	}
	return ctx, ctx.err
}
