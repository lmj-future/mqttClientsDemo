package common

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"demo/config"
)

func EncUpMsg(method string, devSN string, downMsg DownMsg) []byte {
	upMsg := UpMsg{}
	upMsg.DevSN = devSN
	upMsg.Method = method
	devMsgID := Counter.IncrementAndGet(devSN)
	upMsg.DevMsgID = &devMsgID
	switch method {
	case IOT_NET_CFG_UP:
		type wifiConfig struct {
			CountryCode int    `json:"countryCode"`
			SSID        string `json:"SSID"`
			Password    string `json:"password"`
		}
		config := wifiConfig{
			CountryCode: 2510403,
			SSID:        "H3C-T320M-Config",
			Password:    "h3ciotadmin",
		}
		upMsg.WifiConfig = config
	case IOT_GW_CFG_SYNC:
		type gwConfig struct {
			TimeStamp string `json:"timeStamp"`
		}
		timeStamp := strconv.FormatInt(time.Now().UnixNano(), 10)
		if v, ok := DeviceTimeStampMap.Get(devSN); ok {
			timeStamp = v.(string)
		}
		config := gwConfig{
			TimeStamp: timeStamp,
		}
		upMsg.Config = config
	case IOT_GROUP_CFG_SYNC:
		type groupConfig struct {
			Type      int    `json:"type"`
			GroupID   int    `json:"gropuID"`
			TimeStamp string `json:"timeStamp"`
		}
		config := []groupConfig{
			{
				Type:      1,
				GroupID:   1,
				TimeStamp: strconv.FormatInt(time.Now().UnixNano(), 10),
			},
			{
				Type:      2,
				GroupID:   2,
				TimeStamp: strconv.FormatInt(time.Now().UnixNano(), 10),
			},
			{
				Type:      3,
				GroupID:   3,
				TimeStamp: strconv.FormatInt(time.Now().UnixNano(), 10),
			},
			{
				Type:      4,
				GroupID:   4,
				TimeStamp: strconv.FormatInt(time.Now().UnixNano(), 10),
			},
		}
		upMsg.GroupConfig = config
	case IOT_NODE_SYNC:
		type nodeConfig struct {
			NodeSN    string `json:"nodeSN"`
			CanID     int    `json:"canID"`
			TimeStamp string `json:"timeStamp"`
		}
		config := []nodeConfig{
			{
				NodeSN:    devSN,
				CanID:     0,
				TimeStamp: strconv.FormatInt(time.Now().UnixNano(), 10),
			},
		}
		upMsg.NodeConfig = config
	case IOT_MOD_CFG_SYNC:
		type modConfig struct {
			PortID    int    `json:"portID"`
			Type      int    `json:"type"`
			ConfigSrc int    `json:"configSrc"`
			TimeStamp string `json:"timeStamp"`
		}
		type moduleConfig struct {
			NodeSN           string      `json:"nodeSN"`
			CanID            int         `json:"canID"`
			ModuleConfigList []modConfig `json:"moduleConfigList"`
		}
		config := []moduleConfig{
			{
				NodeSN: devSN,
				CanID:  0,
				ModuleConfigList: []modConfig{
					{
						PortID:    0,
						Type:      2,
						ConfigSrc: 2,
						TimeStamp: strconv.FormatInt(time.Now().UnixNano(), 10),
					},
				},
			},
		}
		upMsg.ModuleConfig = config
	case IOT_LOG_CFG_GET:
		type logConfig struct {
			TimeStamp string `json:"timeStamp"`
		}
		config := logConfig{
			TimeStamp: strconv.FormatInt(time.Now().UnixNano(), 10),
		}
		upMsg.LogConfig = config
		upMsg.Version = "1.0"
	case DEV_UPGRADE_PROGRESS_UP:
		upMsg.Method = ""
		upMsg.DevOption = method
		upMsg.DevModel = downMsg.DevModel
		upMsg.Version = downMsg.Version
		step := 100
		upMsg.Step = &step
	case DEV_CONNECT_STATUS:
		connectStatus := 0
		compatibilityVer := 3
		upMsg.ConnectStatus = &connectStatus
		upMsg.CompatibilityVer = &compatibilityVer
	case IOT_VER_INFO_GET:
		upMsg.Version = "1.0"
		upMsg.DevModelList = []string{"T301-IR", "T301-R", "T301-Z", "T320"}
	case BASIC_INFO_UP:
		randomMAC := make([]byte, 6)
		rand.Read(randomMAC)
		nodeMac := hex.EncodeToString(randomMAC)
		nodeMac = nodeMac[0:2] + ":" + nodeMac[2:4] + ":" + nodeMac[4:6] + ":" + nodeMac[6:8] + ":" + nodeMac[8:10] + ":" + nodeMac[10:]
		rand.Read(randomMAC)
		modMac := hex.EncodeToString(randomMAC)
		modMac = modMac[0:4] + "-" + modMac[4:8] + "-" + modMac[8:]
		randomIP := make([]byte, 4)
		rand.Read(randomIP)
		ip := net.IP(randomIP)
		randomGMAC := make([]byte, 8)
		rand.Read(randomGMAC)
		nodeModel := "T320M"
		nodeVersion := "R1268"
		modModel := "T301-R"
		modVersion := "R1245"
		if downMsg.DevModel != "" {
			switch downMsg.DevModel {
			case "T320M":
				nodeModel = downMsg.DevModel
				nodeVersion = downMsg.Version
			case "T301-R":
				modModel = downMsg.DevModel
				modVersion = downMsg.Version
			}
		}
		upMsg.NodeBasicInfoList = []NodeBasicInfo{
			{
				NodeMAC:              nodeMac,
				NodeIP:               ip.String(),
				Gmac:                 hex.EncodeToString(randomGMAC),
				CanID:                0,
				NodeSN:               devSN,
				NodeModel:            nodeModel,
				NodeHwVersion:        "Ver.B",
				NodeSwVersion:        nodeVersion,
				NodeLastRebootReason: 5,
				CanAddr:              0,
				ModBasicInfoList: []ModBasicInfo{
					{
						Type:           2,
						PortID:         0,
						ModModel:       modModel,
						ModHwVersion:   "Ver.A",
						ModSwVersion:   modVersion,
						ModSwInVersion: "V100R001B01D029SP45",
						ModMAC:         modMac,
						ModSN:          devSN,
					},
					{
						Type:           6,
						PortID:         1,
						ModModel:       "T300PB0U",
						ModHwVersion:   "Ver.A",
						ModSwVersion:   "R1245",
						ModSwInVersion: "V100R001B01D029SP45",
						ModMAC:         modMac,
						ModSN:          devSN,
					},
				},
			},
		}
	case STATUS_UP:
		upMsg.NodeStatusList = []NodeStatus{
			{
				CanID:      0,
				NodeSN:     devSN,
				NodeModel:  "T320M",
				NodeStatus: 5,
				NodeAging:  0,
				ModStatusList: []ModStatus{
					{
						PortID:     0,
						ModStatus:  1,
						Type:       2,
						ModModel:   "T301-R",
						ModPullCfg: 1,
					},
					{
						PortID:     1,
						ModStatus:  0,
						Type:       6,
						ModModel:   "T300PB0U",
						ModPullCfg: 0,
					},
				},
			},
		}
	case NODE_TOPO_INFO_UP:
		upMsg.NodeTopoInfoList = []NodeTopoInfo{
			{  
				CanID:     0,
				NodeSN:    devSN,
				NodeModel: "T320M",
				ModTopoInfoList: []ModTopoInfo{
					{
						PortID:    0,
						ModStatus: 1,
						Type:      2,
						ModModel:  "T301-R",
					},
					{
						PortID:    1,
						ModStatus: 0,
						Type:      6,
						ModModel:  "T300PB0U",
					},
				},
			},
		}
	}
	upMsgByte, _ := json.Marshal(upMsg)
	// log.Println("upMsg: ", string(upMsgByte))
	return upMsgByte
}

func UpRawWhenConnect(clientInfo MqttClientInfo) {
	devSN := strings.Replace(clientInfo.ClientId, config.PRODUCT_KEY+"&", "", 1)
	topic := fmt.Sprintf(config.PUB_TOPIC[0], config.PRODUCT_KEY, devSN)
	<-time.After(time.Second)
	payload := EncUpMsg(DEV_CONNECT_STATUS, devSN, DownMsg{})
	if clientInfo.Client.IsConnectionOpen() {
		if config.MQTT_PUBLISH_ENABLE {
			log.Printf("发布topic为 %s 的消息: %s", topic, string(payload))
		}
		clientInfo.Client.Publish(topic, 0x00, false, string(payload))
		<-time.After(time.Microsecond * 500)
	}
	if clientInfo.Client.IsConnectionOpen() {
		payload = EncUpMsg(IOT_VER_INFO_GET, devSN, DownMsg{})
		if config.MQTT_PUBLISH_ENABLE {
			log.Printf("发布topic为 %s 的消息: %s", topic, string(payload))
		}
		clientInfo.Client.Publish(topic, 0x00, false, string(payload))
		<-time.After(time.Microsecond * 100)
	}
	if clientInfo.Client.IsConnectionOpen() {
		payload = EncUpMsg(BASIC_INFO_UP, devSN, DownMsg{})
		if config.MQTT_PUBLISH_ENABLE {
			log.Printf("发布topic为 %s 的消息: %s", topic, string(payload))
		}
		clientInfo.Client.Publish(topic, 0x00, false, string(payload))
		<-time.After(time.Microsecond * 500)
	}
	if clientInfo.Client.IsConnectionOpen() {
		payload = EncUpMsg(STATUS_UP, devSN, DownMsg{})
		if config.MQTT_PUBLISH_ENABLE {
			log.Printf("发布topic为 %s 的消息: %s", topic, string(payload))
		}
		clientInfo.Client.Publish(topic, 0x00, false, string(payload))
	}
}

func UpRawAfterConnect(clientInfo MqttClientInfo, timeSleep time.Duration) {
	devSN := strings.Replace(clientInfo.ClientId, config.PRODUCT_KEY+"&", "", 1)
	topic := fmt.Sprintf(config.PUB_TOPIC[0], config.PRODUCT_KEY, devSN)
	payload := EncUpMsg(NODE_TOPO_INFO_UP, devSN, DownMsg{})
	if clientInfo.Client.IsConnectionOpen() {
		if config.MQTT_PUBLISH_ENABLE {
			log.Printf("发布topic为 %s 的消息: %s", topic, string(payload))
		}
		clientInfo.Client.Publish(topic, 0x00, false, string(payload))
		<-time.After(timeSleep / time.Duration(config.PPS))
	}
	for i := 1; i < config.PPS; i++ {
		if clientInfo.Client.IsConnectionOpen() {
			payload = EncUpMsg(IOT_GW_CFG_SYNC, devSN, DownMsg{})
			if config.MQTT_PUBLISH_ENABLE {
				log.Printf("发布topic为 %s 的消息: %s", topic, string(payload))
			}
			clientInfo.Client.Publish(topic, 0x00, false, string(payload))
			<-time.After(timeSleep / time.Duration(config.PPS))
		}
	}
}
