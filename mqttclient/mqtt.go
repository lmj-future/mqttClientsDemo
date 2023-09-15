package mqttclient

import (
	"log"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"demo/common"
	"demo/config"
	cmap "github.com/orcaman/concurrent-map"
)

var ClientMap cmap.ConcurrentMap

func Init() {
	ClientMap = cmap.New()
}

// 创建全局mqtt publish消息处理 handler
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("订阅到一个未知的topic(谁在搞我?): %s, msg : %s", msg.Topic(), msg.Payload())
}

func connect(clientInfo common.MqttClientInfo, clientOptions *mqtt.ClientOptions) (mqtt.Client, error) {
	// 创建客户端连接
	client := mqtt.NewClient(clientOptions)
	// 客户端连接判断
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		if token.Error().Error() == "not Authorized" {
			log.Printf("MQTT连接被拒绝(设备[%s]没有加到平台吧?)", clientInfo.ClientId)
			config.MQTT_CLIENT_CONNECT_INTERVAL = 10000
		} else {
			log.Printf("MQTT连接时发生错误(设备[%s]连接地址[%s]看看对了吗?): %s", clientInfo.ClientId, clientOptions.Servers[0], token.Error())
		}
		return client, token.Error()
	} else {
		log.Printf("MQTT连接成功[%s]", clientInfo.ClientId)
		return client, nil
	}
}
func connectLoop(clientInfo common.MqttClientInfo, clientOptions *mqtt.ClientOptions, subTopic []string, wg *sync.WaitGroup, reconnect bool) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("连接MQTT时捕获到一个未知错误(你不检查一下?): ", err)
		}
	}()
	reconnectCount := 0
	// 死循环用于连接失败后进行重连，连接成功后跳出循环
	for {
		if cli, err := connect(clientInfo, clientOptions); err != nil {
			// 连接失败，释放连接
			// cli.Disconnect(250)
			if reconnectCount == 0 {
				clientInfo.Client = cli
				// 客户端信息存入Map
				ClientMap.Set(clientInfo.ClientId, clientInfo)
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
						if config.MQTT_SUBSCRIBE_ENABLE {
							log.Printf("订阅到topic为 %s 的消息: %s", msg.Topic(), msg.Payload())
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
			ClientMap.Set(clientInfo.ClientId, clientInfo)
			if wg != nil {
				wg.Done()
			}
			break
		}
	}
}

// MQTT client connect
func Connect(broker string, userName string, password string, clientId string, keepAlive time.Duration, subTopic []string, wg *sync.WaitGroup) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("连接MQTT时捕获到一个未知错误(你不检查一下?): ", err)
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
		ClientId:    clientId,
		Connectting: true,
	}
	// 开始尝试连接循环，直到连接成功
	connectLoop(clientInfo, clientOptions, subTopic, wg, false)

	// 启动协程，用于处理重连
	go func() {
		for {
			if !ClientMap.IsEmpty() {
				// 如果不是正在连接状态，并且连接是断开的，那么将连接状态置成连接中，并释放上次连接，开始下一次连接循环
				if v, ok := ClientMap.Get(clientId); ok && !v.(common.MqttClientInfo).Connectting && !v.(common.MqttClientInfo).Client.IsConnectionOpen() {
					info := v.(common.MqttClientInfo)
					log.Printf("设备MQTT连接断开了[%s]，需要进行重连", info.ClientId)
					// 将连接状态置成正在连接中
					info.Connectting = true
					// 释放连接
					// info.Client.Disconnect(250)
					connectLoop(info, clientOptions, subTopic, nil, true)
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
			if config.MQTT_SUBSCRIBE_ENABLE {
				log.Printf("订阅到topic为 %s 的消息: %s", msg.Topic(), msg.Payload())
			}
			if strings.Contains(msg.Topic(), "/thing/model/down_raw") {
				go common.ProcDownMsg(c, msg)
			} else if strings.Contains(msg.Topic(), "/ota/device/upgrade") {
				go common.ProcUpgradeMsg(c, msg)
			}
		})
	}
}

// Publish Publish
func Publish(clientId string, topic string, msg string) {
	if cli, ok := ClientMap.Get(clientId); ok {
		client := cli.(common.MqttClientInfo).Client
		// 发布消息
		token := client.Publish(topic, 0, false, msg)
		if !token.WaitTimeout(60 * time.Second) {
			log.Printf("[MQTT] Publish timeout: topic: %s, mqttMsg %s", topic, msg)
			client.Disconnect(250)
			<-time.After(time.Second)
			return
		}
		if token.Error() != nil {
			log.Printf("[MQTT] Publish error: topic: %s, mqttMsg %s, error %s", topic, msg, token.Error().Error())
			client.Disconnect(250)
			<-time.After(time.Second)
		} else {
			log.Printf("[MQTT] Publish success: topic: %s, mqttMsg %s", topic, msg)
		}
	} else {
		log.Printf("[MQTT] Publish client is not connect: topic: %s, mqttMsg %s", topic, msg)
	}
}
