package vsctl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ipaoTAT/tools/golib/ovsdb"
	"io"
	"strings"
)

var ErrUnsupport = errors.New("unsupport operation")

type VsCtl struct {
	LogWriter io.Writer
	execUtil
}

func (c *VsCtl) Init() error {
	c.execUtil.logWtr = c.LogWriter
	return nil
}

func (c *VsCtl) Query(ctx context.Context, table string, columns []string, cond map[string]interface{}) (*ovsdb.Rows, error) {
	colStr := strings.Join(columns, ",")
	condStr := strings.Join(c.convertCondition(cond), " ")
	cmd := fmt.Sprintf("ovs-vsctl -f json --column=%s find %s %s", colStr, table, condStr)
	res, err := c.execInShell(ctx, cmd)
	if err != nil {
		return nil, errors.New(err.Error() + ":" + res)
	}
	rows, err := ovsdb.NewRows(bytes.NewReader([]byte(res)))
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *VsCtl) Set(ctx context.Context, table string, index interface{}, kv map[string]interface{}) error {
	kvStrs := c.convertCondition(kv)
	cmds := []string{}
	preCmd := fmt.Sprintf("-- set %s %v ", table, index)
	for _, kvStr := range kvStrs {
		cmds = append(cmds, preCmd+kvStr)
	}
	cmd := "ovs-vsctl -- " + strings.Join(cmds, " ")
	res, err := c.execInShell(ctx, cmd)
	if err != nil {
		return errors.New(err.Error() + ":" + res)
	}
	return nil
}

func (c *VsCtl) Delete(ctx context.Context, table string, index interface{}) error {
	return ErrUnsupport
}

func (c *VsCtl) Insert(ctx context.Context, tables string, kv map[string]interface{}) error {
	return ErrUnsupport
}

func (c *VsCtl) convertCondition(cond map[string]interface{}, superKey ...string) []string {
	var ret []string
	for k, v := range cond {
		subCond, ok := v.(map[string]interface{})
		if ok {
			ret = append(ret, c.convertCondition(subCond, append(superKey, k)...)...)
		} else {
			ret = append(ret, fmt.Sprintf("%s=%v ", strings.Join(append(superKey, k), ":"), v))
		}
	}
	return ret
}
