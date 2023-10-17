package common

type Direction string
type Type string
type Status string
type AttrDataType string
type AttrValue string
type ClusterId string

// zcl小消息回复格式
type DataStructure struct {
	Msgstatus       Status
	MsgAttrDataType AttrDataType
	MsgValue        string //hex格式
}

// zcl下发配置命令消息
type ConfigMsg struct {
	Direction     string
	AttrIdent     string
	AttrDataType  string
	MinnInterval  string
	MaxnInterval  string
	ReportAble    string
	TimeOutPeriod string
}

// 消息头
type ZclHeader struct {
	FrameCtrl      string
	TransactionSec string
	CommandIdent   string
}

// /消息载荷
type ZCLPayLoad struct {
}

var ClusterToAttr map[string]map[string]DataStructure

func init() {
	ClusterToAttr = map[string]map[string]DataStructure{
		"0000": map[string]DataStructure{ //基础信息读取
			"0000": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgUint,
				MsgValue:        "02",
			},
			"0100": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgUint,
				MsgValue:        "11",
			},
			"0200": DataStructure{
				Msgstatus: ResError,
			},
			"0300": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgUint,
				MsgValue:        "10",
			},
			"0400": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgString,
				MsgValue:        "064865696d616e",
			},
			"0500": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgString,
				MsgValue:        "09536d617274506c7567",
			},
			"0600": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgString,
				MsgValue:        "09323031382e312e3131",
			},
			"0700": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgEnum,
				MsgValue:        "01",
			},
		},
		//开关量读取,由服务端下发了对应配置，对下发的配置进行定时上送
		"0006": map[string]DataStructure{
			"0000": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgIsOpen, //反馈开关量的情况
				MsgValue:        "01",
			},
		},
		"0702": map[string]DataStructure{
			"0000": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgUint48,
				MsgValue:        "000000000000",
			},
			"0002": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgMap,
				MsgValue:        "00",
			},
			"0003": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgEnum,
				MsgValue:        "00",
			},
			"0103": DataStructure{
				Msgstatus:       ResError,
				MsgAttrDataType: ResMsgUint24,
				MsgValue:        "010000",
			},
			"0203": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgUint24,
				MsgValue:        "102700",
			},
			"0303": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgMap,
				MsgValue:        "aa",
			},
			"0603": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgMap,
				MsgValue:        "00",
			},
			"0004": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgint24,
				MsgValue:        "2f0000",
			},
		},
		"0b04": map[string]DataStructure{
			"0000": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgMap32,
				MsgValue:        "000000000000",
			},
			"0505": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResmsgUint16,
				MsgValue:        "ba5b",
			},
			"0508": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResmsgUint16,
				MsgValue:        "0200",
			},
			"050b": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgUint24,
				MsgValue:        "2f00",
			},
			"0510": DataStructure{
				Msgstatus: ResError,
			},
			"0600": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResmsgUint16,
				MsgValue:        "0100",
			},
			"0601": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResmsgUint16,
				MsgValue:        "6400",
			},
			"0602": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResmsgUint16,
				MsgValue:        "0100",
			},
			"0603": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResmsgUint16,
				MsgValue:        "6400",
			},
			"0604": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResmsgUint16,
				MsgValue:        "0100",
			},
			"0605": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgint24,
				MsgValue:        "0a00",
			},
			"0800": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgMap,
				MsgValue:        "0000",
			},
			"0801": DataStructure{
				Msgstatus:       ResSuccess,
				MsgAttrDataType: ResMsgMap,
				MsgValue:        "ffff",
			},
		},
	}

}

const (
	//消息传递方向
	DirectionClientServer Direction = "00"
	DirectionServerClient Direction = "01"

	//clusterid的情况
	FrameTypeGlobal Type = "00"
	FrameTypeLocal  Type = "01"

	//回应消息的状态
	ResSuccess Status = "00"
	ResError   Status = "86" //此时后续就不会跟随任何消息

	//回应的消息类型
	ResMsgNoData AttrDataType = "00"
	ResMsgIsOpen AttrDataType = "10"
	ResMsgMap    AttrDataType = "18"
	ResMsgMap16  AttrDataType = "19"
	ResMsgMap32  AttrDataType = "1b"
	ResMsgUint   AttrDataType = "20" //后续接接一个无符号整型
	ResmsgUint16 AttrDataType = "21"
	ResMsgUint24 AttrDataType = "22"
	ResMsgUint48 AttrDataType = "25"
	ResMsgint24  AttrDataType = "2a"
	ResMsgint16  AttrDataType = "29"
	ResMsgEnum   AttrDataType = "30" //后接一个字节的枚举，默认写成字符串
	ResMsgString AttrDataType = "42" //后跟字符串

)
