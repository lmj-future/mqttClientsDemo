
package udpclient

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coocood/freecache"
)

// 0102005219bc03343231393830314132365538323041453030305735343231393830314132365538323041453030305735063cd2e5bb07dc001105543332304d00000000200700420000ffffffff000f347a
func errorCheck(err error, where string, kill bool) {
	if err != nil {
		if kill {

			log.WithError(err).Fatalln("Script Terminated", where)
		} else {
			log.WithError(err).Warnf("@ %s\n", where)
		}
	}
}

func sendJoin(msgType string) {

}

func sendLeave(msgType string) {

}

//当前消息,key表示三级地址的串联，中间需要配置
//只需要管模拟和记录即可
func encKeepAliveMsg(msgType string, dev TerminalInfo, sn string) string {
	var msg strings.Builder
	msg.WriteString(Opts.UdpVer) //消息版本
	msg.WriteString(sn)
	msg.WriteString("03") //地址级数默认走3级
	//默认SN SN MAC
	msg.WriteString("34")
	msg.WriteString(dev.FirstAddr)
	msg.WriteString("34")
	msg.WriteString(dev.SecondAddr)
	msg.WriteString("06")
	msg.WriteString(dev.ThirdAddr)
	msg.WriteString("0011")       //默认
	msg.WriteString("05")         //默认
	msg.WriteString("543332304d") //T320M标识
	//消息头信息
	msg.WriteString("0000")  //控制域
	msg.WriteString("0000")  //紧跟控制域的序列号暂时没什么用
	msg.WriteString(msgType) //消息类型
	//消息体信息
	msg.WriteString(dev.IotModule) //物联网模组id
	msg.WriteString("42")          //控制字H3C归一化报文
	msg.WriteString("0000")        //地址信息
	msg.WriteString("")            //子地址
	msg.WriteString("")            //端口ID
	msg.WriteString("ffffffff")    //厂商topic
	randMsg := getRand(4, false)
	msg.WriteString(randMsg) //消息信息
	message := generateLenOfMsg(msg)
	b, _ := hex.DecodeString(message)
	message += CRC(b) //生成家校验码校验
	//message = "010200525ee603343231393830314132365538323041453030305735343231393830314132365538323041453030305735063cd2e5bb07dc001105543332304d00000000200700420000ffffffff000f6022"

	return message
}



//生成完整消息，载入消息长度
func generateLenOfMsg(msg strings.Builder) string {
	sz := len(msg.String()) + 8 //包含crc
	szOfHex := strconv.FormatInt(int64(sz)/2, 16)
	for len(szOfHex) < 4 {
		szOfHex = "0" + szOfHex
	}
	res := msg.String()
	res = res[0:4] + szOfHex[:] + res[4:]
	return res
}

//验证回复消息
func Check(msg []byte) bool {
	return CRC(append(msg[:0:0], msg[:len(msg)-2]...)) ==
		hex.EncodeToString(append(msg[:0:0], msg[len(msg)-2:]...))
}

//随机生成序列
func getRand(length int, isDigit bool) string {
	rand.Seed(time.Now().UnixNano())
	if length < 1 {
		log.Errorln("范围有误!")
		return ""
	}
	var char string
	if isDigit {
		char = "0123456789"
	} else {
		char = "abcdef0123456789"
	}
	charArr := strings.Split(char, "")
	charlen := len(charArr)
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))
	var rchar string = ""
	for i := 1; i <= length; i++ {
		rchar = rchar + charArr[ran.Intn(charlen)]
	}
	return rchar
}

//检测保活信息的处理
func procKeepAliveMsgFreeCache(key string, t int, clientName string) {
	var timerID = time.NewTimer(time.Duration(t*3+3) * time.Second)
	MsgCheckTimeID.Store(key, timerID)
	<-timerID.C
	timerID.Stop()
	if value, ok := MsgCheckTimeID.Load(key); ok {
		if timerID == value.(*time.Timer) {
			MsgCheckTimeID.Delete(key)
			updateTime, err := KeepAliveTimerFreeCacheGet(key, clientName)
			if err == nil {
				if time.Now().UnixNano()-updateTime > int64(time.Duration(3*t)*time.Second) {
					log.Warnln("device :", key, "has not revc keep alive msg for",
						(time.Now().UnixNano()-updateTime)/int64(time.Second), "seconds. This Client exit!")
					fcache, _ := MsgForEvery.Load(clientName)
					fcache.(*freecache.Cache).Del([]byte(key))
				}
			}
		}
	}
}

//获取对应客户端中的key消息缓存情况
func KeepAliveTimerFreeCacheGet(key, client string) (int64, error) {
	var updateTime int64
	var err error
	fcache, _ := MsgForEvery.Load(client)
	value, err := fcache.(*freecache.Cache).Get([]byte(key))
	if err == nil {
		updateTime, err = strconv.ParseInt(string(value), 10, 64)
	}
	return updateTime, err
}

func (T *TerminalInfo) toString() {
	fmt.Println("firstAddr is : ", T.FirstAddr, "secondAddr is : ", T.SecondAddr, "thirdAddr is :", T.ThirdAddr)
}

//real自增后转化成sz长度的字符串
func makeHex(real string, sz int) string {
	SN, _ := strconv.ParseInt(real, 16, 8)
	SN++
	frameSN := strconv.FormatInt(int64(SN), 16)
	for len(frameSN) < sz {
		frameSN = "0" + frameSN
	}
	return frameSN
}
