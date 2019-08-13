package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ipaoTAT/tools/golib/ovsdb"
	"github.com/ipaoTAT/tools/golib/ovsdb/drivers/vsctl"
)

func main() {
	var db ovsdb.OvsDB = &ovsctl.VsCtl{LogWriter: os.Stdout}
	if err := db.Init(); err != nil {
		fmt.Println(err)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(time.Second))
	var i *ovsdb.Interface
	rows, err := db.Query(ctx, ovsdb.TABLE_INTERFACE, ovsdb.GetTableFields(i), map[string]interface{}{"external_ids": map[string]interface{}{"iface-id": "tap_metadata", "iface-status": "active"}})
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		i = &ovsdb.Interface{}
		if err := rows.Scan(i); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("interface: %+v\n", *i)
		}
	}
	var b *ovsdb.Bridge
	rows, err = db.Query(ctx, ovsdb.TABLE_BRIDGE, ovsdb.GetTableFields(b), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		b = &ovsdb.Bridge{}
		if err := rows.Scan(b); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Bridge: %+v\n", *i)
		}
	}
}
