package DevicePlugin

import (
	"DevicePluginFramework/DeviceManager"
	"fmt"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type NvidiaPlugin struct {
	AiDevicePluginServer
}

func (p *NvidiaPlugin) generateGPU(devs []DeviceManager.Device) []*pluginapi.Device {
	// ...
	apiDevs := make([]*pluginapi.Device, len(devs))
	for _, dev := range devs {
		if !dev.IsAllocated() && dev.IsHealthy() {
			apiDevs = append(apiDevs, &pluginapi.Device{
				ID: fmt.Sprintf("%s-%s",
					dev.GetBrand(),
					dev.GetID()),
				Health: "healthy",
			})
		}
	}
	return apiDevs
}

func (p *NvidiaPlugin) allocateGPU(req *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	resps := make([]*pluginapi.ContainerAllocateResponse, len(req.ContainerRequests))
	for i, req := range req.ContainerRequests {
		for _, devUUID := range req.DevicesIDs {
			dev, _ := p.cache.GetDevice(devUUID)
			if d, ok := dev.(DeviceManager.Device); ok {
				d.SetAllocated(true)
			}
			p.cache.SetDevice(devUUID, dev, -1)
		}
		resps[i] = &pluginapi.ContainerAllocateResponse{
			// ...
		}
	}
	return &pluginapi.AllocateResponse{ContainerResponses: resps}, nil
}
