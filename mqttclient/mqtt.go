package mqttclient

import (
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lmj/mqtt-clients-demo/common"
	"github.com/lmj/mqtt-clients-demo/config"
	"github.com/lmj/mqtt-clients-demo/logger"
	cmap "github.com/orcaman/concurrent-map"
)

var ClientMap cmap.ConcurrentMap
var ClientStopMap cmap.ConcurrentMap

func Init() {
	ClientMap = cmap.New()
	ClientStopMap = cmap.New()
}

// 创建全局mqtt publish消息处理 handler
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	logger.Log.Warnf("订阅到一个未知的topic(谁在搞我?): %s, msg : %s", msg.Topic(), msg.Payload())
}

func connect(clientInfo common.MqttClientInfo, clientOptions *mqtt.ClientOptions, clientState *string) (mqtt.Client, error) {
	// 创建客户端连接
	client := mqtt.NewClient(clientOptions)
	// 客户端连接判断
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		*clientState = token.Error().Error()
		if token.Error().Error() == "not Authorized" {
			logger.Log.Warnf("MQTT连接被拒绝(设备[%s]没有加到平台吧?用户名[%s]密码[%s]填对了吗?)", clientInfo.DevSN, clientOptions.Username, clientOptions.Password)
			config.MQTT_CLIENT_CONNECT_INTERVAL = 10000
		} else {
			logger.Log.Warnf("MQTT连接时发生错误(设备[%s]连接地址%s看看对了吗?): %s", clientInfo.DevSN, clientOptions.Servers, token.Error())
		}
		return client, token.Error()
	} else {
		*clientState = "ok"
		logger.Log.Infof("MQTT连接成功[%s]%s", clientInfo.DevSN, clientOptions.Servers)
		return client, nil
	}
}
func connectLoop(clientInfo common.MqttClientInfo, clientOptions *mqtt.ClientOptions, subTopic []string,
	wg *sync.WaitGroup, reconnect bool, timestamp string, clientState *string) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Log.Errorln("连接MQTT时捕获到一个未知错误(你不检查一下?): ", err)
		}
	}()
	reconnectCount := 0
	// 死循环用于连接失败后进行重连，连接成功后跳出循环
	for {
		// 如果此次开始被停止了，那么就要停止此次开始所做的事情
		if v, ok := ClientStopMap.Get(timestamp); (ok && v.(bool)) || !ok {
			break
		}
		if cli, err := connect(clientInfo, clientOptions, clientState); err != nil {
			// 连接失败，释放连接
			// cli.Disconnect(250)
			if reconnectCount == 0 {
				clientInfo.Client = cli
				// 客户端信息存入Map
				ClientMap.Set(clientInfo.DevSN, clientInfo)
			}
			reconnectCount++

			if reconnectCount > config.MQTT_CLIENT_RECONNECT_COUNT {
				reconnectCount = 0
				// 当重连次数超过MQTT_CLIENT_RECONNECT_COUNT时，需要休眠一下，停顿MQTT_CLIENT_SLEEP_INTERVAL秒后进行下一次重连
				<-time.After(time.Second * time.Duration(config.MQTT_CLIENT_SLEEP_INTERVAL))
			} else {
				// 停顿MQTT_CLIENT_RECONNECT_INTERVAL秒后进行下一次重连
				<-time.After(time.Second * time.Duration(config.MQTT_CLIENT_RECONNECT_INTERVAL))
			}
		} else {
			reconnectCount = 0
			if reconnect {
				// 重连成功后，进行消息订阅
				for _, t := range subTopic {
					cli.Subscribe(t, 0, func(c mqtt.Client, msg mqtt.Message) {
						if config.LOG_MQTT_SUBSCRIBE_ENABLE {
							logger.Log.Infof("订阅到topic为 %s 的消息: %s", msg.Topic(), msg.Payload())
						}
						if strings.Contains(msg.Topic(), "/thing/model/down_raw") {
							go common.ProcDownMsg(c, msg)
						} else if strings.Contains(msg.Topic(), "/ota/device/upgrade") {
							go common.ProcUpgradeMsg(c, msg)
						}
					})
				}
			}
			clientInfo.Client = cli
			// 连接状态置成不在连接中
			clientInfo.Connectting = false
			// 客户端信息存入Map
			ClientMap.Set(clientInfo.DevSN, clientInfo)
			if wg != nil {
				wg.Done()
			}
			break
		}
	}
}

// MQTT client connect
func Connect(devSN string, broker string, userName string, password string, clientId string, keepAlive time.Duration,
	subTopic []string, wg *sync.WaitGroup, timestamp string, clientState *string) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Log.Errorln("连接MQTT时捕获到一个未知错误(你不检查一下?): ", err)
		}
	}()
	// 设置连接参数
	clientOptions := mqtt.NewClientOptions().AddBroker(broker).SetUsername(userName).SetPassword(password)
	// 设置客户端ID
	clientOptions.SetClientID(clientId)
	// 设置保活时长
	clientOptions.SetKeepAlive(keepAlive)
	// 设置handler
	clientOptions.SetDefaultPublishHandler(messagePubHandler)
	// 设置MQTT协议版本号
	clientOptions.SetProtocolVersion(4)
	clientOptions.SetAutoReconnect(false)
	clientInfo := common.MqttClientInfo{
		DevSN:       devSN,
		Connectting: true,
	}
	// 开始尝试连接循环，直到连接成功
	connectLoop(clientInfo, clientOptions, subTopic, wg, false, timestamp, clientState)

	// 启动协程，用于处理重连
	go func() {
		for {
			// 如果此次开始被停止了，那么就要停止此次开始所做的事情
			if v, ok := ClientStopMap.Get(timestamp); ok && v.(bool) {
				break
			}
			if !ClientMap.IsEmpty() {
				// 如果不是正在连接状态，并且连接是断开的，那么将连接状态置成连接中，并释放上次连接，开始下一次连接循环
				if v, ok := ClientMap.Get(devSN); ok && !v.(common.MqttClientInfo).Connectting && !v.(common.MqttClientInfo).Client.IsConnectionOpen() {
					info := v.(common.MqttClientInfo)
					logger.Log.Warnf("设备MQTT连接断开了[%s]，需要进行重连", info.DevSN)
					// 将连接状态置成正在连接中
					info.Connectting = true
					// 释放连接
					// info.Client.Disconnect(250)
					connectLoop(info, clientOptions, subTopic, nil, true, timestamp, clientState)
				}
			}
			// 停顿两秒，避免频繁重连
			<-time.After(2 * time.Second)
		}
	}()
}

func Subscribe(topic []string, client mqtt.Client) {
	// 连接成功后，进行消息订阅
	for _, t := range topic {
		client.Subscribe(t, 0, func(c mqtt.Client, msg mqtt.Message) {
			if config.LOG_MQTT_SUBSCRIBE_ENABLE {
				logger.Log.Infof("订阅到topic为 %s 的消息: %s", msg.Topic(), msg.Payload())
			}
			if strings.Contains(msg.Topic(), "/thing/model/down_raw") {
				go common.ProcDownMsg(c, msg)
			} else if strings.Contains(msg.Topic(), "/ota/device/upgrade") {
				go common.ProcUpgradeMsg(c, msg)
			} else if strings.Contains(msg.Topic(), "/sys/property/down") {
				go common.ProcDownMsg(c, msg)
			} else if strings.Contains(msg.Topic(), "/sys/service/invoke") {
				go common.ProcDownMsg(c, msg)
			}
		})
	}
}
