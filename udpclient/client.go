package udpclient

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/lmj/mqtt-clients-demo/common"
	"github.com/lmj/mqtt-clients-demo/config"
	"github.com/lmj/mqtt-clients-demo/logger"
)

var (
	TnfGroup         []TerminalInfo
	ReStart          int       = -1
	ZclRegularReport int       = -1 //定时上报的帧序列号
)

// 该地方留给触发按键
func StartUDP(ctx context.Context, sig chan os.Signal, timestamp string) {
	TnfGroup = GenMode(config.DEVICE_TOTAL_COUNT)
	for _, v := range TnfGroup {
		go sendMsg(ctx, v, v.Client)
		go v.Client.readFromSocket(ctx, config.UDP_BUFFER_SIZE)
		go v.Client.processPackets(ctx, &v)
		v.Client.msgType <- common.TransFrame.IncrementAndStringGet(v.IP) + config.UDP_KEEPALIVE_MSG
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
				TnfGroup[ID-1].Client.msgType <- common.TransFrame.IncrementAndStringGet(TnfGroup[ID-1].IP) + config.UDP_TERMINAL_NETACCESS
			} else if curStatus == "终端离网" {
				ID, err := strconv.Atoi(TerminalID)
				if err != nil {
					return
				}
				TnfGroup[ID-1].Client.msgType <- common.TransFrame.IncrementAndStringGet(TnfGroup[ID-1].IP) + config.UDP_TERMINAL_NETLEAVE
			} else if curStatus == "设备上线" {
				ID, err := strconv.Atoi(TerminalID)
				if err != nil {
					return
				}
				go TriTimeReport(ctx, 600, TnfGroup[ID-1], false)
			}
		})

	}
	go func() {
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Log.Errorf("ListenAndServe: %s", err)
			}
		}()  //侦听等待处理
		select {  
		case <-ctx.Done():
			if err := srv.Close(); err != nil {
				logger.Log.Errorf("Server Close failed: %s", err)
			}
			log.Println("Server Exited Properly")
		}
	}()
}
