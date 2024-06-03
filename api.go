package main

import (
	"MindustryServer/Mdt"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Api struct {
	Mode           string
	Url            string
	Args           []Args
	SampleResponse string
}
type Args struct {
	Name        string
	Required    bool
	Description string
}

func ApiViewBuild(api Api, err interface{}, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"Api":    api,
		"Err":    err,
		"Config": Cfg,
		"Data":   data,
	}
}

func (t TemplatePointer) ErrorInfoView(w io.Writer, a Api, err interface{}) {
	t.WriteData(w, ApiViewBuild(a, err, ""))
}
func (t TemplatePointer) DataInfoView(w io.Writer, a Api, data interface{}) {
	t.WriteData(w, ApiViewBuild(a, "", data))
}

var Mindustry = Api{
	Mode: "GET",
	Url:  "/api/mdt",
	Args: []Args{{Name: "host", Required: true, Description: "服务器地址"},
		{Name: "mode", Required: false, Description: "输出模式(默认json,img,html,player)"}},
	SampleResponse: "{\n    \"host\": \"p4.simpfun.cn\",\n    \"port\": 8952,\n    \"status\": \"Online\",\n    \"name\": \"[#00ff00]镜影若滴の低配备用服\",\n    \"maps\": \"未知\",\n    \"players\": 0,\n    \"version\": 146,\n    \"wave\": 1,\n    \"vertype\": \"official\",\n    \"gamemode\": {\n        \"name\": \"生存\",\n        \"id\": 0\n    },\n    \"description\": \"[#00ff00]低配但稳定的备用服，[#ffaaff]欢迎加入QQ群：726525226\",\n    \"modename\": \"\",\n    \"limit\": 0,\n    \"ping\": 40\n}",
}

func GetMindustryInfo(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		info Mdt.Info
	)
	if err = r.ParseForm(); err != nil {
		http.Error(w, "参数解析错误", http.StatusInternalServerError)
	}
	host := r.Form.Get("host")
	mode := r.Form.Get("mode")
	if host != "" {
		info, err = Mdt.GetServerInfo(host)
		if !inServers(Cfg.Servers, Server{Host: host, Name: info.Name}) {
			Cfg.Servers = append(Cfg.Servers, Server{Host: host, Name: info.Name})
			SaveConfig(ConfigPath)
			fmt.Printf("配置文件更新完成:\t%s\n", ConfigPath)
		}
		if mode == "" {
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			err = json.NewEncoder(w).Encode(info)
		}
		if mode == "img" {

		}
		if mode == "player" {
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			err = json.NewEncoder(w).Encode(ServersPlayers[host])
		}
		if mode == "html" {
			Apis.Mindustry.DataInfoView(w, Mindustry, map[string]interface{}{
				"Info":     info,
				"UpUrl":    Mindustry.Url + "?host=" + host,
				"Players":  ServersPlayers[host],
				"Interval": Cfg.Interval,
			})
		}
	} else {
		Apis.Info.ErrorInfoView(w, Mindustry, "host为空")
		return
	}
	if err != nil {
		http.Error(w, "无法解析服务器数据", http.StatusInternalServerError)
	}
}
