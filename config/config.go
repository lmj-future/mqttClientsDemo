

package config

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-ini/ini"
)

var MQTT_SERVER_HOST string
var MQTT_SERVER_PORT string
var MQTT_CLIENT_USERNAME string
var MQTT_CLIENT_PASSWORD string
var MQTT_CLIENT_BROKER string
var MQTT_CLIENT_KEEPALIVE int
var MQTT_CLIENT_CONNECT_INTERVAL int
var MQTT_CLIENT_CONNECT_PER_100_INTERVAL int
var MQTT_CLIENT_RECONNECT_INTERVAL int
var MQTT_CLIENT_RECONNECT_COUNT int
var MQTT_CLIENT_SLEEP_INTERVAL int
var PRODUCT_NAME string
var PRODUCT_KEY string
var DEVICE_SN_PRE string
var DEVICE_SN_MID string
var DEVICE_SN_LEN int
var DEVICE_SN_LEFT_LEN string
var DEVICE_SN_SUF_START_BY int
var DEVICE_TOTAL_COUNT int

var PPS int
var PPS_PER int

var SUB_TOPIC []string
var PUB_TOPIC []string

var MQTT_SUBSCRIBE_ENABLE bool
var MQTT_PUBLISH_ENABLE bool

//udp部分需用
var (
	UDP_T320M_NUM int
	UDP_ALIVE_CHECK_TIME int
)

//根据配置文件初始化，初始化后的文件指针
func Init() *ini.File {
	cfg, err := ini.Load("./config/conf.ini")
	if err != nil {
		log.Printf("ini.Load err = %v", err)
		panic(err)
	}
	PRODUCT_NAME = cfg.Section("conf").Key("PRODUCT_NAME").String()
	MQTT_SERVER_HOST = cfg.Section("conf").Key("MQTT_SERVER_HOST").String()
	MQTT_SERVER_PORT = cfg.Section("conf").Key("MQTT_SERVER_PORT").String()
	MQTT_CLIENT_BROKER = fmt.Sprintf("tcp://%s:%s", MQTT_SERVER_HOST, MQTT_SERVER_PORT)
	MQTT_CLIENT_CONNECT_INTERVAL, _ = cfg.Section("conf").Key("MQTT_CLIENT_CONNECT_INTERVAL").Int()
	if MQTT_CLIENT_CONNECT_INTERVAL < 20 {
		MQTT_CLIENT_CONNECT_INTERVAL = 20
	}
	MQTT_CLIENT_CONNECT_PER_100_INTERVAL, _ = cfg.Section("conf").Key("MQTT_CLIENT_CONNECT_PER_100_INTERVAL").Int()
	DEVICE_SN_MID = cfg.Section("conf").Key("DEVICE_SN_MID").String()
	switch PRODUCT_NAME {
	case "T320M", "T320MX", "T320MX-U":
		MQTT_CLIENT_KEEPALIVE = 60
		MQTT_CLIENT_RECONNECT_INTERVAL = 30
		MQTT_CLIENT_RECONNECT_COUNT = 20
		MQTT_CLIENT_SLEEP_INTERVAL = 60
		DEVICE_SN_LEN = 20
		SUB_TOPIC = []string{"/sys/%s/%s/thing/model/down_raw", "/ota/device/upgrade/%s/%s"}
		PUB_TOPIC = []string{"/sys/%s/%s/thing/model/up_raw", "/ota/device/progress/%s/%s"}
		switch PRODUCT_NAME {
		case "T320M":
			MQTT_CLIENT_USERNAME = "EbcMXsMg+c"
			MQTT_CLIENT_PASSWORD = "GnW7YLdtWv"
			PRODUCT_KEY = "kiSHgWsG"
			DEVICE_SN_PRE = "219801A26U"
		case "T320MX":
			MQTT_CLIENT_USERNAME = "kT/7i7p+h6"
			MQTT_CLIENT_PASSWORD = "n7eraCLyKa"
			PRODUCT_KEY = "COc0mEdF"
			DEVICE_SN_PRE = "219801A26N"
		case "T320MX-U":
			MQTT_CLIENT_USERNAME = "Hlhvs9xQWX"
			MQTT_CLIENT_PASSWORD = "EO8e96z0Pc"
			PRODUCT_KEY = "bzwtFAVT"
			DEVICE_SN_PRE = "219801A2YH"
		}
	default:
		MQTT_CLIENT_KEEPALIVE, _ = cfg.Section("conf").Key("MQTT_CLIENT_KEEPALIVE").Int()
		MQTT_CLIENT_RECONNECT_INTERVAL, _ = cfg.Section("conf").Key("MQTT_CLIENT_RECONNECT_INTERVAL").Int()
		MQTT_CLIENT_RECONNECT_COUNT, _ = cfg.Section("conf").Key("MQTT_CLIENT_RECONNECT_COUNT").Int()
		MQTT_CLIENT_SLEEP_INTERVAL, _ = cfg.Section("conf").Key("MQTT_CLIENT_SLEEP_INTERVAL").Int()
		DEVICE_SN_LEN, _ = cfg.Section("conf").Key("DEVICE_SN_LEN").Int()
		SUB_TOPIC = cfg.Section("topic").Key("MQTT_SUBSCRIBE_TOPIC").Strings(";")
		PUB_TOPIC = cfg.Section("topic").Key("MQTT_PUBLISH_TOPIC").Strings(";")
		MQTT_CLIENT_USERNAME = cfg.Section("conf").Key("MQTT_CLIENT_USERNAME").String()
		MQTT_CLIENT_PASSWORD = cfg.Section("conf").Key("MQTT_CLIENT_PASSWORD").String()
		PRODUCT_KEY = cfg.Section("conf").Key("PRODUCT_KEY").String()
		DEVICE_SN_PRE = cfg.Section("conf").Key("DEVICE_SN_PRE").String()
	}
	DEVICE_SN_LEFT_LEN = strconv.FormatInt(int64(DEVICE_SN_LEN-len(DEVICE_SN_PRE)-len(DEVICE_SN_MID)), 10)
	DEVICE_SN_SUF_START_BY, _ = cfg.Section("conf").Key("DEVICE_SN_SUF_START_BY").Int()
	DEVICE_TOTAL_COUNT, _ = cfg.Section("conf").Key("DEVICE_TOTAL_COUNT").Int()

	PPS, _ = cfg.Section("pps").Key("PPS").Int()
	PPS_PER, _ = cfg.Section("pps").Key("PPS_PER").Int()

	MQTT_SUBSCRIBE_ENABLE, _ = cfg.Section("log").Key("MQTT_SUBSCRIBE_ENABLE").Bool()
	MQTT_PUBLISH_ENABLE, _ = cfg.Section("log").Key("MQTT_PUBLISH_ENABLE").Bool()

	UDP_T320M_NUM, _ = cfg.Section("udp").Key("UDP_T320M_NUM").Int()
	UDP_ALIVE_CHECK_TIME, _ = cfg.Section("udp").Key("UDP_ALIVE_CHECK_TIME").Int()
	return cfg
}
