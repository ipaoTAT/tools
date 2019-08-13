package ovsdb

import (
	"bytes"
	"testing"
)

func TestNewRows(t *testing.T) {
	invalidInputs := []string{
		`{`,
		`{}`,
		`{"headings":[]}`,
		`{"data":[]}`,
		`{"headings":[],"data":{}}`,
	}
	validInputs := []string{
		`{"headings":[],"data":[]}`,
		`{"headings":["column1"],"data":[]}`,
	}
	for _, s := range invalidInputs {
		_, err := NewRows(bytes.NewBufferString(s))
		if err == nil {
			t.Fatal("error should not be nil")
		}
		t.Log(err)
	}
	for _, s := range validInputs {
		_, err := NewRows(bytes.NewBufferString(s))
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestRows_Next(t *testing.T) {
	eles := map[string]int{
		`{"headings":[],"data":[]}`:                              0,
		`{"headings":["column1"],"data":[]}`:                     0,
		`{"headings":["column1"],"data":[[]]}`:                   1,
		`{"headings":["column1"],"data":[[1],[2,3]]}`:            2,
		`{"headings":["column1","column2"],"data":[["xsa", 2]]}`: 1,
	}
	for s, v := range eles {
		rows, err := NewRows(bytes.NewBufferString(s))
		if err != nil {
			t.Fatal(err)
			continue
		}
		i := 0
		for ; rows.Next(); i++ {
		}
		if i != v {
			t.Fatal("expect ", v, " Get ", i)
		} else {
			t.Log("expect ", v, " Get ", i)
		}
		if rows.Next() {
			t.Fatal("rows.Next() == true")
		}
	}
}

func TestRows_Scan(t *testing.T) {
	str := `{"headings":["column1","Column2","column3", "Column4"],"data":[["xsa", 2, true, ["set", [1,2]]]]}`
	i := &struct {
		Column1 string `ovsdb:"column1"`
		Column2 int
		Column3 bool `ovsdb:"column3"`
		Column4 []interface{}
		Column5 []interface{} `ovsdb:"-"`
	}{}
	rows, err := NewRows(bytes.NewBufferString(str))
	if err != nil {
		t.Fatal(err)
	} else {
		rows.Next()
		err := rows.Scan(&i)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log("val=", *i)
		}
	}
}
