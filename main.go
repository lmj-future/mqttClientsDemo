package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/pborman/uuid"

	"github.com/lmj/mqtt-clients-demo/common"
	"github.com/lmj/mqtt-clients-demo/config"
	"github.com/lmj/mqtt-clients-demo/logger"
	"github.com/lmj/mqtt-clients-demo/lorcaui"
	"github.com/lmj/mqtt-clients-demo/mqttclient"
)

var clientState string = ""
var onlineDev int32 = 0

func init() {
	config.Init()
	logger.Init()
	mqttclient.Init()
	common.InitCounter()
	common.InitTimeStampMap()
}

func checkConnectOK(index int32) {
	if clientState == "" {
		go func() {
			<-time.After(time.Second)
			checkConnectOK(index)
		}()
	} else if clientState == "ok" {
		leftCount := config.DEVICE_TOTAL_COUNT - int(index)
		lorcaui.Eval(`"CLIENT_STATE"`, `"`+"正在进行客户端连接，大概还需等待"+fmt.Sprintf("%v", leftCount)+"个设备上线"+`"`)
	} else {
		lorcaui.Eval(`"CLIENT_STATE"`, `"`+"正在进行客户端连接，连接失败了["+clientState+"]，请检查一下配置是否正确"+`"`)
	}
}

func getSubTopic(devSN string) []string {
	var subTopic []string
	for _, v := range config.SUB_TOPIC {
		if v != "" {
			if strings.Contains(v, "{devSN}") {
				v = strings.ReplaceAll(v, "{devSN}", devSN)
			}
			if strings.Contains(v, "{devKey}") {
				v = strings.ReplaceAll(v, "{devKey}", config.DEVICE_KEY)
			}
			subTopic = append(subTopic, v)
		}
	}
	return subTopic
}

func getUserName(userName string, devSN string) string {
	if userName != "" {
		if strings.Contains(userName, "{devSN}") {
			userName = strings.ReplaceAll(userName, "{devSN}", devSN)
		}
		if strings.Contains(userName, "{devKey}") {
			userName = strings.ReplaceAll(userName, "{devKey}", config.DEVICE_KEY)
		}
	} else {
		userName = uuid.NewRandom().String()
	}
	return userName
}

func getPassword(password string, devSN string) string {
	if password != "" {
		if strings.Contains(password, "{devSN}") {
			password = strings.ReplaceAll(password, "{devSN}", devSN)
		}
		if strings.Contains(password, "{devKey}") {
			password = strings.ReplaceAll(password, "{devKey}", config.DEVICE_KEY)
		}
	} else {
		password = uuid.NewRandom().String()
	}
	return password
}

func getClientId(userName string, password string, devSN string) string {
	var clientId string
	if config.MQTT_CLIENT_ID != "" {
		clientId = config.MQTT_CLIENT_ID
		if strings.Contains(clientId, "{devSN}") {
			clientId = strings.ReplaceAll(clientId, "{devSN}", devSN)
		}
		if strings.Contains(clientId, "{userName}") {
			clientId = strings.ReplaceAll(clientId, "{userName}", userName)
		}
		if strings.Contains(clientId, "{password}") {
			clientId = strings.ReplaceAll(clientId, "{password}", password)
		}
		if strings.Contains(clientId, "{devKey}") {
			clientId = strings.ReplaceAll(clientId, "{devKey}", config.DEVICE_KEY)
		}
	} else {
		clientId = uuid.NewRandom().String()
	}
	return clientId
}


func mqttConnectFirst(sig chan os.Signal, timestamp string, sn int) bool {
	clientState = ""
	isNeedWait := true
	sufFix := "%0" + config.DEVICE_SN_LEFT_LEN + "d"
	wg := &sync.WaitGroup{}
	wg.Add(1)
	devSN := fmt.Sprintf(config.DEVICE_SN_PRE+config.DEVICE_SN_MID+sufFix, config.DEVICE_SN_SUF_START_BY+sn)
	var clientId string
	var userName string = config.MQTT_CLIENT_USERNAME
	var password string = config.MQTT_CLIENT_PASSWORD
	var subTopic []string
	if config.PRODUCT_NAME == "自定义" {
		subTopic = getSubTopic(devSN)
	} else {
		for _, topic := range config.SUB_TOPIC {
			subTopic = append(subTopic, fmt.Sprintf(topic, config.PRODUCT_KEY, devSN))
		}
	}
	switch config.PRODUCT_NAME {
	case "T320M", "T320MX", "T320MX-U":
		clientId = fmt.Sprintf("%s&%s", config.PRODUCT_KEY, devSN)
	case "示例产品-mqtt":
		clientId = fmt.Sprintf("%s&v5", devSN)
		userName = fmt.Sprintf("%s&%s", config.PRODUCT_KEY, devSN)
		password = config.DEVICE_KEY
	case "自定义":
		userName = getUserName(userName, devSN)
		password = getPassword(password, devSN)
		clientId = getClientId(userName, password, devSN)
	}
	go mqttclient.Connect(devSN, config.MQTT_CLIENT_BROKER, userName, password, clientId,
		time.Duration(config.MQTT_CLIENT_KEEPALIVE)*time.Second, subTopic, wg, timestamp, &clientState, &onlineDev)
	if config.DEVICE_TOTAL_COUNT > 100 {
		// 特定序号进行数据打印
		if onlineDev%100 == 0 {
			checkConnectOK(onlineDev)
		}
	} else {
		if onlineDev == 0 {
			checkConnectOK(onlineDev)
		}
	}
	// 如果此次开始被停止了，那么就要停止此次开始所做的事情
	if v, ok := mqttclient.ClientStopMap.Get(timestamp); ok && v.(bool) {
		stop(sig)
		isNeedWait = false
	}
	// 监听程序退出
	select {
	case <-sig:
		exit()
	default:
	}
	if isNeedWait {
		wg.Wait()
		return true
	}
	return false
}

//一个设备的订阅
func mqttSubscribeSecond(sig chan os.Signal, timestamp string, sn int) bool {
	isNeedWait := true
	sufFix := "%0" + config.DEVICE_SN_LEFT_LEN + "d"
	devSN := fmt.Sprintf(config.DEVICE_SN_PRE+config.DEVICE_SN_MID+sufFix, config.DEVICE_SN_SUF_START_BY+sn)
	var subTopic []string
	if config.PRODUCT_NAME == "自定义" {
		subTopic = getSubTopic(devSN)
	} else {
		for _, topic := range config.SUB_TOPIC {
			subTopic = append(subTopic, fmt.Sprintf(topic, config.PRODUCT_KEY, devSN))
		}
	}
	if v, ok := mqttclient.ClientMap.Get(devSN); ok {
		cli := v.(common.MqttClientInfo)
		mqttclient.Subscribe(subTopic, cli.Client)
		logger.Log.Infof("设备[%s]订阅[%+v]成功", cli.DevSN, subTopic)
	}
	// 如果此次开始被停止了，那么就要停止此次开始所做的事情
	if v, ok := mqttclient.ClientStopMap.Get(timestamp); ok && v.(bool) {
		stop(sig)
		isNeedWait = false
	}
	// 监听程序退出
	select {
	case <-sig:
		exit()
	default:
	}
	return isNeedWait
}

//一个设备的数据交互
func clientBusinessThird(sig chan os.Signal, timestamp string, sn int, wg *sync.WaitGroup) bool {
	isNeedWait := true
	sufFix := "%0" + config.DEVICE_SN_LEFT_LEN + "d"
	devSN := fmt.Sprintf(config.DEVICE_SN_PRE+config.DEVICE_SN_MID+sufFix, config.DEVICE_SN_SUF_START_BY+sn)
	if v, ok := mqttclient.ClientMap.Get(devSN); ok {
		cli := v.(common.MqttClientInfo)
		go common.UpRawWhenConnect(cli)
		logger.Log.Infof("设备[%s]正在进行交互", cli.DevSN)
	}
	// 如果此次开始被停止了，那么就要停止此次开始所做的事情
	if v, ok := mqttclient.ClientStopMap.Get(timestamp); ok && v.(bool) {
		stop(sig)
		isNeedWait = false
	}
	// 监听程序退出
	select {
	case <-sig:
		exit()
	default:
	}
	wg.Done()
	return isNeedWait
}

//数据定时上报异步
func clientTickerData(sig chan os.Signal, timestamp string) {
	ticker := time.NewTicker(time.Duration(config.PPS_PER) * time.Second)
	go func() {
		for {
			// 如果此次开始被停止了，那么就要停止此次开始所做的事情
			if v, ok := mqttclient.ClientStopMap.Get(timestamp); ok && v.(bool) {
				ticker.Stop()
				stop(sig)
				break
			}
			select {
			case <-ticker.C:
				go func() {
					sufFix := "%0" + config.DEVICE_SN_LEFT_LEN + "d"
					for i := 0; i < config.DEVICE_TOTAL_COUNT; i++ {
						go func(index int) {
							devSN := fmt.Sprintf(config.DEVICE_SN_PRE+config.DEVICE_SN_MID+sufFix, config.DEVICE_SN_SUF_START_BY+index)
							if v, ok := mqttclient.ClientMap.Get(devSN); ok {
								cli := v.(common.MqttClientInfo)
								timeSleep := time.Second * time.Duration(config.PPS_PER) / time.Duration(config.DEVICE_TOTAL_COUNT)
								go common.UpRawAfterConnect(cli, timeSleep)
								<-time.After(timeSleep)
							}
						}(i)
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

func checkAndDisplay(timestamp string) {
	//<-time.After(3 * time.Second)
	for {
		// 如果此次开始被停止了，那么就要停止此次开始所做的事情
		if v, ok := mqttclient.ClientStopMap.Get(timestamp); ok && v.(bool) {
			break
		}

		unconnCli := []string{}
		for i := range mqttclient.ClientMap.IterBuffered() {
			cli := i.Val.(common.MqttClientInfo)
			if !cli.Client.IsConnectionOpen() {
				unconnCli = append(unconnCli, cli.DevSN)
			}
		}
		if len(unconnCli) > 0 {
			logger.Log.Warnf("当前离线的设备:%v", unconnCli)
		}
		totalCount := mqttclient.ClientMap.Count()
		onlineCount := mqttclient.ClientMap.Count() - len(unconnCli)
		offlineCount := len(unconnCli)
		pps := float64(config.PPS) / float64(config.PPS_PER)
		totalPps := pps * float64(mqttclient.ClientMap.Count())
		logger.Log.Infof("设备总数[%d] 设备在线数[%d] 设备离线数[%d] 单个设备pps[%.2f] 总pps[%.2f]",
			totalCount, onlineCount, offlineCount, pps, totalPps)
		lorcaui.Eval(`"DEVICE_ONLINE_COUNT"`, `"`+fmt.Sprintf("%d", onlineCount)+`"`)
		lorcaui.Eval(`"DEVICE_OFFLINE_COUNT"`, `"`+fmt.Sprintf("%d", offlineCount)+`"`)
		lorcaui.Eval(`"TOTAL_PPS"`, `"`+fmt.Sprintf("%.2f", totalPps)+`"`)
		<-time.After(3 * time.Second)
	}
}

func stop(sig chan os.Signal) {
	logger.Log.Infoln("客户端触发停止，正在销毁资源！！！请等待...")
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
		var subTopic []string
		if config.PRODUCT_NAME == "自定义" {
			subTopic = getSubTopic(cli.DevSN)
		} else {
			for _, topic := range config.SUB_TOPIC {
				subTopic = append(subTopic, fmt.Sprintf(topic, config.PRODUCT_KEY, cli.DevSN))
			}
		}
		cli.Client.Unsubscribe(subTopic...)
		cli.Client.Disconnect(uint(config.MQTT_CLIENT_CONNECT_INTERVAL))
		logger.Log.Infoln("客户端停止，关闭MQTT连接: ", cli.DevSN)
		<-time.After(time.Duration(config.MQTT_CLIENT_CONNECT_INTERVAL) * time.Microsecond)
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
	logger.Log.Infoln("所有资源已销毁")
}

func exit() {
	logger.Log.Infoln("进程将在销毁所有资源后退出！！！请等待...")
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
		var subTopic []string
		if config.PRODUCT_NAME == "自定义" {
			subTopic = getSubTopic(cli.DevSN)
		} else {
			for _, topic := range config.SUB_TOPIC {
				subTopic = append(subTopic, fmt.Sprintf(topic, config.PRODUCT_KEY, cli.DevSN))
			}
		}
		cli.Client.Unsubscribe(subTopic...)
		cli.Client.Disconnect(uint(config.MQTT_CLIENT_CONNECT_INTERVAL))
		logger.Log.Infoln("进程退出，关闭MQTT连接: ", cli.DevSN)
		<-time.After(time.Duration(config.MQTT_CLIENT_CONNECT_INTERVAL) * time.Microsecond)
		wg.Done()
	}
	wg.Wait()
	logger.Log.Warnln("========================Exit=======================")
	fmt.Println(`
===================================================
:: 进程已优雅地退出，点个赞吧
===================================================
	`)
	os.Exit(0)
}

func do(sig chan os.Signal, timestamp string) {
	wgForBusiness :=  &sync.WaitGroup{} 
	wgForBusiness.Add(config.DEVICE_TOTAL_COUNT)
	for i := 0; i < config.DEVICE_TOTAL_COUNT; i++ {
		go func(index int) {
			//1、先进行MQTT连接
			isNeedNext := mqttConnectFirst(sig, timestamp, index)
			if isNeedNext {
				//其次完成订阅
				if _, ok := mqttclient.ClientStopMap.Get(timestamp); ok {
					isNeedNext = mqttSubscribeSecond(sig, timestamp, index)
				}
			}
			if isNeedNext { 
				//完成数据交互
				if _, ok := mqttclient.ClientStopMap.Get(timestamp); ok {
					isNeedNext = clientBusinessThird(sig, timestamp, index, wgForBusiness)
				}
			}
			if isNeedNext {
				if _, ok := mqttclient.ClientStopMap.Get(timestamp); ok {
					clientTickerData(sig, timestamp)
				}
			}
		}(i)
	}
	// 展示一些数据
	if _, ok := mqttclient.ClientStopMap.Get(timestamp); ok {
		go checkAndDisplay(timestamp)
	}
	wgForBusiness.Wait()  //等待数据交互后显示客户端情况
	lorcaui.Eval(`"CLIENT_STATE"`, `"客户端运行中..."`)
}

// 主程序
func main() {
	fmt.Println(`
===================================================
:: MQTT CLIENTS DEMO ::         (v1.0.2)
===================================================
	`)
	logger.Log.Warnln("========================Start======================")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	signal.Ignore(syscall.SIGPIPE)

	if runtime.GOOS == "windows" {
		// windows环境下支持可视化界面操作
		lorcaui.LorcaUI(sig, do, stop, exit, "--remote-allow-origins=*")
	} else {
		// 其他操作系统暂不支持可视化界面，需要修改配置文件
		mqttclient.ClientStopMap.Set("notSupport", false)
		do(sig, "notSupport")
	}

	<-sig
	exit()
}
