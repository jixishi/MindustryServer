package main

import (
	"MindustryServer/ServerInfo"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type ptrVal struct {
	*string
	*int
	*bool
	*float64
	*float32
}
type Val struct {
	string
	int
	bool
	float64
	float32
}
type Flag struct {
	v    ptrVal
	sStr string
	lStr string
	dv   Val
	help string
}

var Apis ApiTemplate
var ConfigPath = "./config.json"

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func flagInit(f *Flag) {
	if f.v.string != nil {
		flag.StringVar(f.v.string, f.sStr, f.dv.string, "")
		flag.StringVar(f.v.string, f.lStr, f.dv.string, f.help)
	}
	if f.v.bool != nil {
		flag.BoolVar(f.v.bool, f.sStr, f.dv.bool, "")
		flag.BoolVar(f.v.bool, f.lStr, f.dv.bool, f.help)
	}
	if f.v.int != nil {
		flag.IntVar(f.v.int, f.sStr, f.dv.int, "")
		flag.IntVar(f.v.int, f.lStr, f.dv.int, f.help)
	}
}

var newhost string
var host string

func init() {
	flagInit(&Flag{ptrVal{int: &Cfg.Port}, "p", "port", Val{int: 8080}, "Api接口地址"})
	flagInit(&Flag{ptrVal{string: &newhost}, "a", "add", Val{string: ""}, "添加服务器地址"})
	flagInit(&Flag{ptrVal{string: &host}, "h", "host", Val{string: ""}, "服务器地址\tip:port"})
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
func main() {
	flag.Parse()
	if FileExist(ConfigPath) {
		LoadConfig(ConfigPath)
	} else {
		Cfg.Port = 8080
		Cfg.Servers = []Server{{Host: "cn.mindustry.top"}, {Host: "mdtleague.top"}}
		SaveConfig(ConfigPath)
		fmt.Printf("配置文件创建完成:\t%s\n", ConfigPath)
	}
	if newhost != "" {
		if !inServers(Cfg.Servers, Server{Host: newhost}) {
			Cfg.Servers = append(Cfg.Servers, Server{Host: newhost})
			SaveConfig(ConfigPath)
			fmt.Printf("配置文件更新完成:\t%s\n", ConfigPath)
		}
	}
	if host != "" {
		fmt.Printf("正在请求:\t%s\n", host)
		info, _ := ServerInfo.GetServerInfo(host)
		structprint("Info", info)
	} else {
		for i := 0; i < len(Cfg.Servers); i++ {
			info, _ := ServerInfo.GetServerInfo(Cfg.Servers[i].Host)
			Cfg.Servers[i].Name = info.Name
			//structprint(Cfg.Servers[i].Host, info)
		}
	}
	SaveConfig(ConfigPath)
	Apis, _ = initApiTemplate(Cfg.ThemesDir)
	http.HandleFunc(Mindustry.Url, GetMindustryInfo)
	fmt.Printf("ListenAndServe On Port %v \n", Cfg.Port)
	if err := http.ListenAndServe(":"+strconv.Itoa(Cfg.Port), nil); err != nil {
		fmt.Println("ServeErr:", err)
	}
}
