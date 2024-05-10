package DeviceManager

import "k8s.io/klog/v2"

func (d CommonDevice) PrintInfo() {
	klog.Infof("Device %d:\n", d.Id)
	klog.Infof("  Brand: %s\n", d.Brand)
	klog.Infof("  UUID: %s\n", d.Uuid)
	klog.Infof("  Memory: %d MB\n", d.Memory/1024/1024)
}

func (d CommonDevice) GetBrand() string {
	return d.Brand
}

func (d CommonDevice) GetID() int {
	return d.Id
}

func (d CommonDevice) GetUUlD() string {
	return d.Uuid
}

func (d CommonDevice) GetIndex() int {
	return d.Index
}

func (d CommonDevice) GetModel() string {
	return ""
}

func (d CommonDevice) GetMemory() int {
	return d.Memory
}

func (d CommonDevice) IsHealthy() (bool, error) {
	return d.Health, nil
}

func (d CommonDevice) IsAllocated() bool {
	return d.Allocated
}

func (d CommonDevice) SetAllocated(b bool) {
	d.Allocated = b
}
