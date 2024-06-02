package ServerInfo

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
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

type ServerInfo struct {
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

func connectServer(host string) (buf InfoBuffer, err error) {
	socket, err := net.Dial("udp", host)
	if err != nil {
		return InfoBuffer{}, err
	}
	defer socket.Close()
	_, err = socket.Write([]byte{0xFE, 0x01})
	if err != nil {
		return InfoBuffer{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	workDone := make(chan struct{}, 1)
	data := make([]byte, 500)
	go func() {
		_, err = socket.Read(data)
		workDone <- struct{}{}
	}()
	select {
	case <-workDone:
		buf.New(data)
	case <-ctx.Done():
		return InfoBuffer{}, fmt.Errorf("超时:%v", err)
	}
	return buf, nil
}
func ping(host string) (int64, error) {
	startTime := time.Now()
	conn, err := net.DialTimeout("tcp", host, time.Second*2)
	if err != nil {
		return 0, err
	}
	conn.Close()
	conn.RemoteAddr()
	elapsedTime := time.Since(startTime).Milliseconds()
	return elapsedTime, nil
}
func GetServerInfo(host string) (ServerInfo, error) {
	var info ServerInfo
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
	d, err := ping(add)
	if err != nil {
		info.Status = "Offline"
		return info, err
	}
	info.Ping = int(d)
	buf, err := connectServer(add)
	if err != nil {
		info.Status = "Offline"
		return info, err
	}
	info.Status = "Online"
	info.Name = buf.readString()
	info.Maps = buf.readString()
	info.Players = buf.getInt()
	info.Wave = buf.getInt()
	info.Version = buf.getInt()
	info.Vertype = buf.readString()
	info.Gamemode.Id = GamemodeId(buf.get())
	info.Gamemode.Name = info.Gamemode.Id.Name()
	info.Limit = buf.getInt()
	info.Description = buf.readString()
	info.Modename = buf.readString()
	return info, nil
}
