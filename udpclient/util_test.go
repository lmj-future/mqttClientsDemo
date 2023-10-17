package udpclient

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/dyrkin/bin"
)

func TestXxx(t *testing.T) {
	 
	//t1 := []byte{72,101,105,109,97,110}
	//t1 := []byte{83,109,97,114,116,80,108,117,103}
	t1 := "0258"
	b, _ := hex.DecodeString(t1)
	// for i := 0; i < len(t1); i ++ {
	// 	chu, mod := t1[i] / 10, t1[i] % 10
	// 	t1[i] = chu * 16 + mod
	// }
	fmt.Printf("hex.EncodeToString(t1): %v\n", string(b))
}


func TestSss(t *testing.T) {
	timestamp := time.Now().Unix()
	byteTimestamp := []byte{
		byte(timestamp & 0xFF),
		byte((timestamp >> 8) & 0xFF),
		byte((timestamp >> 16) & 0xFF),
		byte((timestamp >> 24) & 0xFF),
	}
	fmt.Println(hex.EncodeToString(byteTimestamp)) //时间戳
}
 
func TestScc(t *testing.T) {
	b, _ := hex.DecodeString("0102006c000a033432313938303141323655434330303030303030313432313938303141323655434330303030303030310850623e987b547500001105543332304d000000002102014286d3ffffffff01040000000686d301010083004ada2d6500818e9010000001001")
	fmt.Println(CRC(b))
}


// Direction Direction
type Direction uint8

// Direction List
const (
	DirectionClientServer Direction = 0x00
	DirectionServerClient Direction = 0x01
)

// Type Type
type Type uint8

//00表示command是通用的
//01表示command是clusterId特殊的
const (
	FrameTypeGlobal Type = 0x00
	FrameTypeLocal  Type = 0x01
)

// Control Control
type Control struct {
	FrameType              Type      `bits:"0b00000011" bitmask:"start"`
	ManufacturerSpecific   uint8     `bits:"0b00000100"`
	Direction              Direction `bits:"0b00001000"`
	DisableDefaultResponse uint8     `bits:"0b00010000"`
	Reserved               uint8     `bits:"0b11100000" bitmask:"end"`
}

// Frame Frame
type Frame struct {
	FrameControl              *Control
	ManufacturerCode          uint16 `cond:"uint:FrameControl.ManufacturerSpecific==1"`
	TransactionSequenceNumber uint8
	CommandIdentifier         uint8
	Payload                   []uint8
}

// Decode Decode
func Decode(buf []uint8) *Frame {
	frame := &Frame{}
	bin.Decode(buf, frame)
	return frame
}

// Encode Encode
func Encode(frame *Frame) []uint8 {
	return bin.Encode(frame)
}