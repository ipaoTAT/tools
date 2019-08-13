package ovsdb

import (
	"context"
	"errors"
)

var ErrUnimplemented = errors.New("Unimplemented ovsdb operation")

// ovsdb handler interface define
type OvsDB interface {
	// function to do some init things
	Init() error

	// function to insert a record into ovsdb table
	Insert(ctx context.Context, table string, kv map[string]interface{}) error

	// function to query records from ovsdb table, return rows
	Query(ctx context.Context, table string, columns []string, condition map[string]interface{}) (*Rows, error)

	// function to update record fields of ovsdb table record
	Set(ctx context.Context, table string, index interface{}, kv map[string]interface{}) error

	// function to delete a recode of ovsdb table
	Delete(ctx context.Context, table string, index interface{}) error
}

// a fake ovsdb implement that doesn't implement any function
type UnimplementedOvsDB byte

func (*UnimplementedOvsDB) Init() error {
	return ErrUnimplemented
}
func (*UnimplementedOvsDB) Insert(ctx context.Context, table string, kv map[string]interface{}) error {
	return ErrUnimplemented
}
func (*UnimplementedOvsDB) Query(ctx context.Context, table string, columns []string, condition map[string]interface{}) (*Rows, error) {
	return nil, ErrUnimplemented
}
func (*UnimplementedOvsDB) Set(ctx context.Context, table string, index interface{}, kv map[string]interface{}) error {
	return ErrUnimplemented
}
func (*UnimplementedOvsDB) Delete(ctx context.Context, table string, index interface{}) error {
	return ErrUnimplemented
}
