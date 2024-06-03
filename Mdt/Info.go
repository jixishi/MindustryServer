package Mdt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type GamemodeId int

const (
	survival GamemodeId = iota
	sandbox
	attack
	pvp
	editor
)

type Gamemode struct {
	Name string     `json:"name"`
	Id   GamemodeId `json:"id"`
}

type Info struct {
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	Status      string   `json:"status"`
	Name        string   `json:"name"`
	Maps        string   `json:"maps"`
	Players     int      `json:"players"`
	Version     int      `json:"version"`
	Wave        int      `json:"wave"`
	Vertype     string   `json:"vertype"`
	Gamemode    Gamemode `json:"gamemode"`
	Description string   `json:"description"`
	Modename    string   `json:"modename"`
	Limit       int      `json:"limit"`
	Ping        int      `json:"ping"`
}

type InfoBuffer struct {
	*bytes.Reader
}

func (g GamemodeId) Name() string {
	switch g {
	case survival:
		return "生存"
	case sandbox:
		return "沙盒"
	case attack:
		return "进攻"
	case pvp:
		return "PVP"
	case editor:
		return "编辑"
	default:
		return "Error"
	}
}

func (r *InfoBuffer) New(b []byte) {
	r.Reader = bytes.NewReader(b)
}

func (r *InfoBuffer) readString() string {
	l, _ := r.ReadByte()
	buf := make([]byte, l)
	r.Read(buf)
	return string(buf)
}

func (r *InfoBuffer) getInt() int {
	var t int32
	binary.Read(r, binary.BigEndian, &t)
	return int(t)
}

func (r *InfoBuffer) get() byte {
	b, _ := r.ReadByte()
	return b
}

func connectServer(host string) (buf InfoBuffer, ping int64, err error) {
	socket, err := net.Dial("udp", host)
	if err != nil {
		return InfoBuffer{}, -1, err
	}
	defer socket.Close()
	startTime := time.Now()
	_, err = socket.Write([]byte{0xFE, 0x01})
	if err != nil {
		return InfoBuffer{}, -1, err
	}
	data := make([]byte, 500)
	socket.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, err = socket.Read(data)
	if err != nil {
		return InfoBuffer{}, -1, err
	}
	ping = time.Since(startTime).Milliseconds()
	buf.New(data)
	return buf, ping, nil
}
func GetServerInfo(host string) (Info, error) {
	var info Info
	ip := strings.Split(host, ":")
	if len(ip) == 1 {
		info.Host = host
		info.Port = 6567
	} else {
		info.Host = ip[0]
		port, err := strconv.Atoi(ip[1])
		if err != nil {
			return info, err
		}
		info.Port = port
	}
	add := fmt.Sprintf("%s:%d", info.Host, info.Port)
	buf, ping, err := connectServer(add)
	if err != nil {
		info.Status = "Offline"
		return info, err
	}
	reg := regexp.MustCompile("(?s)\\[.*?\\]")
	info.Ping = int(ping)
	info.Status = "Online"
	info.Name = reg.ReplaceAllString(buf.readString(), "")
	info.Maps = reg.ReplaceAllString(buf.readString(), "")
	info.Players = buf.getInt()
	info.Wave = buf.getInt()
	info.Version = buf.getInt()
	info.Vertype = buf.readString()
	info.Gamemode.Id = GamemodeId(buf.get())
	info.Gamemode.Name = info.Gamemode.Id.Name()
	info.Limit = buf.getInt()
	info.Description = reg.ReplaceAllString(buf.readString(), "")
	info.Modename = buf.readString()
	return info, nil
}
