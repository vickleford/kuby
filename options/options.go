package options

import (
	"os"
)

type argument struct {
	Name     string // eg kubeconfig for --kubeconfig
	Value    string // determined value
	Default  string // if not given on cmdline
	Override string // environment variable name
}

type ArgTranslator struct {
	args       []string
	arguments  []*argument
	ConfigFile string
	Context    string
}

func (a *ArgTranslator) Add(name, vdefault, override string) {
	arg := new(argument)
	arg.Name = name
	arg.Default = vdefault
	arg.Override = override
	a.arguments = append(a.arguments, arg)
}

func (a *ArgTranslator) Get(key string) string {
	// might should check if it has been parsed yet
	for _, a := range a.arguments {
		if a.Name == key {
			return a.Value
		}
	}
	return ""
}

func (a *ArgTranslator) Parse() {
	for _, flag := range a.arguments {
		flaglen := len(flag.Name) + 2 // adjust for preceeding "--"
		for i, v := range a.args {
			arglen := len(v) // already has preceeding "--"
			if arglen < flaglen {
				continue
			}
			if arglen == flaglen && v == "--"+flag.Name {
				flag.Value = a.args[i+1]
				break
			}
			if arglen > flaglen && v[:flaglen+1] == "--"+flag.Name+"=" {
				flag.Value = v[flaglen+1:]
				break
			}
		}
	}

	for _, flag := range a.arguments {
		if flag.Value != "" {
			// already set, do nothing.
		} else if ev := os.Getenv(flag.Override); ev != "" {
			flag.Value = os.ExpandEnv(ev)
		} else {
			flag.Value = os.ExpandEnv(flag.Default)
		}
	}
}

func New(args []string) *ArgTranslator {
	argtr := new(ArgTranslator)
	argtr.args = args
	return argtr
}
