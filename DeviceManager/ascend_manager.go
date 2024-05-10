package DeviceManager

import (
	"gitee.com/ascend/ascend-npu-exporter/v6/devmanager"
)

type AscendManager struct {
	vendor        string
	devManager    devmanager.DeviceInterface
	devices       map[string]*Device
	healthMonitor map[string]chan<- *Device
	stopCh        chan<- interface{}
}
