package DeviceManager

import (
	"github.com/fsnotify/fsnotify"
)

// 查找设备：DeviceGetCount()+DeviceGetHandleByIndex(Index int)  - GetDeviceList()返回逻辑ID列表
// 产品名: nvml.DeviceGetBrand - GetDevType
// ID: nvml.DeviceGetUUID - GetDeviceList()
// Index: 递增
// 内存: DeviceGetMemoryInfo - GetDeviceMemoryInfo
// 健康状态:  /rm/health.go - GetDeviceHealth
type Device interface {
	GetBrand() string
	GetID() int
	GetUUlD() string
	GetIndex() int
	GetModel() string
	GetMemory() int
	IsHealthy() bool
	IsAllocated() bool
	SetAllocated(bool)
}

type CommonDevice struct {
	Brand     string
	Id        int
	Uuid      string
	Index     int
	Memory    int
	Health    bool
	Allocated bool
}

type InterfaceServer interface {
	Name()
	Start(*fsnotify.Watcher) error
	Stop()
}

type ResourceManager interface {
	Start(stopCh <-chan interface{}) error
	Stop() error
	Name() string
	Vendor() string
	Devices() []*CommonDevice
	//RegisterHealthMonitor(name string, monitorCh chan<- *CommonDevice) error
	//GetDeviceUtilByName(name string) (int, uint64, error)
}
