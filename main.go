package main

import (
	"MindustryServer/Mdt"
	"MindustryServer/config"
	"MindustryServer/utils"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var Apis ApiTemplate
var ServersPlayers map[string][]int

func UpdateInfo() {
	for i := 0; i < len(config.Cfg.Servers); i++ {
		info, _ := Mdt.GetServerInfo(config.Cfg.Servers[i].Host)
		config.Cfg.Servers[i].Name = info.Name
		if len(ServersPlayers[config.Cfg.Servers[i].Host]) >= 20 {
			ServersPlayers[config.Cfg.Servers[i].Host] = ServersPlayers[config.Cfg.Servers[i].Host][1:]
		}
		ServersPlayers[config.Cfg.Servers[i].Host] = append(ServersPlayers[config.Cfg.Servers[i].Host], info.Players)
	}
}

func main() {
	flag.Parse()
	if utils.FileExist(config.ConfigPath) {
		config.LoadConfig()
	} else {
		config.Cfg.Port = 8080
		config.Cfg.Servers = []config.Server{{Host: "cn.mindustry.top"}, {Host: "mdtleague.top"}}
		config.SaveConfig()
		fmt.Printf("配置文件创建完成:\t%s\n", config.ConfigPath)
	}
	if newhost != "" {
		config.UpdateConfig(config.Server{Host: newhost})
	}
	ServersPlayers = make(map[string][]int)
	var ticker *time.Ticker
	done := make(chan bool)
	if chost != "" {
		fmt.Printf("正在请求:\t%s\n", chost)
		//info, _ := Mdt.GetServerInfo(chost)
		//utils.Structprint("Info", info)
		Mdt.GetPlayerList(chost)
	} else {
		ticker = time.NewTicker(time.Duration(config.Cfg.Interval) * time.Second)
		go func() {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					config.LoadConfig()
					UpdateInfo()
					config.SaveConfig()
				}
			}
		}()
	}
	Apis, _ = initApiTemplate(config.Cfg.ThemesDir)
	http.HandleFunc(Mindustry.Url, GetMindustryInfo)
	fmt.Printf("ListenAndServe On Port %v \n", config.Cfg.Port)
	if err := http.ListenAndServe(":"+strconv.Itoa(config.Cfg.Port), nil); err != nil {
		fmt.Println("ServeErr:", err)
	}
	ticker.Stop()
	done <- true
}
