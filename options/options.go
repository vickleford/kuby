package options

import flag "github.com/spf13/pflag"
import "os"

type ArgTranslator struct {
	args       []string
	ConfigFile string
	Context    string
}

func (a *ArgTranslator) Parse() {
	flagset := flag.NewFlagSet(a.args[0], flag.ExitOnError)
	kubeconfFlag := flagset.String("kubeconfig", "", "Specify an alternative path to `kubeconfig`")
	contextFlag := flagset.String("context", "", "Use context from `kubeconfig`")
	flagset.Parse(a.args)

	if *kubeconfFlag != "" {
		a.ConfigFile = *kubeconfFlag
	} else if os.Getenv("KUBECONFIG") != "" {
		a.ConfigFile = os.Getenv("KUBECONFIG")
	} else {
		a.ConfigFile = os.ExpandEnv("${HOME}/.kube/config")
	}

	a.Context = *contextFlag
}

func New(args []string) *ArgTranslator {
	argtr := new(ArgTranslator)
	argtr.args = args
	argtr.Parse()
	return argtr
}
