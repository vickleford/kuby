package main

import flag "github.com/spf13/pflag"
import "os"

type ArgTranslator struct {
	args       []string
	ConfigFile string
}

func (a *ArgTranslator) Parse() {
	flagset := flag.NewFlagSet(a.args[0], flag.ExitOnError)
	kubeconfFlag := flagset.String("kubeconfig", "", "Specify an alternative path to kubeconfig")
	flagset.Parse(a.args)

	if *kubeconfFlag != "" {
		a.ConfigFile = *kubeconfFlag
	} else if os.Getenv("KUBECONFIG") != "" {
		a.ConfigFile = os.Getenv("KUBECONFIG")
	} else {
		a.ConfigFile = "~/.kube/config"
	}
}

func NewArgTranslator(args []string) *ArgTranslator {
	argtr := new(ArgTranslator)
	argtr.args = args
	argtr.Parse()
	return argtr
}
