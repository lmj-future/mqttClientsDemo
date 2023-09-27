package udpclient

import (
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/lmj/mqtt-clients-demo/logger"
)

//ParseUDPMsg ParseUDPMsg
func ParseUDPMsg(udpMsg []byte) JSONInfo {
	var tunnelHeader = parseTunnelHeader(udpMsg)
	var messageHeader = parseMessageHeader(udpMsg, tunnelHeader.TunnelHeaderLen)
	var messagePayload = parseMessagePayload(udpMsg, tunnelHeader.TunnelHeaderLen+messageHeader.MessageHeaderLen)
	return JSONInfo{
		TunnelHeader:   tunnelHeader,   //传输头
		MessageHeader:  messageHeader,  //应用头
		MessagePayload: messagePayload, //应用数据
	}
}

func parseTunnelHeader(udpMsg []byte) TunnelHeader {
	var tunnelHeader = TunnelHeader{
		Version:  hex.EncodeToString(append(udpMsg[:0:0], udpMsg[0:2]...)), //协议版本（2byte）
		FrameLen: hex.EncodeToString(append(udpMsg[:0:0], udpMsg[2:4]...)), //帧长度（2byte）
		FrameSN:  hex.EncodeToString(append(udpMsg[:0:0], udpMsg[4:6]...)), //帧序列号（2byte）
	}
	var linkInfo = LinkInfo{ //链路信息
		AddrNum: hex.EncodeToString(append(udpMsg[:0:0], udpMsg[6:7]...)), //地址个数
	}
	var linkInfoLen = 0
	addrNum, _ := strconv.ParseInt(linkInfo.AddrNum, 16, 0)   //字符转16进制数
	var address = make([]AddrInfo, int(addrNum)) //地址列表
	for index := 0; index < int(addrNum); index++ {
		var addrInfo = AddrInfo{} //n级地址
		addrInfo.AddrInfo = hex.EncodeToString(append(udpMsg[:0:0], udpMsg[7+linkInfoLen:8+linkInfoLen]...))
		addrLen, _ := strconv.ParseInt(addrInfo.AddrInfo, 16, 0) //n级地址长度
		var addrLength = int(addrLen) & 31                       //bit 4-0, length
		if (addrLen >> 5) == 0 {                                 //bit 7-5, 0b000,mac
			addrInfo.MACAddr = hex.EncodeToString(append(udpMsg[:0:0], udpMsg[8+linkInfoLen:8+linkInfoLen+addrLength]...))
		} else if (addrLen >> 5) == 1 { //bit 7-5, 0b001,sn
			addrInfo.SN = hex.EncodeToString(append(udpMsg[:0:0], udpMsg[8+linkInfoLen:8+linkInfoLen+addrLength]...))
		}
		address[index] = addrInfo
		linkInfoLen += (addrLength + 1)
	}
	if addrNum == 3 { //3级
		if address[0].SN != "" {
			linkInfo.ACMac = address[0].SN
		} else {
			linkInfo.ACMac = address[0].MACAddr
		}
		if address[1].SN != "" {
			linkInfo.APMac = address[1].SN
		} else {
			linkInfo.APMac = address[1].MACAddr
		}
		if address[2].SN != "" {
			linkInfo.T300ID = address[2].SN
		} else {
			linkInfo.T300ID = address[2].MACAddr
		}
		linkInfo.FirstAddr = linkInfo.ACMac
		linkInfo.SecondAddr = linkInfo.APMac
		linkInfo.ThirdAddr = linkInfo.T300ID
	} else if addrNum == 2 { //2级
		if address[0].SN != "" {
			linkInfo.ACMac = address[0].SN
		} else {
			linkInfo.ACMac = address[0].MACAddr
		}
		if address[1].SN != "" {
			linkInfo.APMac = address[1].SN
		} else {
			linkInfo.APMac = address[1].MACAddr
		}
		linkInfo.FirstAddr = linkInfo.ACMac
		linkInfo.SecondAddr = linkInfo.APMac
	} else if addrNum == 1 { //1级
		if address[0].SN != "" {
			linkInfo.APMac = address[0].SN
		} else {
			linkInfo.APMac = address[0].MACAddr
		}
		linkInfo.FirstAddr = linkInfo.APMac
	}
	linkInfo.Address = address
	tunnelHeader.LinkInfo = linkInfo
	var extendData = hex.EncodeToString(append(udpMsg[:0:0], udpMsg[7+linkInfoLen:9+linkInfoLen]...))
	isNeedUserName, _ := strconv.ParseInt(extendData, 16, 0)
	isNeedDevType, _ := strconv.ParseInt(extendData, 16, 0)
	isNeedVender, _ := strconv.ParseInt(extendData, 16, 0)
	isNeedSec, _ := strconv.ParseInt(extendData, 16, 0)
	optionType, _ := strconv.ParseInt(extendData, 16, 0)
	var extendInfo = ExtendInfo{
		IsNeedUserName: (isNeedUserName & 32) == 32,
		IsNeedDevType:  (isNeedDevType & 16) == 16,
		IsNeedVender:   (isNeedVender & 8) == 8,
		IsNeedSec:      (isNeedSec & 4) == 4,
		AckOptionType:  int(optionType & 3), //0,1,2
		ExtendData:     extendData,
	}
	tunnelHeader.ExtendInfo = extendInfo //扩展字段
	var secInfo = SecInfo{}
	var secInfoLen = 0
	if extendInfo.IsNeedSec {
		secInfo.SecType = hex.EncodeToString(append(udpMsg[:0:0], udpMsg[9+linkInfoLen:10+linkInfoLen]...))
		secInfo.SecDataLen = hex.EncodeToString(append(udpMsg[:0:0], udpMsg[10+linkInfoLen:11+linkInfoLen]...))
		secDataLen, _ := strconv.ParseInt(secInfo.SecDataLen, 16, 0)
		secInfo.SecData = hex.EncodeToString(append(udpMsg[:0:0], udpMsg[11+linkInfoLen:11+linkInfoLen+int(secDataLen)]...))
		secInfo.SecID = hex.EncodeToString(append(udpMsg[:0:0], udpMsg[11+linkInfoLen+int(secDataLen):12+linkInfoLen+int(secDataLen)]...))
		secInfoLen = 1 + 1 + int(secDataLen) + 1
	}
	var venderInfo = struct {
		VenderIDLen string
		VenderID    string
	}{}
	var venderInfoLen = 0
	if extendInfo.IsNeedVender {
		venderInfo.VenderIDLen = hex.EncodeToString(append(udpMsg[:0:0], udpMsg[9+linkInfoLen+secInfoLen:10+linkInfoLen+secInfoLen]...))
		venderIDLen, _ := strconv.ParseInt(venderInfo.VenderIDLen, 16, 0)
		venderInfo.VenderID = string(append(udpMsg[:0:0], udpMsg[10+linkInfoLen+secInfoLen:10+linkInfoLen+secInfoLen+int(venderIDLen)]...))
		venderInfoLen = 1 + int(venderIDLen)
	}
	var devTypeInfo = struct {
		DevTypeLen string
		DevType    string
	}{}
	var devTypeInfoLen = 0
	if extendInfo.IsNeedDevType {
		devTypeInfo.DevTypeLen = hex.EncodeToString(append(udpMsg[:0:0], udpMsg[9+linkInfoLen+secInfoLen+venderInfoLen:10+linkInfoLen+secInfoLen+venderInfoLen]...))
		devTypeLen, _ := strconv.ParseInt(devTypeInfo.DevTypeLen, 16, 0)
		devTypeInfo.DevType = string(append(udpMsg[:0:0], udpMsg[10+linkInfoLen+secInfoLen+venderInfoLen:10+linkInfoLen+secInfoLen+venderInfoLen+int(devTypeLen)]...))
		devTypeInfoLen = 1 + int(devTypeLen)
	}
	var userNameInfo = struct {
		UserNameLen string
		UserName    string
	}{}
	var userNameInfoLen = 0
	if extendInfo.IsNeedUserName {
		userNameInfo.UserNameLen = hex.EncodeToString(append(udpMsg[:0:0], udpMsg[9+linkInfoLen+secInfoLen+venderInfoLen+devTypeInfoLen:10+linkInfoLen+secInfoLen+venderInfoLen+devTypeInfoLen]...))
		userNameLen, _ := strconv.ParseInt(userNameInfo.UserNameLen, 16, 0)
		userNameInfo.UserName = string(append(udpMsg[:0:0], udpMsg[10+linkInfoLen+secInfoLen+venderInfoLen+devTypeInfoLen:10+linkInfoLen+secInfoLen+venderInfoLen+devTypeInfoLen+int(userNameLen)]...))
		userNameInfoLen = 1 + int(userNameLen)
	}
	tunnelHeader.SecInfo = secInfo       //加密字段
	tunnelHeader.VenderInfo = venderInfo //场景字段
	tunnelHeader.DevTypeInfo = devTypeInfo
	tunnelHeader.UserNameInfo = userNameInfo
	tunnelHeader.TunnelHeaderLen = 9 + linkInfoLen + secInfoLen + venderInfoLen + devTypeInfoLen + userNameInfoLen
	return tunnelHeader
}


func parseMessageHeader(udpMsg []byte, len int)MessageHeader {
	return MessageHeader{
		OptionType:       hex.EncodeToString(append(udpMsg[:0:0], udpMsg[len:len+2]...)),
		SN:               hex.EncodeToString(append(udpMsg[:0:0], udpMsg[len+2:len+4]...)),
		MsgType:          hex.EncodeToString(append(udpMsg[:0:0], udpMsg[len+4:len+6]...)),
		MessageHeaderLen: 6,
	}
}

func parseMessagePayload(udpMsg []byte, length int) MessagePayload {
	var messagePayload = MessagePayload{
		ModuleID: hex.EncodeToString(append(udpMsg[:0:0], udpMsg[length:length+1]...)),
	}
	ctrl, _ := strconv.ParseInt(hex.EncodeToString(append(udpMsg[:0:0], udpMsg[length+1:length+2]...)), 16, 0)
	messagePayload.Ctrl = int(ctrl)
	var ctrlMsg = parseCtrl(append(udpMsg[:0:0], udpMsg[length+2:]...), messagePayload.Ctrl)
	if len(ctrlMsg.Addr) == 4 {
		messagePayload.Address = ctrlMsg.Addr
	} else {
		messagePayload.Address = strings.ToUpper(ctrlMsg.Addr)
	}
	messagePayload.SubAddr = ctrlMsg.SubAddr
	messagePayload.PortID = ctrlMsg.PortID
	messagePayload.Topic = hex.EncodeToString(append(udpMsg[:0:0], udpMsg[length+2+ctrlMsg.CtrlLen:length+2+ctrlMsg.CtrlLen+4]...))
	messagePayload.Data = hex.EncodeToString(append(udpMsg[:0:0], udpMsg[length+2+ctrlMsg.CtrlLen+4:len(udpMsg)-2]...))

	return messagePayload
}

func parseCtrl(udpMsg []byte, ctrl int) CtrlMsg {
	var reserve = (ctrl >> 6) & 3
	var portIDLen = (ctrl >> 4) & 3
	var subAddrLen = (ctrl >> 2) & 3
	var addrLen = ctrl & 3
	switch portIDLen {
	case 0:
		portIDLen = 0
	case 1:
		portIDLen = 2
	case 2:
		portIDLen = 4
	case 3:
		portIDLen = 8
	default:
		portIDLen = 0
		logger.Log.Warnln("invalid portIDLen:", portIDLen)
	}
	switch subAddrLen {
	case 0:
		subAddrLen = 0
	case 1:
		subAddrLen = 2
	case 2:
		subAddrLen = 4
	case 3:
		subAddrLen = 8
	default:
		subAddrLen = 0
		logger.Log.Warnln("invalid subAddrLen:", subAddrLen)
	}
	switch reserve {
	case 0: //����1
		switch addrLen {
		case 0:
			addrLen = 6
		case 1:
			addrLen = 8
		case 2:
			addrLen = 16
		case 3:
			addrLen = 20
		default:
			addrLen = 0
			logger.Log.Warnln("invalid addrLen:", addrLen)
		}
	case 1: //����2
		switch addrLen {
		case 1:
			addrLen = 4
		case 2:
			addrLen = 2
		default:
			addrLen = 0
			logger.Log.Warnln("invalid addrLen:", addrLen)
		}
	}

	return CtrlMsg{
		Addr:    hex.EncodeToString(append(udpMsg[:0:0], udpMsg[0:addrLen]...)),
		SubAddr: hex.EncodeToString(append(udpMsg[:0:0], udpMsg[addrLen:addrLen+subAddrLen]...)),
		PortID:  hex.EncodeToString(append(udpMsg[:0:0], udpMsg[addrLen+subAddrLen:addrLen+subAddrLen+portIDLen]...)),
		CtrlLen: addrLen + subAddrLen + portIDLen,
	}
}