package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type RecvMessage struct {
	DeviceID string            `json:"deviceid"`
	Topic    string            `json:"topic"`
	Payload  string            `json:"payload"` // could be []byte later
	Config   map[string]string `json:"config"`
}

type Command struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser
}

func NewCommand(path string, args []string) *Command {
	return &Command{
		cmd:   exec.Command(path, args...),
		stdin: nil,
	}
}

func (c *Command) Start() error {
	var err error
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	c.cmd.Stdin = nil
	c.stdin, err = c.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("Failed to open a Stdin pipe: %v", err)
	}
	if err := c.cmd.Start(); err != nil {
		return err
	}
	// TODO: Create a command ready protocol that we can wait for
	time.Sleep(time.Second * 2)
	return nil
}

func (c *Command) Exited() bool {
	return c.cmd.ProcessState.Exited()
}

func (c *Command) Stop() error {
	return c.cmd.Process.Signal(os.Interrupt)
}

func (c *Command) Recv(devid, topic string, payload []byte, config map[string]string) error {
	var rm = RecvMessage{
		DeviceID: devid,
		Topic:    topic,
		Payload:  string(payload),
		Config:   config,
	}
	// if c.Exited() {
	// 	return fmt.Errorf("Child command has exited unexpectedly")
	// }

	out, err := json.Marshal(&rm)
	// fmt.Fprintf(c.stdin, "%s\n", string(payload))
	if err != nil {
		return err
	}
	if _, err := c.stdin.Write(out); err != nil {
		return err
	}
	if _, err := c.stdin.Write([]byte("\n")); err != nil {
		return err
	}
	return nil
}
