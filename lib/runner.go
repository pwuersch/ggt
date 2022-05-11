package lib

import (
	"os"
	"os/exec"
	"strings"

	"github.com/pwuersch/ggt/lib/globals"
)

type Runner struct {
	command *exec.Cmd
}

func (r *Runner) WithDir(dir string) *Runner {
	r.command.Dir = dir
	return r
}

func (r *Runner) Run() error {
	r.command.Stderr = os.Stdout
	r.command.Stdout = os.Stdout
	r.command.Stdin = os.Stdin
	return r.command.Run()
}

func NewRunner(name string, args ...string) *Runner {
	runner := &Runner{}
	runner.command = exec.Command(name, args...)

	globals.Debug(strings.Join(append([]string{name}, args...), " "))

	return runner
}
