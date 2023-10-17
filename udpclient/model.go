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
	DataMsgType.UpMsg.StartEvent = "0000"
	DataMsgType.UpMsg.TerminalReportPort = "2206"
	DataMsgType.UpMsg.TerminalInfoUp = "2102"
	DataMsgType.UpMsg.TerminalSvcDiscoverRsp = "2212"
	DataMsgType.UpMsg.TerminalPortBindRsp = "2208"
	DataMsgType.UpMsg.TerminalAccessRsp = "2210"

	DataMsgType.DownMsg.NeedAck = "1111"
	DataMsgType.DownMsg.TerminalGetPort = "2205"
	DataMsgType.DownMsg.TerminalInfoDown = "2101"
	DataMsgType.DownMsg.TerminalSvcDiscoverReq = "2211"
	DataMsgType.DownMsg.TerminalPortBindReq = "2207"
	DataMsgType.DownMsg.TerminalAccessReq = "220f"

	DataMsgType.GeneralAck = "2006"
}

type MessageType struct {
	UpMsg      UpMsg
	DownMsg    DownMsg
	GeneralAck string
}

type UpMsg struct {
	TerminalJoinEvent      string
	TerminalLeaveEvent     string
	DataUpEvent            string
	KeepAliveEvent         string
	StartEvent             string
	TerminalReportPort     string
	TerminalInfoUp         string
	TerminalSvcDiscoverRsp string //终端服务发现回复
	TerminalPortBindRsp    string //终端端口绑定
	TerminalAccessRsp      string //终端允许入网回复
}

type DownMsg struct {
	TerminalInfoDown       string
	TerminalGetPort        string
	NeedAck                string
	TerminalSvcDiscoverReq string //终端服务发现请求
	TerminalPortBindReq    string //终端端口绑定
	TerminalAccessReq      string //终端允许入网
}

// some bytes associated with an address
type packet struct {
	bytes         []byte
	returnAddress *net.UDPAddr
}

type Client struct {
	Connection *net.UDPConn
	messages   chan Message
	packets    chan packet
	Kill       chan bool
	msgType    chan string
	clientname string
}

type Message struct {
	Type    MessageType
	Message []byte
}

// create a new client.
func NewClient() *Client {
	return &Client{
		packets:  make(chan packet),
		messages: make(chan Message),
		Kill:     make(chan bool),
		msgType:  make(chan string), //通道传递
	}
}

// 模拟设备的链路信息
type TerminalInfo struct {
	FirstAddr  string //地址一级
	SecondAddr string
	ThirdAddr  string
	IotModule  string
	key        string
	Client     *Client     //设备对应连接客户端
	msgType    chan string //通道传递
	devSN      string
	DevEUI     string
}

// JSONInfo JSONInfo
type JSONInfo struct {
	TunnelHeader   TunnelHeader   `json:"tunnelHeader"`   //传输头
	MessageHeader  MessageHeader  `json:"messageHeader"`  //应用头
	MessagePayload MessagePayload `json:"messagePayload"` //应用数据
}

// TunnelHeader TunnelHeader
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

// LinkInfo LinkInfo
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

// AddrInfo AddrInfo
type AddrInfo struct {
	AddrInfo string
	MACAddr  string
	SN       string
}

// ExtendInfo ExtendInfo
type ExtendInfo struct {
	IsNeedUserName bool
	IsNeedDevType  bool
	IsNeedVender   bool
	IsNeedSec      bool
	AckOptionType  int
	ExtendData     string
}

// SecInfo SecInfo
type SecInfo struct {
	SecType    string
	SecDataLen string
	SecData    string
	SecID      string
}

// VenderInfo VenderInfo
type VenderInfo struct {
	VenderIDLen string
	VenderID    string
}

// DevTypeInfo DevTypeInfo
type DevTypeInfo struct {
	DevTypeLen string
	DevType    string
}

// UserNameInfo UserNameInfo
type UserNameInfo struct {
	UserNameLen string
	UserName    string
}

// MessageHeader MessageHeader
type MessageHeader struct {
	OptionType       string
	SN               string
	MsgType          string
	MessageHeaderLen int
}

// CtrlMsg CtrlMsg
type CtrlMsg struct {
	Addr    string
	SubAddr string
	PortID  string
	CtrlLen int
}

// MessagePayload MessagePayload
type MessagePayload struct {
	ModuleID string
	Ctrl     int
	Address  string
	SubAddr  string
	PortID   string
	Topic    string
	Data     string
}
