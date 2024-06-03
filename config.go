package main

import (
	"encoding/json"
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

func LoadConfig(path string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&Cfg)
	return err
}
func SaveConfig(path string) error {
	jsonFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	encode := json.NewEncoder(jsonFile)
	encode.SetIndent("", "\t")
	err = encode.Encode(Cfg)
	return nil
}
