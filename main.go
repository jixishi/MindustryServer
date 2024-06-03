package main

import (
	"MindustryServer/Mdt"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var Apis ApiTemplate
var ServersPlayers map[string][]int

func UpdateInfo() {
	for i := 0; i < len(Cfg.Servers); i++ {
		info, _ := Mdt.GetServerInfo(Cfg.Servers[i].Host)
		Cfg.Servers[i].Name = info.Name
		if len(ServersPlayers[Cfg.Servers[i].Host]) >= 20 {
			ServersPlayers[Cfg.Servers[i].Host] = ServersPlayers[Cfg.Servers[i].Host][1:]
		}
		ServersPlayers[Cfg.Servers[i].Host] = append(ServersPlayers[Cfg.Servers[i].Host], info.Players)
	}
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
	ServersPlayers = make(map[string][]int)
	var ticker *time.Ticker
	done := make(chan bool)
	if chost != "" {
		fmt.Printf("正在请求:\t%s\n", chost)
		info, _ := Mdt.GetServerInfo(chost)
		structprint("Info", info)
	} else {
		ticker = time.NewTicker(time.Duration(Cfg.Interval) * time.Second)
		go func() {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					LoadConfig(ConfigPath)
					UpdateInfo()
					SaveConfig(ConfigPath)
				}
			}
		}()
	}
	Apis, _ = initApiTemplate(Cfg.ThemesDir)
	http.HandleFunc(Mindustry.Url, GetMindustryInfo)
	fmt.Printf("ListenAndServe On Port %v \n", Cfg.Port)
	if err := http.ListenAndServe(":"+strconv.Itoa(Cfg.Port), nil); err != nil {
		fmt.Println("ServeErr:", err)
	}
	ticker.Stop()
	done <- true
}
