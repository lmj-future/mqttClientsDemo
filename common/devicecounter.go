package common

import "sync"

var Counter *DeviceCounter

func InitCounter() {
	Counter = NewDeviceCounter()
}

type DeviceCounter struct {
	counters map[string]int
	mutex    sync.Mutex
}

func NewDeviceCounter() *DeviceCounter {
	return &DeviceCounter{
		counters: make(map[string]int),
		mutex:    sync.Mutex{},
	}
}

func (dc *DeviceCounter) Increment(deviceSN string) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	dc.counters[deviceSN]++
}

func (dc *DeviceCounter) GetCount(deviceSN string) int {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	return dc.counters[deviceSN]
}

func (dc *DeviceCounter) IncrementAndGet(deviceSN string) int {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	dc.counters[deviceSN]++
	return dc.counters[deviceSN]
}
