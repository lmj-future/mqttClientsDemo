package lorcaui

import (
	"os"
	"strconv"

	"github.com/lmj/mqtt-clients-demo/config"
	"github.com/lmj/mqtt-clients-demo/mqttclient"
	"github.com/lmj/mqtt-clients-demo/telnet"
	"github.com/lxn/win"
	"github.com/zserge/lorca"
)

var LORCA_UI lorca.UI

func Eval(dataLabal string, innerText string) {
	LORCA_UI.Eval(`
	var e = document.querySelector('div[data-label=` + dataLabal + `]');
	if (e !== null && e !== undefined) {var e2 = e.getElementsByClassName("text");e2[0].innerText = ` + innerText + `;}
	`)
}
//在这里触发按钮的停止和启动
func LorcaUI(sig chan os.Signal, do func(sig chan os.Signal, timestamp string), stop func(sig chan os.Signal), exit func(), args ...string) {
	hdc := win.GetDC(win.HWND(0))
	defer win.ReleaseDC(win.HWND(0), hdc)
	width := win.GetDeviceCaps(hdc, win.HORZRES)
	height := win.GetDeviceCaps(hdc, win.VERTRES)
	x := 1216
	y := 700
	args = append(args, "--window-position="+strconv.Itoa(int(width)/2-x/2)+","+strconv.Itoa(int(height)/2-y/2))
	LORCA_UI, _ = lorca.New("", "", x, y, args...)
	defer LORCA_UI.Close()

	// 以Data URI协议读取
	// ui.Load("data:text/html," + url.PathEscape(newFile))

	// 以文件方式读取
	dir, _ := os.Getwd()
	LORCA_UI.Load("file:///" + dir + "/html/index.html")

	LORCA_UI.Bind("start", func(jsonString string, timestamp string) {
		Eval(`"CLIENT_STATE"`, `"开始连接..."`)
		mqttclient.ClientStopMap.Set(timestamp, false)
		config.SetConfig(jsonString)
		do(sig, timestamp)
	})
	LORCA_UI.Bind("stop", func(timestamp string) {
		Eval(`"CLIENT_STATE"`, `"停止连接，开始销毁资源..."`)
		Eval(`"DEVICE_ONLINE_COUNT"`, `"0"`)
		Eval(`"DEVICE_OFFLINE_COUNT"`, `"0"`)
		Eval(`"TOTAL_PPS"`, `"0"`)
		mqttclient.ClientStopMap.Set(timestamp, true)
		stop(sig)
		Eval(`"CLIENT_STATE"`, `"等待开始"`)
	})
	LORCA_UI.Bind("pubLog", func(enable bool) {
		config.LOG_MQTT_PUBLISH_ENABLE = enable
	})
	LORCA_UI.Bind("subLog", func(enable bool) {
		config.LOG_MQTT_SUBSCRIBE_ENABLE = enable
	})
	LORCA_UI.Bind("telnet", func(enable bool) {
		if enable {
			telnet.Start()
		} else {
			telnet.Stop()
		}
	})

	// Wait for the window to be closed
	select {
	case <-sig:
		LORCA_UI.Close()
		exit()
	case <-LORCA_UI.Done():
		exit()
	}
}
