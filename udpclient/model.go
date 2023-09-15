
package udpclient

import (
	"net"
)

var DataMsgType MessageType

func init() {
	DataMsgType.UpMsg.DataUpEvent = "2102"
	DataMsgType.UpMsg.KeepAliveEvent = "2007"
	DataMsgType.UpMsg.TerminalJoinEvent = "2001"
	DataMsgType.UpMsg.TerminalLeaveEvent = "2002"
	DataMsgType.DownMsg.GeneralAck = "2006"
}

type MessageType struct {
	UpMsg   UpMsg
	DownMsg DownMsg
}

type UpMsg struct {
	TerminalJoinEvent  string
	TerminalLeaveEvent string
	DataUpEvent        string
	KeepAliveEvent     string
}

type DownMsg struct {
	GeneralAck string
}

//some bytes associated with an address
type packet struct {
	bytes         []byte
	returnAddress *net.UDPAddr
}

type Client struct {
	Connection *net.UDPConn
	port       int
	messages   chan Message
	packets    chan packet
	kill       chan bool
	msgType    chan string 
	clientname string
}

type Message struct {
	Type    MessageType
	Message []byte
}

//create a new client.
func NewClient() *Client {
	return &Client{
		packets:  make(chan packet),
		messages: make(chan Message),
		kill:     make(chan bool),
		msgType: make(chan string), //通道传递
	}
}

//模拟设备的链路信息
type TerminalInfo struct {
	FirstAddr  string //地址一级
	SecondAddr string
	ThirdAddr  string
	IotModule  string
	key        string
	client     *Client     //设备对应连接客户端
	msgType    chan string //通道传递

}

//JSONInfo JSONInfo
type JSONInfo struct {
	TunnelHeader   TunnelHeader   `json:"tunnelHeader"`   //传输头
	MessageHeader  MessageHeader  `json:"messageHeader"`  //应用头
	MessagePayload MessagePayload `json:"messagePayload"` //应用数据
}

//TunnelHeader TunnelHeader
type TunnelHeader struct {
	Version         string
	FrameLen        string
	FrameSN         string
	LinkInfo        LinkInfo
	ExtendInfo      ExtendInfo
	SecInfo         SecInfo
	VenderInfo      VenderInfo
	DevTypeInfo     DevTypeInfo
	UserNameInfo    UserNameInfo
	TunnelHeaderLen int
}

//LinkInfo LinkInfo
type LinkInfo struct {
	AddrNum    string
	ACMac      string
	APMac      string
	T300ID     string
	FirstAddr  string
	SecondAddr string
	ThirdAddr  string
	Address    []AddrInfo
}

//AddrInfo AddrInfo
type AddrInfo struct {
	AddrInfo string
	MACAddr  string
	SN       string
}

//ExtendInfo ExtendInfo
type ExtendInfo struct {
	IsNeedUserName bool
	IsNeedDevType  bool
	IsNeedVender   bool
	IsNeedSec      bool
	AckOptionType  int
	ExtendData     string
}

//SecInfo SecInfo
type SecInfo struct {
	SecType    string
	SecDataLen string
	SecData    string
	SecID      string
}

//VenderInfo VenderInfo
type VenderInfo struct {
	VenderIDLen string
	VenderID    string
}

//DevTypeInfo DevTypeInfo
type DevTypeInfo struct {
	DevTypeLen string
	DevType    string
}

//UserNameInfo UserNameInfo
type UserNameInfo struct {
	UserNameLen string
	UserName    string
}

//MessageHeader MessageHeader
type MessageHeader struct {
	OptionType       string
	SN               string
	MsgType          string
	MessageHeaderLen int
}

//CtrlMsg CtrlMsg
type CtrlMsg struct {
	Addr    string
	SubAddr string
	PortID  string
	CtrlLen int
}

//MessagePayload MessagePayload
type MessagePayload struct {
	ModuleID string
	Ctrl     int
	Address  string
	SubAddr  string
	PortID   string
	Topic    string
	Data     string
}
