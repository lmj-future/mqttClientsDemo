package common

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	mathRand "math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lmj/mqtt-clients-demo/config"
	"github.com/lmj/mqtt-clients-demo/logger"
	"github.com/pborman/uuid"
)
var DevSNwithMac = &sync.Map{}  //记录mac信息

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
		randomMMAC, _ := DevSNwithMac.Load(devSN + "M")
		modMac := randomMMAC.(string)
		modMac = modMac[0:4] + "-" + modMac[4:8] + "-" + modMac[8:]
		randomGMAC1, _ := DevSNwithMac.Load(devSN + "G")
		randomGMAC2 := randomGMAC1.(string)
		randomGMAC2 = randomGMAC2[0:4]+ "-"+randomGMAC2[4:8] +"-" + randomGMAC2[8:12] + "-" + randomGMAC2[8:12]
		randomIP := make([]byte, 4)
		rand.Read(randomIP)
		ip := net.IP(randomIP)
		nodeModel := config.PRODUCT_NAME
		nodeVersion := "R1268"
		modModel := "T301-R"
		modVersion := "R1245"
		if downMsg.DevModel != "" {
			switch downMsg.DevModel {
			case config.PRODUCT_NAME:
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
				Gmac:                 randomGMAC1.(string),
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
						Type:           3,
						PortID:         1,
						ModModel:       "T301-Z",
						ModHwVersion:   "Ver.A",
						ModSwVersion:   "R1245",
						ModSwInVersion: "V100R001B01D029SP45",
						ModMAC:         randomGMAC2,
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
				NodeModel:  config.PRODUCT_NAME,
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
						ModStatus:  1,
						Type:       3,
						ModModel:   "T301-Z",
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
				NodeModel: config.PRODUCT_NAME,
				ModTopoInfoList: []ModTopoInfo{
					{
						PortID:    0,
						ModStatus: 1,
						Type:      2,
						ModModel:  "T301-R",
					},
					{
						PortID:    1,
						ModStatus: 1,
						Type:      3,
						ModModel:  "T301-Z",
					},
				},
			},
		}
	}
	upMsgByte, _ := json.Marshal(upMsg)
	return upMsgByte
}

func upRawWhenConnectByT320(clientInfo MqttClientInfo) {
	topic := fmt.Sprintf(config.PUB_TOPIC[0], config.PRODUCT_KEY, clientInfo.DevSN)
	<-time.After(time.Second)
	payload := EncUpMsg(DEV_CONNECT_STATUS, clientInfo.DevSN, DownMsg{})
	if clientInfo.Client.IsConnectionOpen() {
		if config.LOG_MQTT_PUBLISH_ENABLE {
			logger.Log.Infof("发布topic为 %s 的消息: %s", topic, string(payload))
		}
		clientInfo.Client.Publish(topic, 0x00, false, string(payload))
		<-time.After(time.Microsecond * 10)
	}
	if clientInfo.Client.IsConnectionOpen() {
		payload = EncUpMsg(IOT_VER_INFO_GET, clientInfo.DevSN, DownMsg{})
		if config.LOG_MQTT_PUBLISH_ENABLE {
			logger.Log.Infof("发布topic为 %s 的消息: %s", topic, string(payload))
		}
		clientInfo.Client.Publish(topic, 0x00, false, string(payload))
		<-time.After(time.Microsecond * 2000)
	}
	if clientInfo.Client.IsConnectionOpen() {
		payload = EncUpMsg(BASIC_INFO_UP, clientInfo.DevSN, DownMsg{})
		if config.LOG_MQTT_PUBLISH_ENABLE {
			logger.Log.Infof("发布topic为 %s 的消息: %s", topic, string(payload))
		}
		clientInfo.Client.Publish(topic, 0x00, false, string(payload))
		<-time.After(time.Microsecond * 3000)
	}
	if clientInfo.Client.IsConnectionOpen() {
		payload = EncUpMsg(STATUS_UP, clientInfo.DevSN, DownMsg{})
		if config.LOG_MQTT_PUBLISH_ENABLE {
			logger.Log.Infof("发布topic为 %s 的消息: %s", topic, string(payload))
		}
		clientInfo.Client.Publish(topic, 0x00, false, string(payload))
	}
}
func upRawWhenConnectByMqtt0001(clientInfo MqttClientInfo) {
	// devSN := strings.Replace(clientInfo.ClientId, "&v5", "", 1)
	// topic := fmt.Sprintf(config.PUB_TOPIC[0], config.PRODUCT_KEY, devSN)

	// do nothing
}

func UpRawWhenConnect(clientInfo MqttClientInfo) {
	switch config.PRODUCT_NAME {
	case "T320M", "T320MX", "T320MX-U":
		upRawWhenConnectByT320(clientInfo)
	case "示例产品-mqtt":
		upRawWhenConnectByMqtt0001(clientInfo)
	}
}

func upRawAfterConnectByT320(clientInfo MqttClientInfo, timeSleep time.Duration) {
	topic := fmt.Sprintf(config.PUB_TOPIC[0], config.PRODUCT_KEY, clientInfo.DevSN)
	payload := EncUpMsg(NODE_TOPO_INFO_UP, clientInfo.DevSN, DownMsg{})
	if clientInfo.Client.IsConnectionOpen() {
		if config.LOG_MQTT_PUBLISH_ENABLE {
			logger.Log.Infof("发布topic为 %s 的消息: %s", topic, string(payload))
		}
		clientInfo.Client.Publish(topic, 0x00, false, string(payload))
		<-time.After(timeSleep / time.Duration(config.PPS))
	}
	for i := 1; i < config.PPS; i++ {
		if clientInfo.Client.IsConnectionOpen() {
			payload = EncUpMsg(IOT_GW_CFG_SYNC, clientInfo.DevSN, DownMsg{})
			if config.LOG_MQTT_PUBLISH_ENABLE {
				logger.Log.Infof("发布topic为 %s 的消息: %s", topic, string(payload))
			}
			clientInfo.Client.Publish(topic, 0x00, false, string(payload))
			<-time.After(timeSleep / time.Duration(config.PPS))
		}
	}
}

func pubProperty(params map[string]interface{}, clientInfo MqttClientInfo, topic string, timeSleep time.Duration) {
	var msg map[string]interface{}
	json.Unmarshal([]byte(`{"msgid":"","params":{}}`), &msg)
	msg["msgid"] = strconv.Itoa(Counter.IncrementAndGet(clientInfo.DevSN))
	msg["params"] = params
	payload, _ := json.Marshal(msg)
	if config.LOG_MQTT_PUBLISH_ENABLE {
		logger.Log.Infof("发布topic为 %s 的消息: %s", topic, string(payload))
	}
	if clientInfo.Client.IsConnectionOpen() {
		clientInfo.Client.Publish(topic, 0x00, false, string(payload))
	}
	<-time.After(timeSleep / time.Duration(config.PPS))
}
func upRawAfterConnectByMqtt0001(clientInfo MqttClientInfo, timeSleep time.Duration) {
	if config.PROPERTY_UP_ENABLE {
		topic := fmt.Sprintf(config.PUB_TOPIC[0], config.PRODUCT_KEY, clientInfo.DevSN)
		// 采用随机种子，对以下数据生成随机值
		mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
		if config.BOOLKey {
			var params map[string]interface{} = map[string]interface{}{
				"BOOLKey": map[string]bool{"value": mathRand.Intn(2) == 1},
			}
			pubProperty(params, clientInfo, topic, timeSleep)
		}
		if config.DoubleKey {
			var params map[string]interface{} = map[string]interface{}{
				"DoubleKey": map[string]float64{"value": float64(1+mathRand.Intn(254)) + mathRand.Float64()},
			}
			pubProperty(params, clientInfo, topic, timeSleep)
		}
		if config.LongKey {
			var params map[string]interface{} = map[string]interface{}{
				"LongKey": map[string]int64{"value": int64(1 + mathRand.Intn(254))},
			}
			pubProperty(params, clientInfo, topic, timeSleep)
		}
		if config.StringKey {
			var params map[string]interface{} = map[string]interface{}{
				"StringKey": map[string]string{"value": uuid.NewRandom().String()},
			}
			pubProperty(params, clientInfo, topic, timeSleep)
		}
		if config.PowerSwitch {
			var params map[string]interface{} = map[string]interface{}{
				"PowerSwitch": map[string]bool{"value": mathRand.Intn(2) == 1},
			}
			pubProperty(params, clientInfo, topic, timeSleep)
		}
		if config.CustomKey {
			var customValue map[string]interface{}
			json.Unmarshal([]byte(config.CustomValue), &customValue)
			var params map[string]interface{} = make(map[string]interface{})
			for k, v := range customValue {
				params[k] = map[string]interface{}{"value": v}
			}
			pubProperty(params, clientInfo, topic, timeSleep)
		}
	}
}
func getPubTopic(devSN string) []string {
	var pubTopic []string
	for _, v := range config.PUB_TOPIC {
		if v != "" {
			if strings.Contains(v, "{devSN}") {
				v = strings.ReplaceAll(v, "{devSN}", devSN)
			}
			if strings.Contains(v, "{devKey}") {
				v = strings.ReplaceAll(v, "{devKey}", config.DEVICE_KEY)
			}
			pubTopic = append(pubTopic, v)
		}
	}
	return pubTopic
}
func upRawAfterConnectByCustom(clientInfo MqttClientInfo, timeSleep time.Duration) {
	if config.PROPERTY_UP_ENABLE {
		topic := getPubTopic(clientInfo.DevSN)[0]
		if config.CustomKey {
			var customValue map[string]interface{}
			json.Unmarshal([]byte(config.CustomValue), &customValue)
			payload, _ := json.Marshal(customValue)
			for i := 0; i < config.PPS; i++ {
				if config.LOG_MQTT_PUBLISH_ENABLE {
					logger.Log.Infof("发布topic为 %s 的消息: %s", topic, string(payload))
				}
				if clientInfo.Client.IsConnectionOpen() {
					clientInfo.Client.Publish(topic, 0x00, false, string(payload))
				}
				<-time.After(timeSleep / time.Duration(config.PPS))
			}
		}
	}
}

func UpRawAfterConnect(clientInfo MqttClientInfo, timeSleep time.Duration) {
	switch config.PRODUCT_NAME {
	case "T320M", "T320MX", "T320MX-U":
		upRawAfterConnectByT320(clientInfo, timeSleep)
	case "示例产品-mqtt":
		upRawAfterConnectByMqtt0001(clientInfo, timeSleep)
	case "自定义":
		upRawAfterConnectByCustom(clientInfo, timeSleep)
	}
}
