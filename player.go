package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"github.com/google/uuid"
	"hash/crc32"
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
	err := binary.Write(p, binary.BigEndian, &i)
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
func NewPlayer(name, versiontype, locale string, version int) {
	p := newPlayerBuffer()
	p.writeInt(int32(version))
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
}
