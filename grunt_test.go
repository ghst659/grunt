package grunt

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestRunCommand(t *testing.T) {
	cases := []struct {
		cmd []string
		err bool
	}{
		{
			cmd: []string{"/non/existent/command", "arg"},
			err: true,
		},
		{
			cmd: []string{"/bin/false"},
			err: true,
		},
		{
			cmd: []string{"/bin/true"},
			err: false,
		},
		{
			cmd: []string{"/bin/echo", "foo"},
			err: false,
		},
	}
	g := Grunt{
		Noop: false,
		Log:  log.New(os.Stderr, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	for _, c := range cases {
		t.Logf("command: %v", c.cmd)
		err := g.Run(c.cmd)
		if c.err && err == nil {
			t.Errorf("%v: missing expected error", c.cmd)
		}
		if !c.err {
			if err != nil {
				t.Errorf("%v: unexpected error: %v", c.cmd, err)
			}
		}
	}
}

func TestRunFunction(t *testing.T) {
	stateVar := "BEFORE"
	cases := []struct {
		cmd func() error
		err bool
	}{
		{
			cmd: func() error { return nil },
			err: false,
		},
		{
			cmd: func() error { return fmt.Errorf("deliberate error") },
			err: true,
		},
		{
			cmd: func() error {
				stateVar = "TRANSITION 0"
				return nil
			},
			err: false,
		},
		{
			cmd: func() error {
				stateVar = "TRANSITION 1"
				return nil
			},
			err: false,
		},
	}
	g := Grunt{
		Noop: false,
		Log:  log.New(os.Stderr, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := g.Run(c.cmd)
			if c.err && err == nil {
				t.Errorf("missing expected error")
			}
			if !c.err {
				t.Logf("State: %q", stateVar)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestNoop(t *testing.T) {
	stateVar := "BEFORE"
	lastState := stateVar
	cases := []struct {
		dry bool
		cmd func() error
	}{
		{
			dry: true,
			cmd: func() error {
				stateVar = "TRANSITION 0"
				return nil
			},
		},
		{
			dry: false,
			cmd: func() error {
				stateVar = "TRANSITION 1"
				return nil
			},
		},
	}
	logger := log.New(os.Stderr, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile)
	for i, c := range cases {
		lastState = stateVar
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			g := Grunt{
				Noop: c.dry,
				Log:  logger,
			}
			err := g.Run(c.cmd)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if c.dry && stateVar != lastState {
				t.Errorf("failed to no-op: %q", stateVar)
			}
			if !c.dry && stateVar == lastState {
				t.Errorf("failed to  execute: %q", stateVar)
			}
		})
	}
}
