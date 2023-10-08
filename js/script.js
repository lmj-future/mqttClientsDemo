
document.title = "MQTT CLIENTS DEMO v1.0.2"
const hashTable = {};


var timestamp = Date.now().toString();
var config = {
    PPS: 7,
    PPS_PER: 600,
    
    // MQTT_SUBSCRIBE_ENABLE: false,
    // MQTT_PUBLISH_ENABLE: false,
 
    MQTT_SERVER_HOST: "33.33.33.244",
    MQTT_SERVER_PORT: "14005",
    PRODUCT_NAME: "T320M",
    MQTT_CLIENT_CONNECT_INTERVAL: 20,
    MQTT_CLIENT_CONNECT_PER_100_INTERVAL: 1,
    DEVICE_SN_MID: "CC",
    DEVICE_SN_SUF_START_BY: 1,
    DEVICE_TOTAL_COUNT: 1,
    DEVICE_KEY: "123456",
    
    MQTT_CLIENT_USERNAME: "EbcMXsMg+c",
    MQTT_CLIENT_PASSWORD: "GnW7YLdtWv",
    MQTT_CLIENT_ID: "",
    MQTT_CLIENT_KEEPALIVE: 60,
    MQTT_CLIENT_RECONNECT_INTERVAL: 30,
    MQTT_CLIENT_RECONNECT_COUNT: 20,
    MQTT_CLIENT_SLEEP_INTERVAL: 60,
    PRODUCT_KEY: "kiSHgWsG",
    DEVICE_SN_PRE: "219801A26U",
    DEVICE_SN_LEN: 20,
    
    MQTT_SUBSCRIBE_TOPIC: "/sys/%s/%s/thing/model/down_raw;/ota/device/upgrade/%s/%s",
    MQTT_PUBLISH_TOPIC: "/sys/%s/%s/thing/model/up_raw;/ota/device/progress/%s/%s",

    PROPERTY_UP_ENABLE: false,
    BOOLKey: false,
    DoubleKey: false,
    LongKey: false,
    StringKey: false,
    PowerSwitch: false,
    CustomKey: false,
    CustomValue: "{\"key\":\"value\"}"
}

function getConfig() {
    var e01 = document.querySelector('div[data-label="PPS"]');
    if (e01 !== null && e01 !== undefined) {config.PPS = parseInt(e01.querySelector('input').value);}
    var e02 = document.querySelector('div[data-label="PPS_PER"]');
    if (e02 !== null && e02 !== undefined) {config.PPS_PER = parseInt(e02.querySelector('input').value);}
    // var e03 = document.querySelector('div[data-label="MQTT_SUBSCRIBE_ENABLE"]');
    // if (e03 !== null && e03 !== undefined) {config.MQTT_SUBSCRIBE_ENABLE = e03.querySelector('input').value;}
    // var e04 = document.querySelector('div[data-label="MQTT_PUBLISH_ENABLE"]');
    // if (e04 !== null && e04 !== undefined) {config.MQTT_PUBLISH_ENABLE = e04.querySelector('input').value;}
    var e05 = document.querySelector('div[data-label="MQTT_SERVER_HOST"]');
    if (e05 !== null && e05 !== undefined) {config.MQTT_SERVER_HOST = e05.querySelector('input').value;}
    var e06 = document.querySelector('div[data-label="MQTT_SERVER_PORT"]');
    if (e06 !== null && e06 !== undefined) {config.MQTT_SERVER_PORT = e06.querySelector('input').value;}
    var e07 = document.querySelector('div[data-label="PRODUCT_NAME"]');
    if (e07 !== null && e07 !== undefined) {config.PRODUCT_NAME = e07.querySelector('select').value;}
    var e08 = document.querySelector('div[data-label="MQTT_CLIENT_CONNECT_INTERVAL"]');
    if (e08 !== null && e08 !== undefined) {config.MQTT_CLIENT_CONNECT_INTERVAL = parseInt(e08.querySelector('input').value);}
    var e09 = document.querySelector('div[data-label="MQTT_CLIENT_CONNECT_PER_100_INTERVAL"]');
    if (e09 !== null && e09 !== undefined) {config.MQTT_CLIENT_CONNECT_PER_100_INTERVAL = parseInt(e09.querySelector('input').value);}
    var e10 = document.querySelector('div[data-label="DEVICE_SN_MID"]');
    if (e10 !== null && e10 !== undefined) {config.DEVICE_SN_MID = e10.querySelector('input').value;}
    var e11 = document.querySelector('div[data-label="DEVICE_SN_SUF_START_BY"]');
    if (e11 !== null && e11 !== undefined) {config.DEVICE_SN_SUF_START_BY = parseInt(e11.querySelector('input').value);}
    var e12 = document.querySelector('div[data-label="DEVICE_TOTAL_COUNT"]');
    if (e12 !== null && e12 !== undefined) {config.DEVICE_TOTAL_COUNT = parseInt(e12.querySelector('input').value);}
    var e13 = document.querySelector('div[data-label="MQTT_CLIENT_USERNAME"]');
    if (e13 !== null && e13 !== undefined) {config.MQTT_CLIENT_USERNAME = e13.querySelector('input').value;}
    var e14 = document.querySelector('div[data-label="MQTT_CLIENT_PASSWORD"]');
    if (e14 !== null && e14 !== undefined) {config.MQTT_CLIENT_PASSWORD = e14.querySelector('input').value;}
    var e15 = document.querySelector('div[data-label="MQTT_CLIENT_KEEPALIVE"]');
    if (e15 !== null && e15 !== undefined) {config.MQTT_CLIENT_KEEPALIVE = parseInt(e15.querySelector('input').value);}
    var e16 = document.querySelector('div[data-label="MQTT_CLIENT_RECONNECT_INTERVAL"]');
    if (e16 !== null && e16 !== undefined) {config.MQTT_CLIENT_RECONNECT_INTERVAL = parseInt(e16.querySelector('input').value);}
    var e17 = document.querySelector('div[data-label="MQTT_CLIENT_RECONNECT_COUNT"]');
    if (e17 !== null && e17 !== undefined) {config.MQTT_CLIENT_RECONNECT_COUNT = parseInt(e17.querySelector('input').value);}
    var e18 = document.querySelector('div[data-label="MQTT_CLIENT_SLEEP_INTERVAL"]');
    if (e18 !== null && e18 !== undefined) {config.MQTT_CLIENT_SLEEP_INTERVAL = parseInt(e18.querySelector('input').value);}
    var e19 = document.querySelector('div[data-label="PRODUCT_KEY"]');
    if (e19 !== null && e19 !== undefined) {config.PRODUCT_KEY = e19.querySelector('input').value;}
    var e20 = document.querySelector('div[data-label="DEVICE_SN_PRE"]');
    if (e20 !== null && e20 !== undefined) {config.DEVICE_SN_PRE = e20.querySelector('input').value;}
    var e21 = document.querySelector('div[data-label="DEVICE_SN_LEN"]');
    if (e21 !== null && e21 !== undefined) {config.DEVICE_SN_LEN = parseInt(e21.querySelector('input').value);}
    var e22 = document.querySelector('div[data-label="MQTT_SUBSCRIBE_TOPIC"]');
    if (e22 !== null && e22 !== undefined) {config.MQTT_SUBSCRIBE_TOPIC = e22.querySelector('input').value;}
    var e23 = document.querySelector('div[data-label="MQTT_PUBLISH_TOPIC"]');
    if (e23 !== null && e23 !== undefined) {config.MQTT_PUBLISH_TOPIC = e23.querySelector('input').value;}
    var e24 = document.querySelector('div[data-label="DEVICE_KEY"]');
    if (e24 !== null && e24 !== undefined) {config.DEVICE_KEY = e24.querySelector('input').value;}
    var e25 = document.querySelector('div[data-label="PROPERTY_UP_ENABLE"]');
    if (e25 !== null && e25 !== undefined) {config.PROPERTY_UP_ENABLE = e25.querySelector('input').checked;}
    var e26 = document.querySelector('div[data-label="BOOLKey"]');
    if (e26 !== null && e26 !== undefined) {config.BOOLKey = e26.querySelector('input').checked;}
    var e27 = document.querySelector('div[data-label="DoubleKey"]');
    if (e27 !== null && e27 !== undefined) {config.DoubleKey = e27.querySelector('input').checked;}
    var e28 = document.querySelector('div[data-label="LongKey"]');
    if (e28 !== null && e28 !== undefined) {config.LongKey = e28.querySelector('input').checked;}
    var e29 = document.querySelector('div[data-label="StringKey"]');
    if (e29 !== null && e29 !== undefined) {config.StringKey = e29.querySelector('input').checked;}
    var e30 = document.querySelector('div[data-label="PowerSwitch"]');
    if (e30 !== null && e30 !== undefined) {config.PowerSwitch = e30.querySelector('input').checked;}
    var e31 = document.querySelector('div[data-label="CustomKey"]');
    if (e31 !== null && e31 !== undefined) {config.CustomKey = e31.querySelector('input').checked;}
    var e32 = document.querySelector('div[data-label="CustomValue"]');
    if (e32 !== null && e32 !== undefined) {config.CustomValue = e32.querySelector('input').value;}
    var e33 = document.querySelector('div[data-label="MQTT_CLIENT_ID"]');
    if (e33 !== null && e33 !== undefined) {config.MQTT_CLIENT_ID = e33.querySelector('input').value;}

    return JSON.stringify(config);
}

function getTimestamp() {
    return timestamp;
}
function setTimestamp() {
    timestamp = Date.now().toString();
}
function getStartButton() {
    return document.querySelector('div[data-label="开始按钮"]');
}
function getStopButton() {
    return document.querySelector('div[data-label="停止按钮"]');
}
function getTerminalJoinOrLeave() {
    return document.getElementById("u247_input");
}
function getTerminalId() {
    return document.getElementById("u250_input").value;
}

//绑定点击信息
//这里还需要发送对应终端的id过去
function changeTerminalJoiToLeave() {
    var TerButton = getTerminalJoinOrLeave();
    if (TerButton !== null && TerButton !== undefined) {   
        TerButton.addEventListener('click', function() {       
            var TerId = getTerminalId()
            var xhr = new XMLHttpRequest();
            xhr.open("GET","http://localhost:7777/api/data?id=" + TerId + "&isLeave="+ TerButton.value ,true)
            changeButtonStatus(); 
            xhr.send(); 
        });
    }
}

function disableStartButton() {
    var startButton = getStartButton();
    if (startButton !== null && startButton !== undefined) {
        startButton.style.pointerEvents = 'none';  // 禁用鼠标事件
        startButton.style.opacity = '0.5';  // 设置按钮透明度为0.5，呈现灰色效果
    }
}
function enableStartButton() {
    var startButton = getStartButton();
    if (startButton !== null && startButton !== undefined) {
        startButton.style.pointerEvents = 'auto';  // 恢复鼠标事件
        startButton.style.opacity = '1';  // 恢复按钮透明度
    }
}
function disableStopButton() {
    var stopButton = getStopButton();
    if (stopButton !== null && stopButton !== undefined) {
        stopButton.style.pointerEvents = 'none';  // 禁用鼠标事件
        stopButton.style.opacity = '0.5';  // 设置按钮透明度为0.5，呈现灰色效果
    }
}
//进入界面的时候禁用终端入网
function disableTerminalJoinOrLeaveButton() {
    var disButton = getTerminalJoinOrLeave();
    if (disButton !== null && disButton !== undefined) {
        disButton.style.pointerEvents = 'none';  // 禁用鼠标事件
        disButton.style.opacity = '0.5';  // 设置按钮透明度为0.5，呈现灰色效果
    }
}
function enableTerminalJoinOrLeaveButton() {
    var enableButton = getTerminalJoinOrLeave();
    if (enableButton !== null && enableButton !== undefined) {
        enableButton.style.pointerEvents = 'auto';  // 禁用鼠标事件
        enableButton.style.opacity = '1';  // 设置按钮透明度为0.5，呈现灰色效果
    }
}

function enableStopButton() {
    var stopButton = getStopButton();
    if (stopButton !== null && stopButton !== undefined) {
        stopButton.style.pointerEvents = 'auto';  // 恢复鼠标事件
        stopButton.style.opacity = '1';  // 恢复按钮透明度
    }
}
function startButtonEvent() {
    var startButton = getStartButton();
    if (startButton !== null && startButton !== undefined) {
        startButton.addEventListener('click', function() {
            // 处理开始按钮的点击事件
            disableStartButton();
            enableTerminalJoinOrLeaveButton();
            setTimestamp();
            start(getConfig(), getTimestamp());
            enableStopButton();
            initDropList();
            initTerminalStatus();
        });
    }
}
function stopButtonEvent() {
    var stopButton = getStopButton();
    if (stopButton !== null && stopButton !== undefined) {
        stopButton.addEventListener('click', function() {
            // 处理停止按钮的点击事件
            disableTerminalJoinOrLeaveButton();
            disableStopButton();
            stop(getTimestamp());
            enableStartButton();
        });
    }
}

//更改状态为终端离网
function changeButtonStatus() {
    //首先从下拉列表中获取
    var terminalId = document.getElementById("u250_input")
    var JoinOrLeaveButton = document.getElementById("u247_input")
    if (terminalId !== null && terminalId !== undefined) {
        var JoinOrLeavStatus =  hashTable[terminalId.value];
        if (JoinOrLeavStatus !== null && JoinOrLeavStatus !== undefined) {
            hashTable[terminalId.value] = JoinOrLeavStatus === "模拟终端入网" ? "模拟终端离网" : "模拟终端入网" 
            JoinOrLeaveButton.value = JoinOrLeavStatus === "模拟终端入网" ? "模拟终端离网" : "模拟终端入网" 
        }
    }
}

function initTerminalStatus() {
    var sz = document.getElementById("u155_input")
    if (sz !== null && sz !== undefined) {
        for (var i = 1; i <= sz.value; i ++) {
            hashTable[i +""] = "模拟终端入网"
        }
    }
}

function getSubLogCheckboxEvent() {
    var div = document.querySelector('div[data-label="MQTT_SUBSCRIBE_ENABLE"]')
    if (div !== null && div !== undefined) {return div.querySelector('input');}
    return null;
}
function getPubLogCheckboxEvent() {
    var div = document.querySelector('div[data-label="MQTT_PUBLISH_ENABLE"]')
    if (div !== null && div !== undefined) {return div.querySelector('input');}
    return null;
}
function subLogCheckboxEvent() {
    var checkBox = getSubLogCheckboxEvent();
    if (checkBox !== null && checkBox !== undefined) {
        checkBox.addEventListener("change", function() {
            if (checkBox.checked) {
                subLog(true);
            } else {
                subLog(false);
            }
        });
    }
}
function pubLogCheckboxEvent() {
    var checkBox = getPubLogCheckboxEvent();
    if (checkBox !== null && checkBox !== undefined) {
        checkBox.addEventListener("change", function() {
            if (checkBox.checked) {
                pubLog(true);
            } else {
                pubLog(false);
            }
        });
    }
}
function getTelnetCheckboxEvent() {
    var div = document.querySelector('div[data-label="TELNET"]')
    if (div !== null && div !== undefined) {return div.querySelector('input');}
    return null;
}
function telnetCheckboxEvent() {
    var checkBox = getTelnetCheckboxEvent();
    if (checkBox !== null && checkBox !== undefined) {
        checkBox.addEventListener("change", function() {
            if (checkBox.checked) {
                telnet(true);
            } else {
                telnet(false);
            }
        });
    }
}

function initDropList() {
    var selectElement = document.getElementById("u250_input");
    var lengthPre = selectElement.options.length;
    for (var i = lengthPre - 1; i >= 0; i --) {
        selectElement.remove(i);
    }
    var len = document.getElementById("u155_input")
    for (var i = 1; i <= len.value; i ++) {
        var option = document.createElement("option");
        option.text = i;
        selectElement.add(option);
    }
}

function changeWithSelect() {
    var selectElement = document.getElementById('u250_input');
    var TerminalButton = document.getElementById('u247_input');
    selectElement.addEventListener('change', function() {
        TerminalButton.value = hashTable[selectElement.value]
    });
}

enableStartButton();
disableStopButton();
startButtonEvent();
stopButtonEvent();
subLogCheckboxEvent();
pubLogCheckboxEvent();
telnetCheckboxEvent();
disableTerminalJoinOrLeaveButton();
changeTerminalJoiToLeave();
changeWithSelect();