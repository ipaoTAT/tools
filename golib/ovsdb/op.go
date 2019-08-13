package ovsdb

import "context"

type OvsDB interface {
	Init() error
	Insert(ctx context.Context, tables string, kv map[string]interface{}) error
	Query(ctx context.Context, table string, columns []string, condition map[string]interface{}) (*Rows, error)
	Set(ctx context.Context, table string, index interface{}, kv map[string]interface{}) error
	Delete(ctx context.Context, table string, index interface{}) error
}
