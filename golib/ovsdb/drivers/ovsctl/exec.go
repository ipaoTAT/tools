package vsctl

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

type execUtil struct {
	logWtr io.Writer
}

func (c *execUtil) execInShell(ctx context.Context, cmd string) (string, error) {
	const SHELL_CMD = "/bin/sh"
	const SHELL_OPT = "-c"
	return c.exec(ctx, SHELL_CMD, SHELL_OPT, cmd)
}

func (c *execUtil) exec(ctx context.Context, name string, args ...string) (string, error) {
	cmdStr := name + " " + strings.Join(args, " ")
	c.logfln("exec '%s'", cmdStr)
	start := time.Now()
	bt, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	res := strings.Trim(string(bt), "\n")
	if err != nil {
		c.logfln("exec '%s' failed, error: %v, result: %s, elapsed: %+v", cmdStr, res, err, time.Since(start))
		return res, err
	}
	c.logfln("exec '%s' succeed, result: %s, elapsed: %+v", cmdStr, res, time.Since(start))
	return res, nil
}

func (c *execUtil) logfln(format string, v ...interface{}) {
	if c.logWtr == nil {
		return
	}
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(c.logWtr, format, v...)
}
