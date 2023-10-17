package udpclient

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/coocood/freecache"
	"github.com/lmj/mqtt-clients-demo/common"
	"github.com/lmj/mqtt-clients-demo/config"
	"github.com/lmj/mqtt-clients-demo/logger"
)

func errorCheck(err error, where string, kill bool) {
	if err != nil {
		if kill {
			logger.Log.WithError(err).Fatalln("Script Terminated", where)
		} else {
			logger.Log.WithError(err).Warnf("@ %s\n", where)
		}
	}
}

// 当前消息,key表示三级地址的串联，中间需要配置
// payload的使用，实在构造回复消息的时作为消息载体存在
// 终端网络地址也要设备进行标识
func encMsg(msgType string, dev TerminalInfo, FrameSN string, payLoad string) (string,error) {
	var (
		msg     strings.Builder
		randMsg string
		err error
	)
	msg.WriteString(config.UDP_VERSION_TYPE) //消息版本
	msg.WriteString(FrameSN)
	msg.WriteString("03") //地址级数
	msg.WriteString("34")
	msg.WriteString(dev.FirstAddr)
	msg.WriteString("34")
	msg.WriteString(dev.SecondAddr)
	randMsg = GetRand(4, false)                                                           //模拟消息载体信息
	if msgType == DataMsgType.UpMsg.KeepAliveEvent || msgType == DataMsgType.GeneralAck { //基本保活
		msg.WriteString("06")
		mMac, _ := common.DevSNwithMac.Load(dev.devSN + "M")
		msg.WriteString(mMac.(string))
	} else if msgType == DataMsgType.UpMsg.TerminalJoinEvent { //终端入网
		msg.WriteString("08")
		gmac, _ := common.DevSNwithMac.Load(dev.devSN + "G")
		msg.WriteString(gmac.(string))
		randMsg = "86d3" + dev.DevEUI + "8e" + "00" //终端地址，设备标识，设备功能，入网方式
	} else if msgType == DataMsgType.UpMsg.TerminalLeaveEvent { //终端离网
		msg.WriteString("08")
		gmac, _ := common.DevSNwithMac.Load(dev.devSN + "G")
		msg.WriteString(gmac.(string))
		randMsg = dev.DevEUI
	} else if msgType == DataMsgType.UpMsg.TerminalReportPort { //终端上报
		msg.WriteString("08")
		gmac, _ := common.DevSNwithMac.Load(dev.devSN + "G")
		msg.WriteString(gmac.(string))
		//profile编号，终端网络地址，状态成功，终端网络的地址，终端端口数量，终端端口列表
		randMsg = "0104" + "86d3" + "00" + "86d3" + "01" + "01"
	} else if msgType == DataMsgType.UpMsg.TerminalInfoUp { //回应下发命令
		msg.WriteString("08")
		gmac, _ := common.DevSNwithMac.Load(dev.devSN + "G")
		msg.WriteString(gmac.(string))
		profileId, clusterId, zclData := "", "", ""
		if payLoad[:7] == "REGULAR" {
			profileId, clusterId, zclData = "0104", "0006", payLoad[7:]
		} else {
			profileId, clusterId, zclData = payLoad[0:4], payLoad[4:8], payLoad[10:]
		}
		randMsg, err = CreatePreUpData(profileId, clusterId, "86d3", zclData, dev) //86d3早晚需要被替代掉
		if err != nil {
			return "", err
		}
	} else if msgType == DataMsgType.UpMsg.TerminalSvcDiscoverRsp {
		msg.WriteString("08")
		gmac, _ := common.DevSNwithMac.Load(dev.devSN + "G")
		msg.WriteString(gmac.(string))
		//Profile, 发送消息设备网络地址，状态成功，终端网络, 后续消息长度(不含自身，不含crc),终端端口，profile,终端类型id,终端版本，终端Incluster数量, Incluster列表，out数量，out列表
		randMsg = "0104" + "86d3" + "00" + "86d3" + "1a" + "01" + "0104" + "0051" + "00" + "07" + "0000000300040006000907020b040203001900"
	} else if msgType == DataMsgType.UpMsg.TerminalPortBindRsp { //终端绑定端口回复
		msg.WriteString("08")
		gmac, _ := common.DevSNwithMac.Load(dev.devSN + "G")
		msg.WriteString(gmac.(string))
		//Profile, 网络地址，状态
		randMsg = "0104" + "86d3" + "00"
	} else if msgType == DataMsgType.UpMsg.TerminalAccessRsp {
		msg.WriteString("08")
		gmac, _ := common.DevSNwithMac.Load(dev.devSN + "G")
		msg.WriteString(gmac.(string))
		//Profile, 网络地址，状态
		randMsg = "0104" + "86d3" + "00"
	} else if msgType == DataMsgType.ZigbeeGeneralFailed {
		msg.WriteString("08")
		gmac, _ := common.DevSNwithMac.Load(dev.devSN + "G")
		msg.WriteString(gmac.(string))
		randMsg = payLoad + "00000001" //错误码
	}
	msg.WriteString("0011")       //默认
	msg.WriteString("05")         //默认
	msg.WriteString("543332304d") //T320M
	//消息头信息
	msg.WriteString("0000")  //控制域
	msg.WriteString("0000")  //紧跟控制域的序列号暂时没什么用
	msg.WriteString(msgType) //消息类型
	//消息体信息
	msg.WriteString("01")       //物联网模组id  需要更改
	msg.WriteString("42")       //控制字H3C归一化报文
	msg.WriteString("86d3")     //地址信息
	msg.WriteString("")         //子地址
	msg.WriteString("")         //端口ID
	msg.WriteString("ffffffff") //厂商topic
	if len(payLoad) == 4 {      //ack应答专用
		msg.WriteString(payLoad)
	} else {
		msg.WriteString(randMsg)
	}
	message := generateLenOfMsg(msg)
	b, _ := hex.DecodeString(message)
	message += CRC(b)
	return message, nil
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
func GetRand(length int, isDigit bool) string {
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
		return 0, fmt.Errorf("Connect has Interrupted")
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
	return hex.EncodeToString(Msg)
}

// 回复PayLoad
func CreatePreUpData(p, c, addr, zclData string, Ter TerminalInfo) (string, error) {
	var (
		tempString   strings.Builder
		tempInString strings.Builder
	)
	tempString.WriteString(p)
	tempString.WriteString("0000") //gp
	tempString.WriteString(c)
	tempString.WriteString(addr)
	tempString.WriteString("01") //源终端端口
	tempString.WriteString("01") //目的终端端口
	tempString.WriteString("00") //wasBroadcast
	tempString.WriteString("83") //连接质量
	tempString.WriteString("00") //安全使用
	timestamp := time.Now().Unix()
	byteTimestamp := []byte{
		byte(timestamp & 0xFF),
		byte((timestamp >> 8) & 0xFF),
		byte((timestamp >> 16) & 0xFF),
		byte((timestamp >> 24) & 0xFF),
	}
	tempString.WriteString(hex.EncodeToString(byteTimestamp)) //时间戳
	tempString.WriteString("00")                              //传递序列符号
	zclHeader := common.ZclHeader{
		FrameCtrl:      zclData[0:2],
		TransactionSec: zclData[2:4],
		CommandIdent:   zclData[4:6],
	}
	zclRspHeader := common.ZclHeader{
		TransactionSec: zclHeader.TransactionSec,
	}
	postZclData := procReadBasic(zclData, c, 6)
	switch zclHeader.CommandIdent {
	case "00":
		zclRspHeader.FrameCtrl = "18"
		zclRspHeader.CommandIdent = "01"
	case "06":
		zclRspHeader.FrameCtrl = "08"
		zclRspHeader.CommandIdent = "07"
		go func() {
			interval, _ := strconv.ParseInt(zclData[16:18]+zclData[14:16], 16, 32)
			ticker := time.NewTicker(time.Duration(interval/10) * time.Second)
			ZclRegularReport++
			for {
				select {
				case <-ticker.C:
					//理论上应当采用保活帧序列号，但此处使用新起始号，zcl封装信息固定
					s := strconv.FormatInt(int64(ZclRegularReport), 16)
					logger.Log.Infoln("/udpclient/CreatePreUpData start regular send message!")
					if len(s) < 2 {
						s = "0" + s
					}
					message := "0000" + DataMsgType.UpMsg.TerminalInfoUp + "REGULAR" + "08" + s + "0a" + "00001001"
					Ter.Client.msgType <- message
				case <-RegularCh: //终端定时上报程序
					return
				default:
					continue
				}
			}
		}()
		postZclData = "00" //确认
	case "0a":
		zclRspHeader.FrameCtrl = zclHeader.FrameCtrl
		zclRspHeader.CommandIdent = zclHeader.CommandIdent
		postZclData = zclData[6:]
	default:
		// case "0b": 存在问题
	// 	zclRspHeader.FrameCtrl = "08"
	// 	zclRspHeader.CommandIdent = "07"
		return "", fmt.Errorf("/udpclient/encMsg can't recogenize the zcl message")
	}
	 
	tempInString.WriteString(zclRspHeader.FrameCtrl)
	tempInString.WriteString(zclRspHeader.TransactionSec)
	tempInString.WriteString(zclRspHeader.CommandIdent)
	tempInString.WriteString(postZclData)
	ZclCmd := strconv.FormatInt(int64(len(tempInString.String())/2), 16)
	if len(ZclCmd) < 2 {
		ZclCmd = "0" + ZclCmd
	}
	tempString.WriteString(ZclCmd)
	tempString.WriteString(tempInString.String())
	return tempString.String(), nil
}

// 处理普通读请求操作，没有就空
func procReadBasic(zclData, cluster string, preFix int) string {
	n := len(zclData)
	var result strings.Builder
	for i := preFix; i < n && i+4 <= n; i += 4 {
		curRead := zclData[i : i+4]
		result.WriteString(curRead + string(common.ClusterToAttr[cluster][curRead].Msgstatus) + string(common.ClusterToAttr[cluster][curRead].MsgAttrDataType) + string(common.ClusterToAttr[cluster][curRead].MsgValue))
	}
	return result.String()
}
