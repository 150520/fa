package adb

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
)

func swapUint32(n uint32) uint32 {
	var i uint32
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, n)
	binary.Read(buf, binary.BigEndian, &i)
	return i
}

// Packet is a meta for adb connect
type Packet struct {
	Command string
	Arg0    uint32
	Arg1    uint32
	Body    []byte
}

func (pkt Packet) magic() []byte {
	return xorBytes([]byte(pkt.Command), []byte{0xff, 0xff, 0xff, 0xff})
}

func (pkt Packet) checksum() uint32 {
	return checksum(pkt.Body)
}

func (pkt Packet) swapu32(n uint32) uint32 {
	var i uint32
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, n)
	binary.Read(buf, binary.BigEndian, &i)
	return i
}

func (pkt Packet) EncodeToBytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 24+len(pkt.Body)))
	if len(pkt.Command) != 4 {
		panic("Invalid command " + strconv.Quote(pkt.Command))
	}
	binary.Write(buf, binary.LittleEndian, []byte(pkt.Command))
	binary.Write(buf, binary.LittleEndian, pkt.Arg0)
	binary.Write(buf, binary.LittleEndian, pkt.Arg1)
	binary.Write(buf, binary.LittleEndian, uint32(len(pkt.Body)))
	binary.Write(buf, binary.LittleEndian, pkt.checksum())
	binary.Write(buf, binary.LittleEndian, pkt.magic())
	buf.Write(pkt.Body)
	return buf.Bytes()
}

func (pkt Packet) WriteTo(wr io.Writer) (n int, err error) {
	return wr.Write(pkt.EncodeToBytes())
}

func (pkt Packet) DumpToStdout() {
	fmt.Printf("cmd:%s arg0:%d arg1:%d\n", pkt.Command, pkt.Arg0, pkt.Arg1)
	dumper := hex.Dumper(os.Stdout)
	dumper.Write(pkt.Body)
	dumper.Close()
}