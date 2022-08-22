package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type Runner interface {
	Run() (string, error)
}

type Step struct {
	Name       string
	Bin        string
	Args       []string
	successful bool
	log        string
	wd         string
}

// TODO
// type CloneStep struct {
// 	Step
// }
// type PullStep struct {
// 	Step
// }
// type DeployStep struct {
// 	Step
// }
// type NotifyStep struct {
// 	Step
// }

// Run executes the step and returns corresponding output. Satisfies the Runner
// interface.
func (s *Step) Run() (string, error) {
	cmdStr := s.Bin + " " + strings.Join(s.Args, " ")
	cmd := exec.Command(s.Bin, s.Args...)
	cmd.Dir = s.wd

	out, err := cmd.CombinedOutput()
	if err != nil {
		return cmdStr + "\n" + string(out), &StepError{
			name:    s.Name,
			message: "step failed to execute",
			err:     err,
		}
	}
	return cmdStr + "\n" + string(out), nil
}

type StepError struct {
	name    string
	message string
	err     error
}

func (s *StepError) Error() string {
	return fmt.Sprintf("%s: %q", s.message, s.err)
}
