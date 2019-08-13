package ovsdb

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrUnexpected   = errors.New("ovsdb: there is a bug in lib")
	ErrRowsBadIndex = errors.New("ovsdb.rows: index out of range")
	ErrRowsBadData  = errors.New("ovsdb.rows: bad raw ovs json data")
)

func BadValueTypeError(v interface{}) error {
	return fmt.Errorf("ovsdb.value: bad value type %s, value: %v", reflect.TypeOf(v).String(), v)
}

func BadValueLenError() error {
	return errors.New("ovsdb.value: bad value length, need 2")
}

func BadComplexValueType(tp string) error {
	return errors.New("ovsdb.value: bad complex value type " + tp)
}

func BadUUIDValueType(v interface{}) error {
	return fmt.Errorf("ovsdb.value: bad uuid value type %s, expect string", reflect.TypeOf(v).String())
}

func BadMapValueType(v interface{}) error {
	return fmt.Errorf("ovsdb.value: bad map value type %s, expect []interface{}", reflect.TypeOf(v).String())
}

func BadSetValueType(v interface{}) error {
	return fmt.Errorf("ovsdb.value: bad set value type %s, expect []interface{}", reflect.TypeOf(v).String())
}
