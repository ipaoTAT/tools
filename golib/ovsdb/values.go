package ovsdb

import (
	"encoding/json"
	"github.com/pkg/errors"
)

// abstract of ovsdb column values
type Value interface {
	// function to convert go value to ovs ovs value
	Marshal() (interface{}, error)

	// function to convert ovs value to go value and store it
	Unmarshal(interface{}) error

	// function to return go value as a interface
	Value() interface{}
}

func UnmarshalJson(bt []byte) (Value, error) {
	var i interface{}
	if err := json.Unmarshal(bt, &i); err != nil {
		return nil, errors.New("json.Unmarshal: " + err.Error())
	}
	return Unmarshal(i)
}

func Unmarshal(i interface{}) (Value, error) {
	var ret Value = nil
	switch i.(type) {
	case nil:
		ret = &NullValue{}
	case string, int, uint, float64, bool:
		ret = &BasicValue{}
	case []interface{}:
		ret = &ComplexValue{}
	default:
		return nil, BadValueTypeError(i)
	}
	return ret, ret.Unmarshal(i)
}

type NullValue struct{}

func (NullValue) Marshal() (interface{}, error) {
	return nil, nil
}

func (NullValue) Unmarshal(i interface{}) error {
	return nil
}

func (NullValue) Value() interface{} {
	return nil
}

// type basic: integer, real, boolean, string
type BasicValue struct {
	v interface{}
}

func (v BasicValue) Marshal() (interface{}, error) {
	return v.v, nil
}

func (v *BasicValue) Unmarshal(i interface{}) error {
	v.v = i
	return nil
}

func (v *BasicValue) Value() interface{} {
	return v.v
}

const (
	COMPLEX_TYPE_UUID = "uuid"
	COMPLEX_TYPE_SET  = "set"
	COMPLEX_TYPE_MAP  = "map"
)

// type complex value: for uuid, set, map
type ComplexValue struct {
	v Value
}

func (v ComplexValue) Marshal() (interface{}, error) {
	return v.v.Marshal()
}

func (v *ComplexValue) Unmarshal(i interface{}) error {
	tp, _, err := deCapComplexValue(i)
	if err != nil {
		return err
	}
	switch tp {
	case COMPLEX_TYPE_UUID:
		v.v = &UUIDValue{}
	case COMPLEX_TYPE_SET:
		v.v = &SetValue{}
	case COMPLEX_TYPE_MAP:
		v.v = &MapValue{}
	default:
		return BadComplexValueType(tp)
	}
	return v.v.Unmarshal(i)
}

func (v ComplexValue) Value() interface{} {
	return v.v.Value()
}

// type UUID
type UUIDValue struct {
	v string
}

func (v UUIDValue) Marshal() (interface{}, error) {
	return enCapComplexValue(v.v, COMPLEX_TYPE_UUID)
}

func (v *UUIDValue) Unmarshal(i interface{}) error {
	i, err := deCapSpecificTypeComplexValue(i, COMPLEX_TYPE_UUID)
	if err != nil {
		return err
	}
	s, ok := i.(string)
	if !ok {
		return BadUUIDValueType(i)
	}
	v.v = s
	return nil
}

func (v UUIDValue) Value() interface{} {
	return v.v
}

// type set
type SetValue struct {
	v []Value
}

func (s *SetValue) Marshal() (interface{}, error) {
	vals := make([]interface{}, 0, len(s.v))
	for _, v := range s.v {
		val, err := v.Marshal()
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return enCapComplexValue(vals, COMPLEX_TYPE_SET)
}

func (s *SetValue) Unmarshal(i interface{}) error {
	i, err := deCapSpecificTypeComplexValue(i, COMPLEX_TYPE_SET)
	if err != nil {
		return err
	}
	vals, ok := i.([]interface{})
	if !ok {
		return BadSetValueType(i)
	}
	s.v = make([]Value, 0, len(vals))
	for _, v := range vals {
		ov, err := Unmarshal(v)
		if err != nil {
			return err
		}
		s.v = append(s.v, ov)
	}
	return nil
}

func (s *SetValue) Value() interface{} {
	ret := make([]interface{}, 0, len(s.v))
	for _, v := range s.v {
		ret = append(ret, v.Value())
	}
	return ret
}

// type map
type MapValue struct {
	v map[string]Value
}

func (s *MapValue) Marshal() (interface{}, error) {
	vals := make([]interface{}, 0, len(s.v))
	for k, v := range s.v {
		val, err := v.Marshal()
		if err != nil {
			return nil, err
		}
		vals = append(vals, []interface{}{k, val})
	}
	return enCapComplexValue(vals, COMPLEX_TYPE_MAP)
}

func (s *MapValue) Unmarshal(i interface{}) error {
	i, err := deCapSpecificTypeComplexValue(i, COMPLEX_TYPE_MAP)
	if err != nil {
		return err
	}
	vals, ok := i.([]interface{})
	if !ok {
		return BadMapValueType(i)
	}
	s.v = make(map[string]Value, len(vals))
	for _, val := range vals {
		kv, ok := val.([]interface{})
		if !ok || len(kv) != 2 {
			return BadValueLenError()
		}
		k, ok := kv[0].(string)
		if !ok {
			return BadValueTypeError(kv[0])
		}
		v, err := Unmarshal(kv[1])
		if err != nil {
			return err
		}
		s.v[k] = v
	}
	return nil
}

func (s *MapValue) Value() interface{} {
	ret := make(map[string]interface{})
	for k, v := range s.v {
		ret[k] = v.Value()
	}
	return ret
}

// function to encap complex value with given type
func enCapComplexValue(v interface{}, tp string) (interface{}, error) {
	return []interface{}{tp, v}, nil
}

// function to decap complex value with given type
func deCapComplexValue(i interface{}) (string, interface{}, error) {
	v, ok := i.([]interface{})
	if !ok {
		return "", nil, BadValueTypeError(i)
	}
	if len(v) != 2 {
		return "", nil, BadValueLenError()
	}
	_tp, ok := v[0].(string)
	if !ok {
		return "", nil, BadValueTypeError(v[0])
	}
	return _tp, v[1], nil
}

// function to decap complex value with given type
func deCapSpecificTypeComplexValue(i interface{}, tp string) (interface{}, error) {
	_tp, val, err := deCapComplexValue(i)
	if err != nil {
		return nil, err
	}
	if _tp != tp {
		return nil, BadComplexValueType(tp)
	}
	return val, nil
}
