package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func Structprint(name string, a any) {
	bs, _ := json.Marshal(a)
	var out bytes.Buffer
	json.Indent(&out, bs, "", "\t")
	fmt.Printf("%s:%v\n", name, out.String())
}

func HostPreProcessing(host string) string {
	ip := strings.Split(host, ":")
	var port int
	if len(ip) == 1 {
		host = host
		port = 6567
	} else {
		host = ip[0]
		port, _ = strconv.Atoi(ip[1])
	}
	return fmt.Sprintf("%s:%d", host, port)
}
