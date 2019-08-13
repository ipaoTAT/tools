package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	m := []string{"abc", "def"}
	bt, _ := json.MarshalIndent(m, "", "")
	fmt.Println(string(bt))
}