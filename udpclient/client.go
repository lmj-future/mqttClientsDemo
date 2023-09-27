package udpclient

import (
	"net/http"
	"os"
	"time"

	"github.com/coocood/freecache"
	"github.com/lmj/mqtt-clients-demo/config"
	"github.com/lmj/mqtt-clients-demo/logger"
)

var TnfGroup []TerminalInfo

// 该地方留给触发按键
func StartUDP(sig chan os.Signal, timestamp string) {
	TnfGroup = GenMode(config.UDP_T320M_NUM)
	for _, v := range TnfGroup {
		go sendMsg(v, v.Client)
		go v.Client.readFromSocket(config.UDP_BUFFER_SIZE)
		go v.Client.processPackets()
		v.Client.msgType <- "00002007"
	}

	//间隔检测内存使用情况
	go func() {
		var tickerID = time.NewTicker(time.Duration(config.UDP_ALIVE_CHECK_TIME) * 5 * time.Second)
		defer tickerID.Stop()
		for {
			<-tickerID.C
			var cacheNum int64 = 0
			ClientForEveryMsg.Range(func(key, value interface{}) bool {
				cacheNum += value.(*freecache.Cache).EntryCount()
				return true
			})
			logger.Log.Errorln("ClientForEveryMsg's EntryCount: ", cacheNum)
			// 监听程序退出
			select {
			case <-sig:
				return
			default:
				continue
			}
		}
	}()

	go func() {
		for {
			http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
				curStatus := r.FormValue("isLeave")
				if curStatus == "模拟终端入网" {
					for _, v := range TnfGroup {
						go func(t1 TerminalInfo) { //模拟终端离网
							t1.Client.msgType <- "00002001"
						}(v)
					}
				} else {
					for _, v := range TnfGroup {
						go func(t1 TerminalInfo) { //模拟终端入网
							t1.Client.msgType <- "00002002"
						}(v)
					}
				}
			})
			http.ListenAndServe(":7777", nil)
			select {
			case <-sig:
				return
			default:
				continue
			}
		}
	}()

}
