package common

import mqtt "github.com/eclipse/paho.mqtt.golang"

type MqttClientInfo struct {
	DevSN       string
	Client      mqtt.Client
	Connectting bool
}

type ModTopoInfo struct {
	PortID    int    `json:"portID"`
	ModStatus int    `json:"modStatus"`
	Type      int    `json:"type"`
	ModModel  string `json:"modModel"`
}

type NodeTopoInfo struct {
	CanID           int           `json:"canID"`
	NodeSN          string        `json:"nodeSN"`
	NodeModel       string        `json:"nodeModel"`
	ModTopoInfoList []ModTopoInfo `json:"modTopoInfoList"`
}

type ModStatus struct {
	PortID     int    `json:"portID"`
	ModStatus  int    `json:"modStatus"`
	Type       int    `json:"type"`
	ModModel   string `json:"modModel"`
	ModPullCfg int    `json:"modPullCfg"`
}

type NodeStatus struct {
	CanID         int         `json:"canID"`
	NodeSN        string      `json:"nodeSN"`
	NodeModel     string      `json:"nodeModel"`
	NodeStatus    int         `json:"nodeStatus"`
	NodeAging     int         `json:"nodeAging"`
	ModStatusList []ModStatus `json:"modStatusList"`
}

type ModBasicInfo struct {
	Type           int    `json:"type"`
	PortID         int    `json:"portID"`
	ModModel       string `json:"modModel"`
	ModHwVersion   string `json:"modHwVersion"`
	ModSwVersion   string `json:"modSwVersion"`
	ModSwInVersion string `json:"modSwInVersion"`
	ModMAC         string `json:"modMAC"`
	ModSN          string `json:"modSN"`
}

type NodeBasicInfo struct {
	NodeMAC              string         `json:"nodeMAC"`
	NodeIP               string         `json:"nodeIP"`
	Gmac                 string         `json:"gmac"`
	CanID                int            `json:"canID"`
	NodeSN               string         `json:"nodeSN"`
	NodeModel            string         `json:"nodeModel"`
	NodeHwVersion        string         `json:"nodeHwVersion"`
	NodeSwVersion        string         `json:"nodeSwVersion"`
	NodeLastRebootReason int            `json:"nodeLastRebootReason"`
	CanAddr              int            `json:"canAddr"`
	ModBasicInfoList     []ModBasicInfo `json:"modBasicInfoList"`
}

type UpMsg struct {
	DevSN             string          `json:"devSN,omitempty"`
	Method            string          `json:"method,omitempty"`
	DevOption         string          `json:"devOption,omitempty"`
	DevMsgID          *int            `json:"devMsgID"`
	ConnectStatus     *int            `json:"connectStatus"`
	CompatibilityVer  *int            `json:"compatibilityVer"`
	Version           string          `json:"version,omitempty"`
	DevModelList      []string        `json:"devModelList,omitempty"`
	NodeBasicInfoList []NodeBasicInfo `json:"nodeBasicInfoList,omitempty"`
	NodeStatusList    []NodeStatus    `json:"nodeStatusList,omitempty"`
	NodeTopoInfoList  []NodeTopoInfo  `json:"nodeTopoInfoList,omitempty"`

	DevModel string `json:"devModel,omitempty"`
	Step     *int   `json:"step"`

	Config       interface{} `json:"config,omitempty"`
	LogConfig    interface{} `json:"logConfig,omitempty"`
	WifiConfig   interface{} `json:"wifiConfig,omitempty"`
	GroupConfig  interface{} `json:"groupConfig,omitempty"`
	NodeConfig   interface{} `json:"nodeConfig,omitempty"`
	ModuleConfig interface{} `json:"moduleConfig,omitempty"`
}

type Param struct {
	Step string `json:"step,omitempty"`
	Desc string `json:"desc,omitempty"`
}

type ProgressMsg struct {
	Id     int   `json:"id"`
	Params Param `json:"params,omitempty"`
}

type DownMsg struct {
	DevSN        string       `json:"devSN,omitempty"`
	Method       string       `json:"method,omitempty"`
	DevOption    string       `json:"devOption,omitempty"`
	Maintenance  string       `json:"maintenance,omitempty"`
	DevMsgID     *int         `json:"devMsgID"`
	MessageID    *int         `json:"messageID"`
	Config       *interface{} `json:"config,omitempty"`
	LogConfig    *interface{} `json:"logConfig,omitempty"`
	Version      string       `json:"version,omitempty"`
	WifiConfig   *interface{} `json:"wifiConfig,omitempty"`
	GroupConfig  *interface{} `json:"groupConfig,omitempty"`
	NodeConfig   *interface{} `json:"nodeConfig,omitempty"`
	ModuleConfig *interface{} `json:"moduleConfig,omitempty"`

	DevType   *int   `json:"devType"`
	DevModel  string `json:"devModel,omitempty"`
	Size      *int   `json:"size"`
	Sign      string `json:"sign,omitempty"`
	SignType  string `json:"signType,omitempty"`
	RebootDev *int   `json:"rebootDev"`
	Url       string `json:"url,omitempty"`
}

type UpgradeData struct {
	Size    *int   `json:"size"`
	Version string `json:"version,omitempty"`
	Url     string `json:"url,omitempty"`
	Md5     string `json:"md5,omitempty"`
}

type UpgradeMsg struct {
	Code   string      `json:"code,omitempty"`
	Data   UpgradeData `json:"data,omitempty"`
	Method string      `json:"method,omitempty"`
}

type ConfigResult struct {
	Result string `json:"result,omitempty"`
}

type NodeConfigResult struct {
	CanID  int    `json:"canID"`
	NodeSN string `json:"nodeSN,omitempty"`
	Result string `json:"result,omitempty"`
}

type ModConfigResult struct {
	PortID int    `json:"portID"`
	Result string `json:"result,omitempty"`
}

type ModuleConfigResult struct {
	ModuleConfigList []ModConfigResult `json:"moduleConfigList,omitempty"`
	CanID            int               `json:"canID"`
	NodeSN           string            `json:"nodeSN,omitempty"`
}

type DownMsgRsp struct {
	DevSN        string               `json:"devSN,omitempty"`
	Method       string               `json:"method,omitempty"`
	DevOption    string               `json:"devOption,omitempty"`
	Maintenance  string               `json:"maintenance,omitempty"`
	DevMsgID     int                  `json:"devMsgID"`
	MessageID    int                  `json:"messageID"`
	Result       string               `json:"result,omitempty"`
	RespCode     int                  `json:"respCode"`
	Config       ConfigResult         `json:"config,omitempty"`
	LogConfig    ConfigResult         `json:"logConfig,omitempty"`
	Version      string               `json:"version,omitempty"`
	WifiConfig   ConfigResult         `json:"wifiConfig,omitempty"`
	GroupConfig  ConfigResult         `json:"groupConfig,omitempty"`
	NodeConfig   []NodeConfigResult   `json:"nodeConfig,omitempty"`
	ModuleConfig []ModuleConfigResult `json:"moduleConfig,omitempty"`

	DevModel string `json:"devModel,omitempty"`
}

// mqtt0001
type PropertyUp struct {
	MsgId  string      `json:"msgid"`
	Params interface{} `json:"params"`
}
