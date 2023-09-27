package udpclient

import (
	"encoding/hex"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/coocood/freecache"
	"github.com/lmj/mqtt-clients-demo/common"
	"github.com/lmj/mqtt-clients-demo/config"
	"github.com/lmj/mqtt-clients-demo/logger"
)

// 0102005219bc03343231393830314132365538323041453030305735343231393830314132365538323041453030305735063cd2e5bb07dc001105543332304d00000000200700420000ffffffff000f347a
func errorCheck(err error, where string, kill bool) {
	if err != nil {
		if kill {
			logger.Log.WithError(err).Fatalln("Script Terminated", where)
		} else {
			logger.Log.WithError(err).Warnf("@ %s\n", where)
		}
	}
}


//当前消息,key表示三级地址的串联，中间需要配置
//只需要管模拟和记录即可
//该函数需要适配不同的消息类型1，2，3，4分别对应相关的需求命令
//payload的使用，实在构造回复消息的时作为消息载体存在
func encMsg(msgType string, dev TerminalInfo, FrameSN string, payLoad string) string {
	var (
		msg     strings.Builder
		randMsg string
	)
	msg.WriteString(config.UDP_VERSION_TYPE) //消息版本
	msg.WriteString(FrameSN)
	msg.WriteString("03") //地址级数
	msg.WriteString("34")
	msg.WriteString(dev.FirstAddr)
	msg.WriteString("34")
	msg.WriteString(dev.SecondAddr)
	randMsg = getRand(4, false) //模拟消息载体信息
	if msgType == DataMsgType.UpMsg.KeepAliveEvent || msgType == DataMsgType.GeneralAck {  //基本保活
		msg.WriteString("06")
		mMac, _ := common.DevSNwithMac.Load(dev.devSN + "M")
		msg.WriteString(mMac.(string))

	} else if msgType == DataMsgType.UpMsg.TerminalJoinEvent { //终端入网
		msg.WriteString("08")
		gmac, _ := common.DevSNwithMac.Load(dev.devSN + "G")
		msg.WriteString(gmac.(string))
		randMsg = "86d3" + gmac.(string) + "8e" + "00" //终端地址，设备标识，设备功能，入网方式
	} else if msgType == DataMsgType.UpMsg.TerminalLeaveEvent { //终端离网
		msg.WriteString("08")
		gmac, _ := common.DevSNwithMac.Load(dev.devSN + "G")
		msg.WriteString(gmac.(string))
		randMsg = "86d3" + gmac.(string) + "8e" + "00" //终端地址，设备标识，设备功能，入网方式
	} else if msgType == "2206" { //终端端口汇报
		msg.WriteString("08")
		gmac, _ := common.DevSNwithMac.Load(dev.devSN + "G")
		msg.WriteString(gmac.(string))
		//profile编号，终端网络地址，状态成功，终端网络的地址，终端端口数量，终端端口列表
		randMsg = "0104" + "86d3" + "00" + "86d3" + "01" + "01"
	}
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
	msg.WriteString("86d3")        //地址信息
	msg.WriteString("")            //子地址
	msg.WriteString("")            //端口ID
	msg.WriteString("ffffffff")    //厂商topic
	if payLoad == "" {
		msg.WriteString(randMsg)       //消息信息
	} else {
		msg.WriteString(payLoad)
	}
	message := generateLenOfMsg(msg)
	b, _ := hex.DecodeString(message)
	message += CRC(b)
	return message
}

// 生成完整消息，载入消息长度
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

// 验证回复消息
func Check(msg []byte) bool {
	return CRC(append(msg[:0:0], msg[:len(msg)-2]...)) ==
		hex.EncodeToString(append(msg[:0:0], msg[len(msg)-2:]...))
}

// 随机生成序列
func getRand(length int, isDigit bool) string {
	time.Sleep(time.Nanosecond)
	if length < 1 {
		logger.Log.Errorln("范围有误!")
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

// 检测保活信息的处理
func procKeepAliveMsgFreeCache(key string, t int, clientName string) {
	logger.Log.Infoln("Start detecting message survival!")
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
					logger.Log.Warnln("device :", key, "has not revc keep alive msg for",
						(time.Now().UnixNano()-updateTime)/int64(time.Second), "seconds. This Client exit!")
					fcache, _ := ClientForEveryMsg.Load(clientName)
					fcache.(*freecache.Cache).Del([]byte(key))
					fcache, _ = MsgAllClientPayload.Load(clientName)
					fcache.(*freecache.Cache).Del([]byte(key))
				}
			} else if err.Error() != "Entry not found" {
				logger.Log.Warnln("/udpclient/procKeepAliveMsgFreeCache", err)
			}
		}
	}
}

// 获取对应客户端中的key消息缓存的截至过期时间
func KeepAliveTimerFreeCacheGet(key, client string) (int64, error) {
	var updateTime int64
	var err error
	fcache, ok := ClientForEveryMsg.Load(client)
	if !ok {
		return 0, errors.New("Connect has Interrupted!")
	}
	value, err := fcache.(*freecache.Cache).Get([]byte(key))
	if err != nil {
		return 0, err
	} 
	updateTime, err = strconv.ParseInt(string(value), 10, 64)
	return updateTime, err
}


// real自增后转化成sz长度的字符串
func makeHex(real string, sz int) string {
	SN, _ := strconv.ParseInt(real, 16, 8)
	SN++
	frameSN := strconv.FormatInt(int64(SN), 16)
	for len(frameSN) < sz {
		frameSN = "0" + frameSN
	}
	return frameSN
}

/*对缓存中取到的消息进行确认，并生成需要重发的消息*/
func CreateNewMsg(FrameSN string, Msg []byte) string {
	szOfMsg := len(Msg)
	byteOfSN, _ := hex.DecodeString(FrameSN)
	Msg[4], Msg[5] = byteOfSN[0], byteOfSN[1]
	byteofCRC, _ := hex.DecodeString(CRC(Msg[0 : szOfMsg-2]))
	Msg[szOfMsg-2], Msg[szOfMsg-1] = byteofCRC[0], byteofCRC[1]
	return  hex.EncodeToString(Msg)
}

