package common

const (
	// 上行method
	DEV_CONNECT_STATUS   = "devConnectStatus"
	IOT_VER_INFO_GET     = "iotVerInfoGet"
	BASIC_INFO_UP        = "basicInfoUp"
	STATUS_UP            = "statusUp"
	STATUS_UP_BY_GENERAL = "statusUpByGeneral"
	NODE_TOPO_INFO_UP    = "nodeTopoInfoUp"
	IOT_NET_CFG_UP       = "iotNetCfgUp"
	IOT_GW_CFG_SYNC      = "iotGwCfgSync"
	IOT_GROUP_CFG_SYNC   = "iotGroupCfgSync"
	IOT_NODE_SYNC        = "iotNodeSync"
	IOT_MOD_CFG_SYNC     = "iotModCfgSync"
	IOT_LOG_CFG_GET      = "iotLogCfgGet"

	// 下行method
	SYNC_TIME        = "syncTime"
	IOT_LOG_CFG_PUSH = "iotLogCfgPush"
	SET_IOT_GW_CFG   = "setIotGwCfg"
	SET_NET_CFG      = "setNetCfg"
	SET_IOT_NODE_CFG = "setIotNodeCfg"
	SET_IOT_MOD_CFG  = "setIotModCfg"

	// 上行devOption
	DEV_UPGRADE_PROGRESS_UP = "devUpgradeProgressUp"

	// 下行devOption
	DEV_UPGRADE = "devUpgrade"

	// propertyUp
	PROPERTY_UP = "propertyUp"

	//downLoad
	DOWNLOADFILE_PATH    = "./storageFile/"
	DOWNLOADFILE_POSTFIX = ".bin"
)
