package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func structprint(name string, a any) {
	bs, _ := json.Marshal(a)
	var out bytes.Buffer
	json.Indent(&out, bs, "", "\t")
	fmt.Printf("%s:%v\n", name, out.String())
}
func inServers(arr []Server, target Server) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
