package ctxmgr

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
