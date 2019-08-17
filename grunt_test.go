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
	}
	g := Grunt{
		Noop: false,
		Log:  log.New(os.Stderr, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	for _, c := range cases {
		err := g.Run(c.cmd)
		if c.err && err == nil {
			t.Errorf("missing expected error")
		}
		if !c.err {
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}
