// Package grunt implements a subprocess runner with a no-op flag.
package grunt

import (
	"log"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
)

type Grunt struct {
	Noop bool
	Log  *log.Logger
}

func (g *Grunt) Run(r interface{}) (e error) {
	switch r := r.(type) {
	case []string:
		e = g.doCmd(r)
	case func() error:
		e = g.doFun(r)
	}
	return
}

func (g *Grunt) doFun(f func() error) (e error) {
	prefix := "FUNCALL"
	if g.Noop {
		prefix = "NOOP"
	}
	g.Log.Printf("%s: %s", prefix, getFuncName(f))
	if !g.Noop {
		if e = f(); e != nil {
			g.Log.Print("ERROR: %v", e)
		}
	}
	return
}

func (g *Grunt) doCmd(cmdvec []string) (e error) {
	prefix := "SUBPROCESS"
	if g.Noop {
		prefix = "NOOP"
	}
	g.Log.Printf("%s: %s", prefix, strings.Join(cmdvec, " "))
	if !g.Noop {
		cmd := exec.Cmd{
			Path:   cmdvec[0],
			Args:   cmdvec,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
		if e = cmd.Run(); e != nil {
			g.Log.Print("ERROR: %v", e)
		}
	}
	return
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
