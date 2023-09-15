
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	cmap "github.com/orcaman/concurrent-map"

	"demo/common"
	"demo/config"
	"demo/mqttclient"
	"demo/udpclient"
)

//初始化
func init() {
	config.Init()
	mqttclient.Init()    //初始化的时候开一个map锁
	common.InitCounter() //以下都是
	common.InitTimeStampMap()
}

func mqttConnectFirst(sig chan os.Signal) {
	sufFix := "%0" + config.DEVICE_SN_LEFT_LEN + "d"
	wg := &sync.WaitGroup{}
	wg.Add(config.DEVICE_TOTAL_COUNT)
	for i := 0; i < config.DEVICE_TOTAL_COUNT; i++ {
		devSN := fmt.Sprintf(config.DEVICE_SN_PRE+config.DEVICE_SN_MID+sufFix, config.DEVICE_SN_SUF_START_BY+i)
		clientId := fmt.Sprintf("%s&%s", config.PRODUCT_KEY, devSN)
		subTopic := []string{
			fmt.Sprintf(config.SUB_TOPIC[0], config.PRODUCT_KEY, devSN),
			fmt.Sprintf(config.SUB_TOPIC[1], config.PRODUCT_KEY, devSN),
		}
		go mqttclient.Connect(config.MQTT_CLIENT_BROKER, config.MQTT_CLIENT_USERNAME, config.MQTT_CLIENT_PASSWORD, clientId,
			time.Duration(config.MQTT_CLIENT_KEEPALIVE)*time.Second, subTopic, wg)
		// 每100个设备连接完后，停顿MQTT_CLIENT_CONNECT_PER_100_INTERVAL秒，再继续下一组100个设备的连接，防止瞬时连接量太大
		if i%100 == 0 {
			<-time.After(time.Second * time.Duration(config.MQTT_CLIENT_CONNECT_PER_100_INTERVAL))
		}
		<-time.After(time.Millisecond * time.Duration(config.MQTT_CLIENT_CONNECT_INTERVAL))

		// 监听程序退出
		select {
		case <-sig:
			exit()
		default:
			continue
		}
	}
	wg.Wait()
}

func mqttSubscribeSecond(sig chan os.Signal) {
	sufFix := "%0" + config.DEVICE_SN_LEFT_LEN + "d"
	wg := &sync.WaitGroup{}
	wg.Add(config.DEVICE_TOTAL_COUNT)
	for i := 0; i < config.DEVICE_TOTAL_COUNT; i++ {
		devSN := fmt.Sprintf(config.DEVICE_SN_PRE+config.DEVICE_SN_MID+sufFix, config.DEVICE_SN_SUF_START_BY+i)
		clientId := fmt.Sprintf("%s&%s", config.PRODUCT_KEY, devSN)
		subTopic := []string{
			fmt.Sprintf(config.SUB_TOPIC[0], config.PRODUCT_KEY, devSN),
			fmt.Sprintf(config.SUB_TOPIC[1], config.PRODUCT_KEY, devSN),
		}
		if v, ok := mqttclient.ClientMap.Get(clientId); ok {
			cli := v.(common.MqttClientInfo)
			mqttclient.Subscribe(subTopic, cli.Client)
			log.Printf("设备[%s]订阅[%+v]成功", cli.ClientId, subTopic)
			<-time.After(time.Millisecond * time.Duration(config.MQTT_CLIENT_CONNECT_INTERVAL))
		}
		wg.Done()

		// 监听程序退出
		select {
		case <-sig:
			exit()
		default:
			continue
		}
	}
	wg.Wait()
}

func clientBusinessThird(sig chan os.Signal) {
	sufFix := "%0" + config.DEVICE_SN_LEFT_LEN + "d"
	wg := &sync.WaitGroup{}
	wg.Add(config.DEVICE_TOTAL_COUNT)
	for i := 0; i < config.DEVICE_TOTAL_COUNT; i++ {
		devSN := fmt.Sprintf(config.DEVICE_SN_PRE+config.DEVICE_SN_MID+sufFix, config.DEVICE_SN_SUF_START_BY+i)
		clientId := fmt.Sprintf("%s&%s", config.PRODUCT_KEY, devSN)
		if v, ok := mqttclient.ClientMap.Get(clientId); ok {
			cli := v.(common.MqttClientInfo)
			go common.UpRawWhenConnect(cli)
			log.Printf("设备[%s]正在进行交互", cli.ClientId)
			<-time.After(time.Millisecond * time.Duration(config.MQTT_CLIENT_CONNECT_INTERVAL))
		}
		wg.Done()
		// 监听程序退出
		select {
		case <-sig:
			exit()
		default:
			continue
		}
	}
	wg.Wait()
}

func clientTickerData(sig chan os.Signal) {
	ticker := time.NewTicker(time.Duration(config.PPS_PER) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				go func() {
					sufFix := "%0" + config.DEVICE_SN_LEFT_LEN + "d"
					for i := 0; i < config.DEVICE_TOTAL_COUNT; i++ {
						devSN := fmt.Sprintf(config.DEVICE_SN_PRE+config.DEVICE_SN_MID+sufFix, config.DEVICE_SN_SUF_START_BY+i)
						clientId := fmt.Sprintf("%s&%s", config.PRODUCT_KEY, devSN)
						if v, ok := mqttclient.ClientMap.Get(clientId); ok {
							cli := v.(common.MqttClientInfo)
							timeSleep := time.Second * time.Duration(config.PPS_PER) / time.Duration(config.DEVICE_TOTAL_COUNT)
							go common.UpRawAfterConnect(cli, timeSleep)
							<-time.After(timeSleep)
						}
					}
				}()
			// 监听程序退出
			case <-sig:
				ticker.Stop()
				exit()
			default:
				continue
			}
		}
	}()
}

func checkAndDisplay() {
	<-time.After(3 * time.Second)
	for {
		unconnCli := []string{}
		for i := range mqttclient.ClientMap.IterBuffered() {
			cli := i.Val.(common.MqttClientInfo)
			if !cli.Client.IsConnectionOpen() {
				unconnCli = append(unconnCli, cli.ClientId)
			}
		}
		if len(unconnCli) > 0 {
			log.Printf("当前离线的设备:%v", unconnCli)
		}
		log.Printf("设备总数[%d] 设备在线数[%d] 设备离线数[%d] 单个设备pps[%.2f] 总pps[%.2f]",
			mqttclient.ClientMap.Count(), mqttclient.ClientMap.Count()-len(unconnCli), len(unconnCli),
			float64(config.PPS)/float64(config.PPS_PER),
			float64(config.PPS)/float64(config.PPS_PER)*float64(mqttclient.ClientMap.Count()))

		<-time.After(10 * time.Second)
	}
}

//提取mqtt的对应引用，释放mqttclient资源
//取消订阅并关闭来连接
func exit() {
	log.Println("进程将在销毁所有资源后退出！！！请等待...")
	<-time.After(3 * time.Second)
	clientMap := cmap.New()
	for i := range mqttclient.ClientMap.IterBuffered() {
		clientMap.Set(i.Key, i.Val)
	}
	mqttclient.ClientMap.Clear()
	wg := &sync.WaitGroup{}
	wg.Add(clientMap.Count())
	for i := range clientMap.IterBuffered() {
		cli := i.Val.(common.MqttClientInfo)
		devSN := strings.Replace(cli.ClientId, config.PRODUCT_KEY+"&", "", 1)
		subTopic := []string{
			fmt.Sprintf(config.SUB_TOPIC[0], config.PRODUCT_KEY, devSN),
			fmt.Sprintf(config.SUB_TOPIC[1], config.PRODUCT_KEY, devSN),
		}
		cli.Client.Unsubscribe(subTopic...)
		cli.Client.Disconnect(uint(config.MQTT_CLIENT_CONNECT_INTERVAL))
		log.Println("进程退出, 关闭MQTT连接: ", cli.ClientId)
		<-time.After(time.Duration(config.MQTT_CLIENT_CONNECT_INTERVAL) * time.Microsecond)
		wg.Done()
	}
	//udp客户端的关闭
	<-time.After(1 * time.Second)
	clientMapUDP := cmap.New()
	for j := range udpclient.ClientMapUDP.IterBuffered() {
		clientMap.Set(j.Key,j.Val)
	}
	udpclient.ClientMapUDP.Clear()
	wg.Add(clientMapUDP.Count())
	for  k := range clientMapUDP.IterBuffered() {
		c := k.Val.(udpclient.Client)
		c.Connection.Close()
		log.Println("进程退出,关闭UDP客户端连接: ", k.Key)
		<-time.After(time.Duration(200) * time.Microsecond)
		wg.Done()
	}
	wg.Wait()
	log.Println("进程已优雅的退出，点个赞吧")
	os.Exit(0)
}

//主程序
func main() {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	signal.Ignore(syscall.SIGPIPE)
	udpclient.StartUDP()
	// 1、先进行MQTT连接
	// mqttConnectFirst(sig)

	// // 2、再进行MQTT消息订阅
	// mqttSubscribeSecond(sig)

	// // 3、然后进行数据交互
	// clientBusinessThird(sig)
	// log.Printf("所有设备交互已完成，你可以进行操作了")

	// // 4、启动定时器，对数据进行定时上送
	// clientTickerData(sig)

	// // 5、展示一些数据
	// go checkAndDisplay()

	

	<-sig
	exit()
}
