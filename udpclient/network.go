package udpclient

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"github.com/coocood/freecache"
	"github.com/lmj/mqtt-clients-demo/common"
	"github.com/lmj/mqtt-clients-demo/config"
	"github.com/lmj/mqtt-clients-demo/logger"
)

// 保活信息记录
var MsgCheckTimeID = &sync.Map{}
var ClientForEveryMsg = &sync.Map{}
var MsgAllClientPayload = &sync.Map{} //记录一个客户端的保活消息载体

func (c *Client) setupConnection(address string) {
	addr, err := net.ResolveUDPAddr("udp4", address)
	errorCheck(err, "setupConnection", true)
	logger.Log.Warnln("> server address: ", addr.String(), " ... connecting ")
	conn, err := net.DialUDP("udp4", nil, addr)
	c.Connection = conn
	//also listen from requests from the server on a random port
	listeningAddress, err := net.ResolveUDPAddr("udp4", ":0")
	errorCheck(err, "setupConnection", true)
	logger.Log.Infoln("...CONNECTED! ")
	conn, err = net.ListenUDP("udp4", listeningAddress)
	errorCheck(err, "setupConnection", true)
	logger.Log.Infoln("listening on: local: ", conn.LocalAddr())
}

func (c *Client) readFromSocket(buffersize int) {
	for {
		var b = make([]byte, buffersize)
		n, addr, err := c.Connection.ReadFromUDP(b[:])
		if err != nil {
			return
		}
		if n > 0 {
			pack := packet{b[0:n], addr}
			select {
			case c.packets <- pack:
				continue
			case <-c.Kill:
				c.Connection.Close()
				return
			}
		}
		select {
		case <-c.Kill:
			c.Connection.Close()
			return
		default:
			continue
		}
	}
}

 

// 回了ack的话就更新缓存中的时间
// 修改一下其中的序列号即可重新发送
// 这里仅仅只将待传递的消息送出去 ---》 帧序列和待发送的类型
func (c *Client) processPackets() {
	for pack := range c.packets {
		logger.Log.Warnln("/udpclient/processPackets/ Receive from", pack.returnAddress.IP.String(), ":", pack.returnAddress.Port, "Starting proc msg:", hex.EncodeToString(pack.bytes))
		jsoninfo, message, cacheKey := ParseUDPMsg(pack.bytes), "", ""
		if jsoninfo.MessageHeader.MsgType == DataMsgType.GeneralAck { //下行报文ack，处理这些ack的思路就是在缓存中清理
			cacheKey = jsoninfo.TunnelHeader.FrameSN + jsoninfo.MessagePayload.Data
			fcache, _ := ClientForEveryMsg.Load(c.clientname)
			_, err1 := fcache.(*freecache.Cache).Get([]byte(cacheKey))
			mcache, _ := MsgAllClientPayload.Load(c.clientname)
			_, err2 := mcache.(*freecache.Cache).Get([]byte(cacheKey))
			if err1 != nil || err2 != nil {
				logger.Log.Errorln("DevEUI", cacheKey, "/udpclient/processPackets can't find Message, it has expired!")
			} else {
				if jsoninfo.MessagePayload.Data == DataMsgType.UpMsg.KeepAliveEvent {
					message = cacheKey
					c.msgType <- message
				} else { //目前来说，有且仅有这一种，除了保活消息以外其他消息从缓存里删除掉
					fcache.(*freecache.Cache).Del([]byte(cacheKey))
					mcache.(*freecache.Cache).Del([]byte(cacheKey))
				}
			}
		} else if jsoninfo.MessageHeader.MsgType == DataMsgType.DownMsg.TerminalGetPort { //终端端口获取
			message = jsoninfo.TunnelHeader.FrameSN + DataMsgType.UpMsg.TerminalReportPort + DataMsgType.DownMsg.TerminalGetPort
			c.msgType <- message
		} else {
			//todo 尚未知道
		}
	}
}

// 产生设备，这里会复用T320M的设备信息，
// 并且构建新的客户端列表
func GenMode(nums int) []TerminalInfo {
	defer func() {
		err := recover()
		if err != nil {
			logger.Log.Errorln("/udpclient/GenMode: An error occurred while generating the device: ", err)
		}
	}()
	logger.Log.Infoln("/udpclient/GenMode: Generate new terminal information based on T320 device information....")
	TnfGroup = make([]TerminalInfo, nums)
	for i := 0; i < nums; i++ {
		sufFix := "%0" + config.DEVICE_SN_LEFT_LEN + "d"
		devSN := fmt.Sprintf(config.DEVICE_SN_PRE+config.DEVICE_SN_MID+sufFix, config.DEVICE_SN_SUF_START_BY+i)
		TnfGroup[i] = TerminalInfo{}
		TnfGroup[i].FirstAddr = getRand(40, false) //模拟终端设备序列号
		Mac, ok := common.DevSNwithMac.Load(devSN + "M")
		if !ok {
			TnfGroup[i].ThirdAddr = getRand(12, false)
		}
		realMac := Mac.(string)
		TnfGroup[i].ThirdAddr = realMac
		TnfGroup[i].IotModule = getRand(2, true)
		TnfGroup[i].SecondAddr = TnfGroup[i].FirstAddr
		TnfGroup[i].key = TnfGroup[i].FirstAddr
		TnfGroup[i].Client = NewClient()
		TnfGroup[i].Client.setupConnection(config.UDP_SERVER_HOST + ":" + strconv.Itoa(config.UDP_SERVER_PORT))
		TnfGroup[i].msgType = make(chan string)
		TnfGroup[i].Client.clientname = "TerminalInfo" + strconv.Itoa(i)
		TnfGroup[i].devSN = devSN
		ClientForEveryMsg.Store(TnfGroup[i].Client.clientname, freecache.NewCache(5*1024*1024))
		MsgAllClientPayload.Store(TnfGroup[i].Client.clientname, freecache.NewCache(5*1024*1024))
	}
	return TnfGroup
}

/**
 * @Descrption: 发送消息，主动触发机制，通道触发
 * @param {TerminalInfo} Terminal
 * @param {*Client} c
 * @return {*}
 */
func sendMsg(Terminal TerminalInfo, c *Client) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Log.Errorln("/udpclient/sendMsg: send message procedure exist problem, errors' are", err)
		}
	}()
	for msg := range c.msgType { //主动发起保活、处理上行报文的时候，发送的消息都需要ack处理,)
		go func() {
			prepareSend, FrameSN, msgType, msgLoad := "", strings.Repeat(msg[0:4], 1), strings.Repeat(msg[4:8], 1), strings.Repeat(msg[8:], 1)
			if msgType == DataMsgType.UpMsg.KeepAliveEvent { //主动保活
				logger.Log.Infoln("DevEUI:", Terminal.key, "/udpclient/sendMsg: proc keepalive message")
				cacheKey := FrameSN + msgType
				msgCache, hasClient1 := MsgAllClientPayload.Load(c.clientname)
				_, err := msgCache.(*freecache.Cache).Get([]byte(cacheKey))
				fcache, hasClient2 := ClientForEveryMsg.Load(c.clientname)
				_, err1 := fcache.(*freecache.Cache).Get([]byte(cacheKey))
				if !hasClient1 || !hasClient2 {
					logger.Log.Errorln("DevEUI:", Terminal.key, "/udpclient/sendMsg: There is no corresponding client object in the cache, This Client is Closed!")
					c.Kill <- true
					return
				}
				if FrameSN == DataMsgType.UpMsg.StartEvent && err != nil && err1 != nil { //没缓存且首次发送
					prepareSend = encMsg(msgType, Terminal, FrameSN, "")
					logger.Log.Infoln("DevEUI:", Terminal.key, "/udpclient/sendMsg: first generate message!")
				} else { //非首次内容则需要缓存提取，仅保活消息
					logger.Log.Infoln("DevEUI:", Terminal.key, "/udpclient/sendMsg: again generate message!")
					msgLoad, err := msgCache.(*freecache.Cache).Get([]byte(cacheKey))
					if err != nil {
						logger.Log.Errorln("DevEUI:", Terminal.key, "/udpclient/sendMsg: The current message has expired and there are network fluctuations")
						return
					}
					FrameSN = makeHex(FrameSN, 4)
					prepareSend = CreateNewMsg(FrameSN, msgLoad)
					msgCache.(*freecache.Cache).Del([]byte(cacheKey)) //清理缓存
					fcache.(*freecache.Cache).Del([]byte(cacheKey))
					time.Sleep(time.Second * time.Duration(config.UDP_ALIVE_CHECK_TIME)) //重新生成消息后再清理缓存
				}
			} else if msgType == DataMsgType.UpMsg.TerminalJoinEvent { 
				logger.Log.Infoln("DevEUI:", Terminal.key, "/udpclient/sendMsg: proc TerminalJoin message")
				prepareSend = encMsg(msgType, Terminal, FrameSN, "")
			} else if msgType == DataMsgType.UpMsg.TerminalLeaveEvent { 
				logger.Log.Infoln("DevEUI:", Terminal.key, "/udpclient/sendMsg: proc TerminalLeave message")
				prepareSend = encMsg(msgType, Terminal, FrameSN, "")
			} else if msgType == DataMsgType.UpMsg.TerminalReportPort { //收到终端上报数据时, 需要主动回复ack
				logger.Log.Infoln("DevEUI:", Terminal.key, "/udpclient/sendMsg: send ack message about TerminalReportPort")
				go sendACK(DataMsgType.GeneralAck, DataMsgType.DownMsg.TerminalGetPort, FrameSN, c, Terminal)
				logger.Log.Infoln("DevEUI:", Terminal.key, "/udpclient/sendMsg: send TerminalReportPort message")
				prepareSend = encMsg(DataMsgType.UpMsg.TerminalReportPort, Terminal, FrameSN, msgLoad)
			}
			byteOfPrepareSend, _ := hex.DecodeString(prepareSend)
			_, err := c.Connection.Write(byteOfPrepareSend)
			logger.Log.Infoln("/udpclient/sendMsg: DevEUI:", Terminal.key, "/udpclient/sendMsg , the msg content is :", prepareSend)
			if err != nil {
				logger.Log.Errorln("/udpclient/sendMsg: The connection", c.clientname, "is disconnected")
				c.Kill <- true
				return
			}
			updateTime, keyOfSend := time.Now().UnixNano(), FrameSN+msgType
			fcache, _ := ClientForEveryMsg.Load(c.clientname)
			fcache.(*freecache.Cache).Set([]byte(keyOfSend), []byte(strconv.FormatInt(updateTime, 10)), config.UDP_ALIVE_CHECK_TIME*5)
			Msgfcache, _ := MsgAllClientPayload.Load(c.clientname)
			Msgfcache.(*freecache.Cache).Set([]byte(keyOfSend), byteOfPrepareSend, config.UDP_ALIVE_CHECK_TIME*5)
			go procKeepAliveMsgFreeCache(keyOfSend, config.UDP_ALIVE_CHECK_TIME, c.clientname)
			go reSendMsg(Terminal, c, 1, keyOfSend) //所有消息都触发重发机制，除了额外封装的ack信息
		}()
	}
}

// 保活消息重发
func reSendMsg(Terminal TerminalInfo, c *Client, resendTime int, key string) {
	if resendTime == 4 {
		return
	}
	timer := time.NewTimer((time.Second * time.Duration(config.UDP_ALIVE_CHECK_TIME))) //每15秒重发
	defer timer.Stop()
	<-timer.C
	msgCache, hasClient := MsgAllClientPayload.Load(c.clientname)
	if !hasClient {
		logger.Log.Warnln("/udpclient/reSendMsg: The connection", c.clientname, "is disconnected")
		return
	}
	data, err := msgCache.(*freecache.Cache).Get([]byte(key))
	if err == nil {
		preSend := hex.EncodeToString(data)
		logger.Log.Warnln("/udpclient/reSendMsg: Devui:", key, "Generate ReSend Message", resendTime, "Times")
		byteOfPresend, errIn := hex.DecodeString(preSend)
		if errIn != nil {
			logger.Log.Errorln("/udpclient/reSendMsg: ReDecode is Fail, Illegal Message!")
			return
		}
		_, errIn = c.Connection.Write(byteOfPresend)
		if errIn != nil {
			logger.Log.Errorln("/udpclient/reSendMsg: The connection", c.clientname, " is disconnected")
			c.Kill <- true
			return
		}
		reSendMsg(Terminal, c, resendTime+1, key)
	}
}

/*
仅用于设备侧的ack消息回复
*/
func sendACK(msgType, load, FrameSN string, c *Client, ter TerminalInfo) {
	prepareMsg := encMsg(msgType, ter, FrameSN, load)
	byteOfPrepareSend, _ := hex.DecodeString(prepareMsg)
	_, err := c.Connection.Write(byteOfPrepareSend)
	if err != nil {
		logger.Log.Errorln("/udpclient/sendACK: The connection", c.clientname, "is disconnected")
		c.Kill <- true
		return
	}
}
