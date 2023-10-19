package common

import (
	"strconv"
	"sync"
)

var Counter *DeviceCounter
var TransFrame *DeviceCounter

func InitCounter() {
	Counter = NewDeviceCounter()
	TransFrame = NewDeviceCounter()
}

type DeviceCounter struct {
	Counters map[string]int
	mutex    sync.Mutex
}

func NewDeviceCounter() *DeviceCounter {
	return &DeviceCounter{
		Counters: make(map[string]int),
		mutex:    sync.Mutex{},
	}
}

func (dc *DeviceCounter) Increment(deviceSN string) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()
	dc.Counters[deviceSN]++
	if dc.Counters[deviceSN] > 0x3f3f3f3f {
		dc.Counters[deviceSN] = 1
	}
}

// 自用
func (dc *DeviceCounter) IncrementString(addr string) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()
	dc.Counters[addr]++
}

func (dc *DeviceCounter) IncrementAndStringGet(deviceSN string) string {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()
	hexNum := strconv.FormatInt(int64(dc.Counters[deviceSN]), 16)
	for len(hexNum) < 4 {
		hexNum = "0" + hexNum
	}
	dc.Counters[deviceSN]++
	if dc.Counters[deviceSN] > 0x3f3f3f3f {
		dc.Counters[deviceSN] = 1
	}
	return hexNum
}

func (dc *DeviceCounter) GetCount(deviceSN string) int {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()
	return dc.Counters[deviceSN]
}

func (dc *DeviceCounter) GetCountString(deviceSN string) string {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()
	hexNum := strconv.FormatInt(int64(dc.Counters[deviceSN]), 16)
	for len(hexNum) < 4 {
		hexNum = "0" + hexNum
	}
	return hexNum
}

func (dc *DeviceCounter) IncrementAndGet(deviceSN string) int {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()
	dc.Counters[deviceSN]++
	if dc.Counters[deviceSN] > 0x3f3f3f3f {
		dc.Counters[deviceSN] = 1
	}
	return dc.Counters[deviceSN]
}
