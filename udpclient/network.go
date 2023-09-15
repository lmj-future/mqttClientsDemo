
package udpclient

import (
	//"fmt"
	"encoding/hex"
	"net"
	"strconv"
	"sync"
	"time"

	"demo/config"

	log "github.com/Sirupsen/logrus"
	"github.com/coocood/freecache"
)

//保活信息记录
var MsgCheckTimeID = &sync.Map{}
var MsgForEvery = &sync.Map{}

func (c *Client) setupConnection(address string) {

	addr, err := net.ResolveUDPAddr("udp4", address)

	errorCheck(err, "setupConnection", true)
	log.Printf("> server address: %s ... connecting ", addr.String())

	conn, err := net.DialUDP("udp4", nil, addr)
	c.Connection = conn

	//also listen from requests from the server on a random port
	listeningAddress, err := net.ResolveUDPAddr("udp4", ":0")
	errorCheck(err, "setupConnection", true)
	log.Printf("...CONNECTED! ")

	conn, err = net.ListenUDP("udp4", listeningAddress)
	errorCheck(err, "setupConnection", true)

	log.Printf("listening on: local:%s\n", conn.LocalAddr())

}

func (c *Client) readFromSocket(buffersize int) {
	for {
		var b = make([]byte, buffersize)
		n, addr, err := c.Connection.ReadFromUDP(b[:])
		errorCheck(err, "readFromSocket", false)
		if n > 0 {
			pack := packet{b[0:n], addr}
			select {
			case c.packets <- pack:
				continue
			case <-c.kill:
				break
			}
		}
		select {
		case <-c.kill:
			break
		default:
			continue
		}
	}
}

//回了ack的话就更新缓存中的时间
//修改一下其中的序列号即可重新发送
func (c *Client) processPackets() {
	for pack := range c.packets {
		log.Warnln("Receive from ", pack.returnAddress.IP.String(), ":", pack.returnAddress.Port, " Starting proc msg Content is :", hex.EncodeToString(pack.bytes))
		jsoninfo := ParseUDPMsg(pack.bytes)
		key := jsoninfo.TunnelHeader.LinkInfo.FirstAddr + jsoninfo.TunnelHeader.FrameSN
		fcache, _ := MsgForEvery.Load(c.clientname)
		_, err := fcache.(*freecache.Cache).Get([]byte(key))
		if err != nil { // 没有找到键值，过期了
			log.Errorln("Message has Expired!")
		} else {
			fcache.(*freecache.Cache).Del([]byte(key))
			SN := makeHex(jsoninfo.TunnelHeader.FrameSN, 4)
			c.msgType <- SN //通知发送
		}
	}
}

//产生设备
func GenMode(nums int) []TerminalInfo {
	TnfGroup := make([]TerminalInfo, nums)
	for i := 0; i < nums; i++ {
		TnfGroup[i] = TerminalInfo{}
		TnfGroup[i].FirstAddr = getRand(40, false) //模拟获取设备
		time.Sleep(time.Microsecond * 50)
		TnfGroup[i].ThirdAddr = getRand(12, false)
		time.Sleep(time.Microsecond * 50)
		TnfGroup[i].IotModule = getRand(2, true)
		TnfGroup[i].SecondAddr = TnfGroup[i].FirstAddr
		TnfGroup[i].key = TnfGroup[i].FirstAddr
		TnfGroup[i].client = NewClient()
		TnfGroup[i].client.setupConnection(Opts.ServerAddress)
		TnfGroup[i].msgType = make(chan string)
		TnfGroup[i].client.clientname = TnfGroup[i].FirstAddr + strconv.Itoa(i)
		MsgForEvery.Store(TnfGroup[i].client.clientname, freecache.NewCache(10*1024*1024))
		ClientMapUDP.Set(strconv.Itoa(i), &TnfGroup[i].client)
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
	var (
		err           error
		byteOfPresend []byte
	)
	for SN := range c.msgType {
		if SN != "0000" { //非首轮消息停5s
			time.Sleep(5 * time.Second)
		}
		preSend := encKeepAliveMsg(DataMsgType.UpMsg.KeepAliveEvent, Terminal, SN)
		key := Terminal.FirstAddr + SN
		log.Warnln("Devui: ", key, "Generate Send Message...")
		byteOfPresend, err = hex.DecodeString(preSend)
		if err != nil {
			log.Errorln("Decode is Fail, Illegal Message!")
			continue //继续监听
		}
		_, err = c.Connection.Write(byteOfPresend)
		if err != nil {
			log.Errorln("The connection ", c.clientname, " is disconnected")
			close(Terminal.msgType)
			return
		}
		updateTime := time.Now().UnixNano() //设置消息过期时间
		fcache, _ := MsgForEvery.Load(c.clientname)
		fcache.(*freecache.Cache).Set([]byte(key), []byte(strconv.FormatInt(updateTime, 10)), config.UDP_ALIVE_CHECK_TIME*5)
		go procKeepAliveMsgFreeCache(key, config.UDP_ALIVE_CHECK_TIME, c.clientname) //这里不进行等待,转而去等待处理下一条消息
		go reSendMsg(Terminal, c, 1, SN, key)
	}
}

//保活消息重发,重发以后就走到keepalive环节了，如果也没有
func reSendMsg(Terminal TerminalInfo, c *Client, resendTime int, SN, key string) {
	if resendTime == 4 {
		return
	}
	timer := time.NewTimer((time.Second * time.Duration(config.UDP_ALIVE_CHECK_TIME))) //每五秒重发
	defer timer.Stop()
	<-timer.C
	fcache, _ := MsgForEvery.Load(c.clientname)
	_, err := fcache.(*freecache.Cache).Get([]byte(key))
	if err == nil { //重发
		preSend := encKeepAliveMsg(DataMsgType.UpMsg.KeepAliveEvent, Terminal, SN) //封装消息
		log.Warnln("Devui: ", key, "Generate ReSend Message", resendTime ,"Times")
		byteOfPresend, errIn := hex.DecodeString(preSend)
		if errIn != nil {
			log.Errorln("ReDecode is Fail, Illegal Message!")
			return
		}
		_, errIn = c.Connection.Write(byteOfPresend) //直接写，但没有任何反馈
		if errIn != nil {
			log.Errorln("The connection ", c.clientname, " is disconnected")
			close(Terminal.msgType)
			return
		}
		reSendMsg(Terminal, c, resendTime+1, SN, key)
	}
}
