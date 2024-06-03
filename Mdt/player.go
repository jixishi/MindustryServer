package Mdt

import "C"
import (
	"MindustryServer/utils"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"hash/crc32"
	"net"
	"time"
)

//(ip.src == 185.248.33.51) or (ip.addr == 185.248.33.51)

// 玩家数据包
// 0000   8a ab 07 d4 c5 a0 0a bf 44 d3 a8 18 08 00 45 00   ........D.....E.
// 0010   00 72 b6 ee 40 00 80 06 00 00 c0 a8 64 d0 af b2   .r..@.......d...
// 0020   16 06 19 7d 1a 0a eb d6 7f 17 c1 89 d2 0b 50 18   ...}..........P.
// 0030   04 05 eb 95 00 00 00 48 03 00 42 01 f0 33 00 00   .......H..B..3..
// 0040   00 92 01 00 08 6f 66 66 69 63 69 61 6c 01 00 03   .....official...
// 0050   6a 78 73 01 00 05 7a 68 5f 43 4e 01 00 0c 4a 38   jxs...zh_CN...J8
// 0060   32 66 4e 51 57 56 54 34 30 3d 58 4c 28 7a 02 aa   2fNQWVT40=XL(z..
// 0070   99 7f 00 00 00 00 26 08 1a e1 00 ff a1 08 ff 00   ......&.........
//
//	public void write(Writes buffer){
//	           buffer.i(Version.build);
//	           TypeIO.writeString(buffer, versionType);
//	           TypeIO.writeString(buffer, name);
//	           TypeIO.writeString(buffer, locale);
//	           TypeIO.writeString(buffer, usid);
//
//	           byte[] b = Base64Coder.decode(uuid);
//	           buffer.b(b);
//	           CRC32 crc = new CRC32();
//	           crc.update(Base64Coder.decode(uuid), 0, b.length);
//	           buffer.l(crc.getValue());
//
//	           buffer.b(mobile ? (byte)1 : 0);
//	           buffer.i(color);
//	           buffer.b((byte)mods.size);
//	           for(int i = 0; i < mods.size; i++){
//	               TypeIO.writeString(buffer, mods.get(i));
//	           }
//	       }
type PlayerBuffer struct {
	*bytes.Buffer
}

func newPlayerBuffer() *PlayerBuffer {
	p := new(PlayerBuffer)
	return p
}
func (p *PlayerBuffer) writeString(s string) {
	p.writeInt(int32(len(s)))
	p.WriteString(s)
}

func (p *PlayerBuffer) set(b byte) {
	p.WriteByte(b)
}
func (p *PlayerBuffer) writeInt(i int32) {
	err := binary.Write(p, binary.LittleEndian, &i)
	if err != nil {
		return
	}
}
func (p *PlayerBuffer) writel(i uint32) {
	err := binary.Write(p, binary.LittleEndian, &i)
	if err != nil {
		return
	}
}

type ClientBuf struct {
	*bytes.Buffer
	tid []byte
	uid []byte
	cnt int32
}

//func newClientBuf() *ClientBuf {
//	p := new(ClientBuf)
//	return p
//}
//func (c *ClientBuf) GetClientID() {
//	c.ReadByte()
//	l, _ := c.ReadByte()
//	b := make([]byte, l)
//	c.Read(b)
//	switch b[1] {
//	case 3:
//		c.uid = b[2:]
//	case 4:
//		c.tid = b[2:]
//	}
//	return
//}
//
//func (c *ClientBuf) RegTCP(C net.Conn) {
//	b := []byte{0xfe, 0x04}
//	C.Write(append(b, c.tid...))
//}
//func (c *ClientBuf) RegUDP(C net.Conn) {
//	b := []byte{0xfe, 0x03}
//	C.Write(append(b, c.uid...))
//}

func NewPlayer(name, versiontype, locale string, version int32) PlayerBuffer {
	p := newPlayerBuffer()
	p.writeInt(146)
	if versiontype == "" {
		versiontype = "official"
	}
	if locale == "" {
		locale = "zh_CN"
	}
	p.writeString(versiontype)
	p.writeString(name)
	p.writeString(locale)
	p.writeString("xxcadawdadawd")
	b, _ := base64.StdEncoding.DecodeString(uuid.New().String())
	p.Write(b)
	p.writel(crc32.ChecksumIEEE(b))
	p.WriteByte(0)
	p.writeInt(0)
	p.WriteByte(0)
	return *p
}

func GetPlayerList(host string) {
	host = utils.HostPreProcessing(host)
	tbuf := make([]byte, 1024)
	//ubuf := make([]byte, 2048)
	udp, _ := net.Dial("udp", host)
	tcp, _ := net.Dial("tcp", host)
	defer udp.Close()
	defer tcp.Close()
	tcp.Read(tbuf)
	fmt.Printf("\n%x\n", tbuf)
	b := bytes.NewReader(tbuf)
	b.ReadByte()
	l, _ := b.ReadByte()
	ib := make([]byte, l)
	b.Read(ib)
	var uid, tid []byte
	switch ib[1] {
	case 3:
		uid = ib[2:]
		fmt.Printf("\nuid:%x\n", uid)
	case 4:
		tid = ib[2:]
		fmt.Printf("\ntid:%x\n", tid)
	}
	ib = []byte{0xfe, 0x03}
	udp.Write(append(ib, tid...))
	time.Sleep(100 * time.Millisecond)
	udp.Write(append(ib, tid...))
	tcp.Read(tbuf)
	fmt.Printf("\n%x\n", tbuf)
	b = bytes.NewReader(tbuf)
	b.ReadByte()
	l, _ = b.ReadByte()
	ib = make([]byte, l)
	b.Read(ib)
	switch ib[1] {
	case 3:
		uid = ib[2:]
		fmt.Printf("\nuid:%x\n", uid)
	case 4:
		tid = ib[2:]
		fmt.Printf("\ntid:%x\n", tid)
	}
	//c := newClientBuf()
	//c.Buffer.Write(tbuf)
	//c.GetClientID()
	//c.RegTCP(udp)
	//c.RegTCP(udp)
	//c.GetClientID()
	p := NewPlayer("jxs", "", "", 146)
	tcp.Write(p.Bytes())
	tcp.Read(tbuf)
	fmt.Println(string(tbuf))
	tcp.Read(tbuf)
	fmt.Println(string(tbuf))
}
