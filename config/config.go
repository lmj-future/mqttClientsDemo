package config

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
)

var MQTT_SERVER_HOST string
var MQTT_SERVER_PORT string
var MQTT_CLIENT_USERNAME string
var MQTT_CLIENT_PASSWORD string
var MQTT_CLIENT_ID string
var MQTT_CLIENT_BROKER string
var MQTT_CLIENT_KEEPALIVE int
var MQTT_CLIENT_CONNECT_INTERVAL int
var MQTT_CLIENT_CONNECT_PER_100_INTERVAL int
var MQTT_CLIENT_RECONNECT_INTERVAL int
var MQTT_CLIENT_RECONNECT_COUNT int
var MQTT_CLIENT_SLEEP_INTERVAL int
var PRODUCT_NAME string
var PRODUCT_KEY string
var DEVICE_KEY string
var DEVICE_SN_PRE string
var DEVICE_SN_MID string
var DEVICE_SN_LEN int
var DEVICE_SN_LEFT_LEN string
var DEVICE_SN_SUF_START_BY int
var DEVICE_TOTAL_COUNT int

var PROPERTY_UP_ENABLE bool
var BOOLKey bool
var DoubleKey bool
var LongKey bool
var StringKey bool
var PowerSwitch bool
var CustomKey bool
var CustomValue string

var PPS int
var PPS_PER int

var SUB_TOPIC []string
var PUB_TOPIC []string

var LOG_MQTT_SUBSCRIBE_ENABLE bool
var LOG_MQTT_PUBLISH_ENABLE bool
var LOG_PATH string
var LOG_LEVEL string
var LOG_AGE int
var LOG_SIZE int
var LOG_BACKUP_COUNT int
var LOG_STORE string

//控制消息
var (
	DevFinshedConn chan int
	DevFinshedSub  chan int
	DevIsNeedNext = true
)

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
	DEVICE_KEY = cfg.Section("conf").Key("DEVICE_KEY").String()
	DEVICE_SN_PRE = cfg.Section("conf").Key("DEVICE_SN_PRE").String()

	DEVICE_SN_SUF_START_BY, _ = cfg.Section("conf").Key("DEVICE_SN_SUF_START_BY").Int()
	DEVICE_TOTAL_COUNT, _ = cfg.Section("conf").Key("DEVICE_TOTAL_COUNT").Int()

	PROPERTY_UP_ENABLE, _ = cfg.Section("conf").Key("PROPERTY_UP_ENABLE").Bool()
	BOOLKey, _ = cfg.Section("conf").Key("BOOLKey").Bool()
	DoubleKey, _ = cfg.Section("conf").Key("DoubleKey").Bool()
	LongKey, _ = cfg.Section("conf").Key("LongKey").Bool()
	StringKey, _ = cfg.Section("conf").Key("StringKey").Bool()
	PowerSwitch, _ = cfg.Section("conf").Key("PowerSwitch").Bool()
	CustomKey, _ = cfg.Section("conf").Key("CustomKey").Bool()
	CustomValue = cfg.Section("conf").Key("CustomValue").String()

	PPS, _ = cfg.Section("pps").Key("PPS").Int()
	PPS_PER, _ = cfg.Section("pps").Key("PPS_PER").Int()

	LOG_MQTT_SUBSCRIBE_ENABLE, _ = cfg.Section("log").Key("LOG_MQTT_SUBSCRIBE_ENABLE").Bool()
	LOG_MQTT_PUBLISH_ENABLE, _ = cfg.Section("log").Key("LOG_MQTT_PUBLISH_ENABLE").Bool()
	LOG_PATH = cfg.Section("log").Key("LOG_PATH").String()
	LOG_LEVEL = cfg.Section("log").Key("LOG_LEVEL").String()
	LOG_AGE, _ = cfg.Section("log").Key("LOG_AGE").Int()
	LOG_SIZE, _ = cfg.Section("log").Key("LOG_SIZE").Int()
	LOG_BACKUP_COUNT, _ = cfg.Section("log").Key("LOG_BACKUP_COUNT").Int()
	LOG_STORE = cfg.Section("log").Key("LOG_STORE").String()

	switch PRODUCT_NAME {
	case "T320M", "T320MX", "T320MX-U":
		MQTT_CLIENT_KEEPALIVE = 60
		MQTT_CLIENT_RECONNECT_INTERVAL = 30
		MQTT_CLIENT_RECONNECT_COUNT = 20
		MQTT_CLIENT_SLEEP_INTERVAL = 60
		DEVICE_SN_LEN = 20
		SUB_TOPIC = []string{
			"/sys/%s/%s/thing/model/down_raw",
			"/ota/device/upgrade/%s/%s",
		}
		PUB_TOPIC = []string{
			"/sys/%s/%s/thing/model/up_raw",
			"/ota/device/progress/%s/%s",
		}
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
	case "示例产品-mqtt":
		SUB_TOPIC = []string{
			"v5/%s/%s/sys/property/down/#",
			"v5/%s/%s/sys/service/invoke/#",
		}
		PUB_TOPIC = []string{
			"v5/%s/%s/sys/property/up",
			"v5/%s/%s/sys/event/up",
		}
		PRODUCT_KEY = "mqtt0001"
		MQTT_CLIENT_USERNAME = PRODUCT_KEY + "&%s"
		MQTT_CLIENT_PASSWORD = DEVICE_KEY
	case "自定义":
	}

	DEVICE_SN_LEFT_LEN = strconv.FormatInt(int64(DEVICE_SN_LEN-len(DEVICE_SN_PRE)-len(DEVICE_SN_MID)), 10)
	DevFinshedSub = make(chan int, 100)
	DevFinshedConn = make(chan int, 100)
	return cfg
}

func SetConfig(jsonString string) {
	type config struct {
		PPS     int `json:"PPS"`
		PPS_PER int `json:"PPS_PER"`
		// MQTT_SUBSCRIBE_ENABLE                bool   `json:"MQTT_SUBSCRIBE_ENABLE"`
		// MQTT_PUBLISH_ENABLE                  bool   `json:"MQTT_PUBLISH_ENABLE"`
		MQTT_SERVER_HOST                     string `json:"MQTT_SERVER_HOST"`
		MQTT_SERVER_PORT                     string `json:"MQTT_SERVER_PORT"`
		PRODUCT_NAME                         string `json:"PRODUCT_NAME"`
		MQTT_CLIENT_CONNECT_INTERVAL         int    `json:"MQTT_CLIENT_CONNECT_INTERVAL"`
		MQTT_CLIENT_CONNECT_PER_100_INTERVAL int    `json:"MQTT_CLIENT_CONNECT_PER_100_INTERVAL"`
		DEVICE_SN_MID                        string `json:"DEVICE_SN_MID"`
		DEVICE_SN_SUF_START_BY               int    `json:"DEVICE_SN_SUF_START_BY"`
		DEVICE_TOTAL_COUNT                   int    `json:"DEVICE_TOTAL_COUNT"`
		DEVICE_KEY                           string `json:"DEVICE_KEY"`
		MQTT_CLIENT_USERNAME                 string `json:"MQTT_CLIENT_USERNAME"`
		MQTT_CLIENT_PASSWORD                 string `json:"MQTT_CLIENT_PASSWORD"`
		MQTT_CLIENT_ID                       string `json:"MQTT_CLIENT_ID"`
		MQTT_CLIENT_KEEPALIVE                int    `json:"MQTT_CLIENT_KEEPALIVE"`
		MQTT_CLIENT_RECONNECT_INTERVAL       int    `json:"MQTT_CLIENT_RECONNECT_INTERVAL"`
		MQTT_CLIENT_RECONNECT_COUNT          int    `json:"MQTT_CLIENT_RECONNECT_COUNT"`
		MQTT_CLIENT_SLEEP_INTERVAL           int    `json:"MQTT_CLIENT_SLEEP_INTERVAL"`
		PRODUCT_KEY                          string `json:"PRODUCT_KEY"`
		DEVICE_SN_PRE                        string `json:"DEVICE_SN_PRE"`
		DEVICE_SN_LEN                        int    `json:"DEVICE_SN_LEN"`
		MQTT_SUBSCRIBE_TOPIC                 string `json:"MQTT_SUBSCRIBE_TOPIC"`
		MQTT_PUBLISH_TOPIC                   string `json:"MQTT_PUBLISH_TOPIC"`

		PROPERTY_UP_ENABLE bool   `json:"PROPERTY_UP_ENABLE"`
		BOOLKey            bool   `json:"BOOLKey"`
		DoubleKey          bool   `json:"DoubleKey"`
		LongKey            bool   `json:"LongKey"`
		StringKey          bool   `json:"StringKey"`
		PowerSwitch        bool   `json:"PowerSwitch"`
		CustomKey          bool   `json:"CustomKey"`
		CustomValue        string `json:"CustomValue"`
	}
	var c config
	err := json.Unmarshal([]byte(jsonString), &c)
	if err == nil {
		PPS = c.PPS
		PPS_PER = c.PPS_PER
		// MQTT_SUBSCRIBE_ENABLE = c.MQTT_SUBSCRIBE_ENABLE
		// MQTT_PUBLISH_ENABLE = c.MQTT_PUBLISH_ENABLE
		MQTT_SERVER_HOST = c.MQTT_SERVER_HOST
		MQTT_SERVER_PORT = c.MQTT_SERVER_PORT
		MQTT_CLIENT_BROKER = fmt.Sprintf("tcp://%s:%s", MQTT_SERVER_HOST, MQTT_SERVER_PORT)
		PRODUCT_NAME = c.PRODUCT_NAME
		MQTT_CLIENT_CONNECT_INTERVAL = c.MQTT_CLIENT_CONNECT_INTERVAL
		if MQTT_CLIENT_CONNECT_INTERVAL < 20 {
			MQTT_CLIENT_CONNECT_INTERVAL = 20
		}
		MQTT_CLIENT_CONNECT_PER_100_INTERVAL = c.MQTT_CLIENT_CONNECT_PER_100_INTERVAL
		DEVICE_SN_MID = c.DEVICE_SN_MID
		DEVICE_SN_SUF_START_BY = c.DEVICE_SN_SUF_START_BY
		DEVICE_TOTAL_COUNT = c.DEVICE_TOTAL_COUNT
		DEVICE_KEY = c.DEVICE_KEY
		MQTT_CLIENT_USERNAME = c.MQTT_CLIENT_USERNAME
		MQTT_CLIENT_PASSWORD = c.MQTT_CLIENT_PASSWORD
		MQTT_CLIENT_KEEPALIVE = c.MQTT_CLIENT_KEEPALIVE
		MQTT_CLIENT_RECONNECT_INTERVAL = c.MQTT_CLIENT_RECONNECT_INTERVAL
		MQTT_CLIENT_RECONNECT_COUNT = c.MQTT_CLIENT_RECONNECT_COUNT
		MQTT_CLIENT_SLEEP_INTERVAL = c.MQTT_CLIENT_SLEEP_INTERVAL
		PRODUCT_KEY = c.PRODUCT_KEY
		DEVICE_SN_PRE = c.DEVICE_SN_PRE
		DEVICE_SN_LEN = c.DEVICE_SN_LEN
		SUB_TOPIC = strings.Split(c.MQTT_SUBSCRIBE_TOPIC, ";")
		PUB_TOPIC = strings.Split(c.MQTT_PUBLISH_TOPIC, ";")

		PROPERTY_UP_ENABLE = c.PROPERTY_UP_ENABLE
		BOOLKey = c.BOOLKey
		DoubleKey = c.DoubleKey
		LongKey = c.LongKey
		StringKey = c.StringKey
		PowerSwitch = c.PowerSwitch
		CustomKey = c.CustomKey
		CustomValue = c.CustomValue
		switch PRODUCT_NAME {
		case "T320M", "T320MX", "T320MX-U":
			MQTT_CLIENT_KEEPALIVE = 60
			MQTT_CLIENT_RECONNECT_INTERVAL = 30
			MQTT_CLIENT_RECONNECT_COUNT = 20
			MQTT_CLIENT_SLEEP_INTERVAL = 60
			DEVICE_SN_LEN = 20
			SUB_TOPIC = []string{
				"/sys/%s/%s/thing/model/down_raw",
				"/ota/device/upgrade/%s/%s",
			}
			PUB_TOPIC = []string{
				"/sys/%s/%s/thing/model/up_raw",
				"/ota/device/progress/%s/%s",
			}
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
		case "示例产品-mqtt":
			SUB_TOPIC = []string{
				"v5/%s/%s/sys/property/down/#",
				"v5/%s/%s/sys/service/invoke/#",
			}
			PUB_TOPIC = []string{
				"v5/%s/%s/sys/property/up",
				"v5/%s/%s/sys/event/up",
			}
			PRODUCT_KEY = "mqtt0001"
			MQTT_CLIENT_USERNAME = PRODUCT_KEY + "&%s"
			MQTT_CLIENT_PASSWORD = DEVICE_KEY
		case "自定义":
			MQTT_CLIENT_ID = c.MQTT_CLIENT_ID
		}
		DEVICE_SN_LEFT_LEN = strconv.FormatInt(int64(DEVICE_SN_LEN-len(DEVICE_SN_PRE)-len(DEVICE_SN_MID)), 10)
	}
}
