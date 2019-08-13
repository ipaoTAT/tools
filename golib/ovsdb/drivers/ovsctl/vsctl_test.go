package vsctl

import (
	"testing"
	"github.com/ipaoTAT/tools/golib/ovsdb"
	"os"
	"context"
	"time"
)

func TestVsCtl_Query(t *testing.T) {
	var db ovsdb.OvsDB = &VsCtl{LogWriter: os.Stdout}
	if err := db.Init(); err != nil {
		t.Fatal(err)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(time.Second))
	var i *ovsdb.Interface
	rows, err := db.Query(ctx, ovsdb.TABLE_INTERFACE, ovsdb.GetTableFields(i), map[string]interface{}{"external_ids": map[string]interface{}{"iface-id": "tap_metadata", "iface-status": "active"}})
	if err != nil {
		t.Fatal(err)
		return
	}
	for rows.Next() {
		i = &ovsdb.Interface{}
		if err := rows.Scan(i); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("interface: %+v", *i)
		}
	}
	var b *ovsdb.Bridge
	rows, err = db.Query(ctx, ovsdb.TABLE_BRIDGE, ovsdb.GetTableFields(b), nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	for rows.Next() {
		b = &ovsdb.Bridge{}
		if err := rows.Scan(b); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("Bridge: %+v", *b)
		}
	}
}

func TestVsCtl_Set(t *testing.T) {
	var db ovsdb.OvsDB = &VsCtl{LogWriter: os.Stdout}
	if err := db.Init(); err != nil {
		t.Fatal(err)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(time.Second))
	db.Query()
}