package ovsdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

const (
	TABLE_INTERFACE = "interface"
	TABLE_BRIDGE    = "bridge"
)

const OVSDB_TAG_COLUMN = "ovsdb"
const OVSDB_TAG_SKIP = "-"

// bridge table record
type Bridge struct {
	Uuid         string                 `ovsdb:"_uuid"`
	Name         string                 `ovsdb:"name"`
	DatapathType string                 `ovsdb:"datapath_type"`
	ExternalIds  map[string]interface{} `ovsdb:"external_ids"`
	FailMode     string                 `ovsdb:"fail_mode"`
	OtherConfig  map[string]interface{} `ovsdb:"other_config"`
	Protocols    []interface{}          `ovsdb:"protocols"`
}

// interface table record
type Interface struct {
	Uuid        string                 `ovsdb:"_uuid"`
	Name        string                 `ovsdb:"name"`
	ExternalIds map[string]interface{} `ovsdb:"external_ids"`
	Ofport      float64                `ovsdb:"ofport"`
	Type        string                 `ovsdb:"type"`
	Options     map[string]interface{} `ovsdb:"options"`
}

// function to get columns of ovsdb table
func GetTableFields(i interface{}) []string {
	tp := reflect.TypeOf(i).Elem()
	ret := make([]string, 0, tp.NumField())
	for i := 0; i < tp.NumField(); i++ {
		tag, ok := tp.Field(i).Tag.Lookup(OVSDB_TAG_COLUMN)
		if ok {
			if tag != OVSDB_TAG_SKIP {
				ret = append(ret, tag)
			}
		} else {
			ret = append(ret, tp.Field(i).Name)
		}
	}
	return ret
}

// function to store columns value into dst
func setTableFields(dst interface{}, kv map[string]interface{}) error {
	// prepare tag-value map
	tp := reflect.TypeOf(dst)
	if tp.Kind() != reflect.Ptr || tp.Elem().Kind() != reflect.Struct {
		// if destination value is not struct pointer, use json convert value directly
		return convertJsonValue(dst, kv)
	}
	tp = tp.Elem()
	val := reflect.ValueOf(dst).Elem()
	vMap := make(map[string]*reflect.Value, tp.NumField())
	for i := 0; i < tp.NumField(); i++ {
		_v := val.Field(i)
		tag, ok := tp.Field(i).Tag.Lookup(OVSDB_TAG_COLUMN)
		if !ok {
			// if not tag, use field name as key
			tag = tp.Field(i).Name
		}
		if tag != OVSDB_TAG_SKIP {
			vMap[tag] = &_v
		}
	}
	// store field values
	for k, v := range kv {
		ptr, ok := vMap[k]
		if !ok || ptr == nil {
			continue
		}
		dstVal := *ptr
		if !dstVal.CanAddr() {
			return fmt.Errorf("ovsdb.tables: cannot get addr for field '%s' type %s", k, dstVal.Type().String())
		}
		// convert value 'v' into dstVal
		if err := convertJsonValue(dstVal.Addr().Interface(), v); err != nil {
			return err
		}
	}
	return nil
}

// function to convert src value into dst use json
func convertJsonValue(dst interface{}, src interface{}) error {
	bt, err := json.Marshal(src)
	if err != nil {
		return errors.New("ovsdb.tables: json.Marshal failed, " + err.Error())
	}
	if err := json.Unmarshal(bt, dst); err != nil {
		errors.New("ovsdb.tables: json.Unmarshal failed, " + err.Error())
	}
	return nil
}
