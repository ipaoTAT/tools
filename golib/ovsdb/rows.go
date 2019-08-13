package ovsdb

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"sync"
	"fmt"
)

type Rows struct {
	sync.Mutex
	h    []string
	data []interface{}
	cur  int
}

func NewRows(r io.Reader) (*Rows, error) {
	bt, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	val := &struct {
		Headers []string      `json:"headings"`
		Data    []interface{} `json:"data"`
	}{}
	err = json.Unmarshal(bt, val)
	if err != nil {
		return nil, fmt.Errorf("ovsdb.rows: %s", err.Error())
	}
	if val.Headers == nil || val.Data == nil {
		return nil, ErrRowsBadData
	}
	return &Rows{h: val.Headers, data: val.Data, cur: -1}, nil
}

func (r *Rows) Next() bool {
	r.Lock()
	defer r.Unlock()
	if r.cur+1 >= len(r.data) {
		return false
	}
	r.cur++
	return true
}

func (r *Rows) currentRow() (*Row, error) {
	r.Lock()
	defer r.Unlock()
	if r.cur < 0 || r.cur >= len(r.data) {
		return nil, ErrRowsBadIndex
	}
	return &Row{r.h, r.data[r.cur]}, nil
}

func (r *Rows) Scan(i interface{}) error {
	row, err := r.currentRow()
	if err != nil {
		return err
	}
	return row.Scan(i)
}

func (r *Rows) Err() error {
	return nil
}

type Row struct {
	h []string
	d interface{}
}

func (r *Row) Scan(i interface{}) error {
	vals, ok := r.d.([]interface{})
	if !ok {
		return errors.New("ovsdb.rows: bad value type")
	}
	if len(vals) != len(r.h) {
		return errors.New("ovsdb.rows: bad value length")
	}
	kv := make(map[string]interface{}, len(r.h))
	for i, v := range vals {
		n, err := Unmarshal(v)
		if err != nil {
			return err
		}
		kv[r.h[i]] = n.Value()
	}
	return setTableFields(i, kv)
}
