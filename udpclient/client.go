package udpclient

import (
	"demo/config"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coocood/freecache"
	"github.com/jessevdk/go-flags"
	cmap "github.com/orcaman/concurrent-map"
)

var ClientMapUDP cmap.ConcurrentMap

var Opts struct {
	ServerAddress string `short:"s" long:"serveraddress" default:"33.33.33.244:3501" description:"The Server's Address"`
	Buffer        int    `short:"b" long:"buffer" default:"1024" description:"max buffer size for the socket io"`
	Quiet         bool   `short:"q" long:"quiet" description:"whether to print logging info or not"`
	UdpVer        string `short:"u" long:"udpver" default:"0102" description:"The version number of the udp message"`
}

var TnfGroup []TerminalInfo

func init() {
	ClientMapUDP = cmap.New()
	_, err := flags.Parse(&Opts)
	errorCheck(err, "init", true)
	if Opts.Quiet {
		log.SetLevel(log.WarnLevel)
	}
	formatter := &log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	}
	log.SetFormatter(formatter)
}

func StartUDP() {
	termiGroup := GenMode(config.UDP_T320M_NUM)
	for _, v := range termiGroup {
		go sendMsg(v, v.client)
		go v.client.readFromSocket(Opts.Buffer)
		go v.client.processPackets()
		v.client.msgType <- "0000" //触发
	}

	//间隔检测内存使用情况
	go func() {
		var tickerID = time.NewTicker(time.Duration(10) * time.Second)
		defer tickerID.Stop()
		for {
			<-tickerID.C
			var cacheNum int64 = 0
			MsgForEvery.Range(func(key, value interface{}) bool {
				cacheNum += value.(*freecache.Cache).EntryCount()
				return true
			})
			log.Println("MsgForEvery' EntryCount: ", cacheNum)
		}
	}()
}
