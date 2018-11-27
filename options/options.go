package options

import (
	"os"
)

type option struct {
	name     string // eg kubeconfig for --kubeconfig
	value    string // determined value
	vdefault string // if not given on cmdline
	override string // environment variable name
}

type ArgTranslator struct {
	cliargs    []string
	options    []*option
	ConfigFile string
	Context    string
}

func (a *ArgTranslator) Add(name, vdefault, override string) {
	arg := new(option)
	arg.name = name
	arg.vdefault = vdefault
	arg.override = override
	a.options = append(a.options, arg)
}

func (a *ArgTranslator) Get(key string) string {
	// might should check if it has been parsed yet
	for _, a := range a.options {
		if a.name == key {
			return a.value
		}
	}
	return ""
}

func (a *ArgTranslator) Parse() {
	for _, flag := range a.options {
		flaglen := len(flag.name) + 2 // adjust for preceeding "--"
		for i, arg := range a.cliargs {
			arglen := len(arg) // already has preceeding "--"
			if arglen < flaglen {
				continue
			}
			if arglen == flaglen && arg == "--"+flag.name {
				flag.value = a.cliargs[i+1]
				break
			}
			if arglen > flaglen && arg[:flaglen+1] == "--"+flag.name+"=" {
				flag.value = arg[flaglen+1:]
				break
			}
		}
	}

	for _, flag := range a.options {
		if flag.value != "" {
			// already set, do nothing.
		} else if ev := os.Getenv(flag.override); ev != "" {
			flag.value = os.ExpandEnv(ev)
		} else {
			flag.value = os.ExpandEnv(flag.vdefault)
		}
	}
}

func New(args []string) *ArgTranslator {
	argtr := new(ArgTranslator)
	argtr.cliargs = args
	return argtr
}
