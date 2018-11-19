package kubectl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// something else needs to find *which* kubectl to run....
// so what do i want this to do?
// run a kubectl command, pass through all os.Args[1:]
// stream the stdout, stderr
// care if the kubectl returns an error? just so kuby can exit unsuccessfully
// possibly pass through os.Env to the command it runs.
// need to make sure it expands environment variables. pass everything or just KUBECONFIG?

type Kubectl interface {
	Run() error
}

// i might need to hook up stdin later...
// eg: cat foo.yaml | kuby create pod
type kubectl struct {
	stderrPipe io.ReadCloser
	stdoutPipe io.ReadCloser
	stdinPipe  io.WriteCloser
	cmd        *exec.Cmd
}

// huh... one thing i could do *for testability* is just return a wrapper
// that somebody else can run?????
func (k *kubectl) wrap(path string) {
	k.cmd = exec.Command(path, os.Args[1:]...)

	cpy := os.Environ()
	envvars := make([]string, len(cpy))
	for i, envvar := range cpy {
		envvars[i] = envvar
	}
	// looks like the env vars don't need to be expanded here.
	k.cmd.Env = envvars
}

func (k *kubectl) Run() (err error) {
	// not sure if i want these here or in wrap
	k.stderrPipe, err = k.cmd.StderrPipe()
	if err != nil {
		return err
	}
	k.stdoutPipe, err = k.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	k.stdinPipe, err = k.cmd.StdinPipe()
	if err != nil {
		return err
	}

	// set up scanners
	stdoutScanner := bufio.NewScanner(k.stdoutPipe)
	stderrScanner := bufio.NewScanner(k.stderrPipe)
	stdinScanner := bufio.NewScanner(os.Stdin)

	// stream pipes
	go func() {
		for stdoutScanner.Scan() {
			fmt.Printf("%s\n", stdoutScanner.Text())
		}
	}()

	go func() {
		for stderrScanner.Scan() {
			fmt.Fprintf(os.Stderr, "%s\n", stderrScanner.Text())
		}
	}()

	go func() {
		defer k.stdinPipe.Close()
		for stdinScanner.Scan() {
			io.WriteString(k.stdinPipe, stdinScanner.Text())
		}
	}()

	// start command
	err = k.cmd.Start()
	if err != nil {
		return err
	}

	// wait for command
	err = k.cmd.Wait()
	if err != nil {
		return err
	}

	return err
}

func New(path string) Kubectl {
	kubectl := new(kubectl)
	kubectl.wrap(path)
	return kubectl
}
