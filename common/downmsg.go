package common

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"demo/config"
	cmap "github.com/orcaman/concurrent-map"
)

var DeviceTimeStampMap cmap.ConcurrentMap

func InitTimeStampMap() {
	DeviceTimeStampMap = cmap.New()
}

func ParseDownMsg(payload []byte) DownMsg {
	downMsg := DownMsg{}
	err := json.Unmarshal(payload, &downMsg)
	if err != nil {
		return DownMsg{}
	}
	return downMsg
}

func ParseUpgradeMsg(payload []byte) UpgradeMsg {
	upgradeMsg := UpgradeMsg{}
	err := json.Unmarshal(payload, &upgradeMsg)
	if err != nil {
		return UpgradeMsg{}
	}
	return upgradeMsg
}

func ProcDownMsg(c mqtt.Client, m mqtt.Message) {
	replyTopic := strings.Replace(m.Topic(), "down_raw", "up_raw", 1)
	downMsg := ParseDownMsg(m.Payload())
	downMsgRsp := DownMsgRsp{}
	downMsgRsp.DevSN = downMsg.DevSN
	if downMsg.DevMsgID != nil {
		downMsgRsp.DevMsgID = *downMsg.DevMsgID
	}
	if downMsg.MessageID != nil {
		downMsgRsp.MessageID = *downMsg.MessageID
	}
	if downMsg.Method != "" {
		method := downMsg.Method
		if method[len(method)-3:] == "Rsp" {
			// log.Println("Recv response from cloud, do nothing")
			return
		} else {
			switch method {
			case SYNC_TIME:
				downMsgRsp.Result = "success"
			case IOT_LOG_CFG_PUSH:
				if downMsg.LogConfig != nil {
					downMsgRsp.LogConfig = ConfigRsult{Result: "success"}
				}
				if downMsg.Version != "" {
					downMsgRsp.Version = downMsg.Version
				}
			case SET_IOT_GW_CFG:
				if downMsg.Config != nil {
					config := *downMsg.Config
					cfg := config.(map[string]interface{})
					DeviceTimeStampMap.Set(downMsg.DevSN, cfg["timeStamp"].(string))
					downMsgRsp.Config = ConfigRsult{Result: "success"}
				}
				if downMsg.GroupConfig != nil {
					downMsgRsp.GroupConfig = ConfigRsult{Result: "success"}
				}
			case SET_NET_CFG:
				if downMsg.WifiConfig != nil {
					downMsgRsp.WifiConfig = ConfigRsult{Result: "success"}
				}
			}
			downMsgRsp.RespCode = 0
			downMsgRsp.Method = method + "Rsp"
		}
	} else if downMsg.DevOption != "" {
		devOption := downMsg.DevOption
		if devOption[len(devOption)-3:] == "Rsp" {
			// log.Println("Recv response from cloud, do nothing")
			return
		} else {
			switch devOption {
			case DEV_UPGRADE:
				downMsgRsp.DevModel = downMsg.DevModel
				downMsgRsp.Version = downMsg.Version
				downMsgRsp.Result = "success"

				<-time.After(3 * time.Second)
				payload := EncUpMsg(DEV_UPGRADE_PROGRESS_UP, downMsg.DevSN, downMsg)
				if c.IsConnectionOpen() {
					if config.MQTT_PUBLISH_ENABLE {
						log.Printf("发布topic为 %s 的消息: %s", replyTopic, string(payload))
					}
					c.Publish(replyTopic, 0x00, false, string(payload))
				}
				<-time.After(3 * time.Second)
				payload = EncUpMsg(BASIC_INFO_UP, downMsg.DevSN, downMsg)
				if c.IsConnectionOpen() {
					if config.MQTT_PUBLISH_ENABLE {
						log.Printf("发布topic为 %s 的消息: %s", replyTopic, string(payload))
					}
					c.Publish(replyTopic, 0x00, false, string(payload))
				}
			}
			downMsgRsp.RespCode = 0
			downMsgRsp.Method = devOption + "Rsp"
		}
	}
	payload, _ := json.Marshal(downMsgRsp)
	if token := c.Publish(replyTopic, 0, false, string(payload)); token.Wait() && token.Error() != nil {
		log.Printf("reply to topic %v \n payload %v \n err = %v", replyTopic, string(payload), token.Error())
	} else {
		if config.MQTT_PUBLISH_ENABLE {
			log.Printf("发布topic为 %s 的消息: %s", replyTopic, string(payload))
		}
	}
}

func ProcUpgradeMsg(c mqtt.Client, m mqtt.Message) {
	replyTopic := strings.Replace(m.Topic(), "upgrade", "progress", 1)
	upgradeMsg := ParseUpgradeMsg(m.Payload())
	if upgradeMsg.Method == "/ota/device/upgrade" {
		devSN := replyTopic[len(replyTopic)-config.DEVICE_SN_LEN:]
		progressMsg := ProgressMsg{
			Id: 1,
			Params: Param{
				Step: "1",
				Desc: "upgrading",
			},
		}
		payload, _ := json.Marshal(progressMsg)
		if c.IsConnectionOpen() {
			if config.MQTT_PUBLISH_ENABLE {
				log.Printf("发布topic为 %s 的消息: %s", replyTopic, string(payload))
			}
			c.Publish(replyTopic, 0, false, string(payload))
		}
		timer := time.NewTimer(time.Second * 3)
		<-timer.C
		progressMsg.Id = 2
		progressMsg.Params.Step = "50"
		payload, _ = json.Marshal(progressMsg)
		if c.IsConnectionOpen() {
			if config.MQTT_PUBLISH_ENABLE {
				log.Printf("发布topic为 %s 的消息: %s", replyTopic, string(payload))
			}
			c.Publish(replyTopic, 0, false, string(payload))
		}
		timer.Reset(time.Second * 3)
		<-timer.C
		progressMsg.Id = 3
		progressMsg.Params.Step = "100"
		payload, _ = json.Marshal(progressMsg)
		if c.IsConnectionOpen() {
			if config.MQTT_PUBLISH_ENABLE {
				log.Printf("发布topic为 %s 的消息: %s", replyTopic, string(payload))
			}
			c.Publish(replyTopic, 0, false, string(payload))
		}
		timer.Reset(time.Second * 3)
		<-timer.C
		payload = EncUpMsg(BASIC_INFO_UP, devSN, DownMsg{Version: upgradeMsg.Data.Version, DevModel: "T320M"})
		if c.IsConnectionOpen() {
			if config.MQTT_PUBLISH_ENABLE {
				log.Printf("发布topic为 %s 的消息: %s", replyTopic, string(payload))
			}
			topic := fmt.Sprintf(config.PUB_TOPIC[0], config.PRODUCT_KEY, devSN)
			c.Publish(topic, 0, false, string(payload))
		}
	}
}
