package udpclient

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/lmj/mqtt-clients-demo/config"
	"github.com/lmj/mqtt-clients-demo/logger"
)

var (
	TnfGroup         []TerminalInfo
	ServerCh         chan bool //udp服务关闭
	RegularCh        chan bool //定时上报
	ReStart          int       = -1
	ZclRegularReport int       = -1 //定时上报的帧序列号
)

// 该地方留给触发按键
func StartUDP(sig chan os.Signal, timestamp string) {
	ServerCh = make(chan bool, 2)
	RegularCh = make(chan bool, 2)
	TnfGroup = GenMode(config.DEVICE_TOTAL_COUNT)
	for _, v := range TnfGroup {
		go sendMsg(v, v.Client)
		go v.Client.readFromSocket(config.UDP_BUFFER_SIZE)
		go v.Client.processPackets()
		v.Client.msgType <- config.UDP_KEEPALIVE_MSG
	}
	srv := &http.Server{
		Addr: ":7777",
	}
	if ReStart == 0 {
		http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
			TerminalID, curStatus := r.FormValue("id"), r.FormValue("isLeave")
			if curStatus == "终端入网" {
				ID, err := strconv.Atoi(TerminalID)
				if err != nil {
					return
				}
				TnfGroup[ID-1].Client.msgType <- config.UDP_TERMINAL_NETACCESS
			} else if curStatus == "终端离网" {
				ID, err := strconv.Atoi(TerminalID)
				if err != nil {
					return
				}
				TnfGroup[ID-1].Client.msgType <- config.UDP_TERMINAL_NETLEAVE
			} else if curStatus == "设备上线"{
				ID, err := strconv.Atoi(TerminalID)
				if err != nil {
					return
				}
				go TriTimeReport(600, TnfGroup[ID-1], false)
			}
		})

	}
	go func() {
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Log.Errorf("ListenAndServe: %s", err)
			}
		}()
		<-ServerCh
		if err := srv.Close(); err != nil {
			logger.Log.Errorf("Server Close failed: %s", err)
		}
		log.Println("Server Exited Properly")
	}()
}
