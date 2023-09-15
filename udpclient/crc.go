
package udpclient

import (
	 
	"strconv"
	"strings"

	 
)

//CRC CRC
func CRC(msg []byte) string {
	var temp = 0
	var crc = 0xffff
	for i := 0; i < len(msg); i++ {
		crc ^= int(msg[i])
		for j := 0; j < 8; j++ {
			temp = 1 & crc
			crc >>= 1
			if temp == 1 {
				crc ^= 0xa001
			}
		}
	}
	crc ^= 0xffff
	var builder strings.Builder
	builder.WriteString("0000")
	builder.WriteString(strconv.FormatInt(int64(crc), 16))

	return strings.Repeat(builder.String()[builder.Len()-4:], 1)
}