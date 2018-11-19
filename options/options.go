package options

import (
	"os"

	flag "github.com/spf13/pflag"
)

type ArgTranslator struct {
	lookfor    []string
	args       []string
	ConfigFile string
	Context    string
}

func (a *ArgTranslator) OldParse() {
	flagset := flag.NewFlagSet(a.args[0], flag.ContinueOnError)
	kubeconfFlag := flagset.String("kubeconfig", "", "Specify an alternative path to `kubeconfig`")
	contextFlag := flagset.String("context", "", "Use context from `kubeconfig`")
	flagset.Parse(a.args[1:])

	if *kubeconfFlag != "" {
		a.ConfigFile = *kubeconfFlag
	} else if os.Getenv("KUBECONFIG") != "" {
		a.ConfigFile = os.Getenv("KUBECONFIG")
	} else {
		a.ConfigFile = os.ExpandEnv("${HOME}/.kube/config")
	}

	a.Context = *contextFlag
}

func (a *ArgTranslator) Parse() {
	values := make(map[string]string)
	a.lookfor = []string{"--context", "--kubeconfig"}
	for _, flag := range a.lookfor {
		// fmt.Printf("Looking for %s\n", flag)
		for i, v := range a.args {
			// fmt.Printf("Checking %s... ", v)
			arglen := len(v)
			if arglen < len(flag) {
				// fmt.Printf("Skipping %s\n", v)
				continue
			}
			if arglen == len(flag) && v == flag {
				// fmt.Printf("Found %s, setting to %s\n", v, a.args[i+1])
				values[v] = a.args[i+1]
				break
			}
			if arglen > len(flag) && v[:len(flag)+1] == v+"=" {
				// fmt.Printf("Found %s=...\n", v)
				values[v] = v[len(flag)+1:]
				break
			}
		}
	}
	a.Context = values["--context"]
	a.ConfigFile = values["--kubeconfig"]
}

func New(args []string) *ArgTranslator {
	argtr := new(ArgTranslator)
	argtr.args = args
	argtr.Parse()
	return argtr
}
