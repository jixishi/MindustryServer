package main

import (
	"MindustryServer/ServerInfo"
	"encoding/json"
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

func ApiViewBuild(api Api, err interface{}) map[string]interface{} {
	return map[string]interface{}{
		"Api":    api,
		"Err":    err,
		"Config": Cfg,
	}
}

func (a Api) ErrorInfoView(w io.Writer, err interface{}) {
	Apis.Info.WriteData(w, ApiViewBuild(a, err))
}

var Mindustry = Api{
	Mode:           "GET",
	Url:            "/api/mdt",
	Args:           []Args{{Name: "host", Required: true, Description: "服务器地址"}},
	SampleResponse: "{\n    \"host\": \"p4.simpfun.cn\",\n    \"port\": 8952,\n    \"status\": \"Online\",\n    \"name\": \"[#00ff00]镜影若滴の低配备用服\",\n    \"maps\": \"未知\",\n    \"players\": 0,\n    \"version\": 146,\n    \"wave\": 1,\n    \"vertype\": \"official\",\n    \"gamemode\": {\n        \"name\": \"生存\",\n        \"id\": 0\n    },\n    \"description\": \"[#00ff00]低配但稳定的备用服，[#ffaaff]欢迎加入QQ群：726525226\",\n    \"modename\": \"\",\n    \"limit\": 0,\n    \"ping\": 40\n}",
}

func GetMindustryInfo(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		info ServerInfo.ServerInfo
	)
	if err = r.ParseForm(); err != nil {
		http.Error(w, "参数解析错误", http.StatusInternalServerError)
	}
	host := r.Form.Get("host")
	if host != "" {
		info, err = ServerInfo.GetServerInfo(host)
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//}
		// 如果需要返回JSON数据给客户端，可以使用以下代码
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		err = json.NewEncoder(w).Encode(info)
	} else {
		Mindustry.ErrorInfoView(w, "host为空")
		return
	}
	if err != nil {
		http.Error(w, "无法解析服务器数据", http.StatusInternalServerError)
	}
}
