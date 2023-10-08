package udpclient

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/lmj/mqtt-clients-demo/config"
	"github.com/lmj/mqtt-clients-demo/logger"
)

var TnfGroup []TerminalInfo
var ServerCh chan bool

// 该地方留给触发按键
func StartUDP(sig chan os.Signal, timestamp string) {
	ServerCh = make(chan bool)
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
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		TerminalID, curStatus := r.FormValue("id"), r.FormValue("isLeave")
		if curStatus == "模拟终端入网" {
			ID, err := strconv.Atoi(TerminalID)
			if err != nil {
				return
			}
			TnfGroup[ID-1].Client.msgType <- config.UDP_TERMINAL_NETACCESS
		} else {
			ID, err := strconv.Atoi(TerminalID)
			if err != nil {
				return
			}
			TnfGroup[ID-1].Client.msgType <- config.UDP_TERMINAL_NETLEAVE
		}
	})
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Errorf("ListenAndServe: %s", err)
		}
		<-ServerCh
		if err := srv.Shutdown(nil); err != nil {
			logger.Log.Errorf("Server shutdown failed: %s", err)
		}
		log.Println("Server Exited Properly")
	}()
}
