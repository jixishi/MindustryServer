package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Servers   []Server `json:"servers"`
	Port      int      `json:"port"`
	ThemesDir string   `json:"themesDir"`
	Interval  int      `json:"interval"`
}
type Server struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

var Cfg Config
var ConfigPath = "./config.json"

func InServers(arr []Server, target Server) bool {
	for _, v := range arr {
		if v.Host == target.Host {
			return true
		}
	}
	return false
}
func LoadConfig() error {
	jsonFile, err := os.Open(ConfigPath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&Cfg)
	return err
}
func SaveConfig() error {
	jsonFile, err := os.Create(ConfigPath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	encode := json.NewEncoder(jsonFile)
	encode.SetIndent("", "\t")
	err = encode.Encode(Cfg)
	return nil
}

func UpdateConfig(server Server) {
	if !InServers(Cfg.Servers, server) {
		Cfg.Servers = append(Cfg.Servers, server)
		SaveConfig()
		fmt.Printf("配置文件更新完成:\t%v\n", server)
	}
}
