package autotest

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
)

type trigger struct {
	pattern *regexp.Regexp
	eventID EventID
}

// Process represents a process.
type Process struct {
	command string
	args    []string

	cmd    *exec.Cmd
	stdin  io.Writer
	stdout io.Reader
	stderr io.Reader

	triggers []trigger

	WorkingDir string
	Verbose    bool
}

// NewProcess creates a new process.
func NewProcess(command string, args ...string) *Process {
	return &Process{command, args, nil, nil, nil, nil, make([]trigger, 0), "", false}
}

// Start starts the process.
func (p *Process) Start() (err error) {
	if len(p.WorkingDir) > 0 {
		if err = os.Chdir(p.WorkingDir); err != nil {
			return
		}
	}
	p.cmd = exec.Command(p.command, p.args...)
	p.stdin, err = p.cmd.StdinPipe()
	if err != nil {
		return
	}
	p.stdout, err = p.cmd.StdoutPipe()
	if err != nil {
		return
	}
	p.stderr, err = p.cmd.StderrPipe()
	if err != nil {
		return
	}

	go p.listenOutput(p.stdout, os.Stdout)
	go p.listenOutput(p.stderr, os.Stderr)

	if err = p.cmd.Start(); err != nil {
		return
	}
	return nil
}

// Stop stops the process.
func (p *Process) Stop() (err error) {
	if p.cmd.Process != nil {
		err = p.cmd.Process.Kill()
		p.cmd = nil
		p.stdin = nil
		p.stdin = nil
		p.stdout = nil
		p.stderr = nil
		return
	}
	return errors.New("Not started")
}

// Wait waits until the process stops or is killed.
func (p *Process) Wait() (err error) {
	if p.cmd.Process != nil {
		return p.cmd.Wait()
	}
	return errors.New("Not started")
}

// Send sends the given command to the process' stdin.
func (p Process) Send(command string) {
	writer := bufio.NewWriter(p.stdin)
	writer.WriteString(command)
	writer.WriteString("\n")
	writer.Flush()
}

// AddTrigger registers a trigger that will be used to detect event from the process' output.
func (p *Process) AddTrigger(pattern string, event EventID) error {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("Invalid pattern: %v", err)
	}
	t := trigger{r, event}
	p.triggers = append(p.triggers, t)
	return nil
}

// listenOutput reads the given Reader and looks for trigger patterns. The function stops when the
// Reader is closed (when the process is stopped).
func (p Process) listenOutput(reader io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		if p.Verbose {
			fmt.Fprintln(output, text)
		}
		for _, t := range p.triggers {
			if matches := t.pattern.FindStringSubmatch(text); len(matches) > 0 {
				p.fireEvent(t.eventID, matches[1:])
			}
		}
	}
}

func (p Process) fireEvent(id EventID, args []string) {
	// XXX: can we make this not blocking?
	events <- event{id, args}
}
