package DeviceManager

import (
	"DevicePluginFramework/pb"
	"context"
	"fmt"
	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
	"github.com/NVIDIA/go-nvlib/pkg/nvml"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"
	"net"
	"strconv"
	"time"
)

type NvidiaManager struct {
	vendor            string
	devManager        nvml.Interface
	devices           map[string]*CommonDevice
	healthMonitor     map[string]chan<- *CommonDevice
	stopCh            chan<- interface{}
	skipMig           bool
	healthCheckTicker *time.Ticker
	queryServer       pb.QueryDevicesServeServer
}

// NewNvidiaManager creates a new NvidiaManager.
func NewNvidiaManager() *NvidiaManager {
	return &NvidiaManager{
		vendor:        "NVIDIA",
		devices:       make(map[string]*CommonDevice),
		healthMonitor: make(map[string]chan<- *CommonDevice),
	}
}

// Start initializes the NvidiaManager.
func (m *NvidiaManager) Start(stopCh <-chan interface{}) error {
	// Initialize NVML library
	if err := m.devManager.Init(); err != nvml.SUCCESS {
		klog.Fatalf("Failed to initialize NVML: %v", err)
		return err
	}
	defer m.devManager.Shutdown()
	// Discover devices
	m.discoverDevices()
	// Start health check
	go m.startHealthCheck(stopCh)
	// Start info check
	go m.startInfoCheck(stopCh)
	// Start gRPC server
	m.startQueryDeviceServe()
	return nil
}

// discoverDevices discovers NVIDIA GPU devices.
func (m *NvidiaManager) discoverDevices() {
	// Get GPU count
	count, err := m.devManager.DeviceGetCount()
	if err != nvml.SUCCESS {
		klog.Fatalf("Failed to get device count: %v", err)
	}
	klog.Infof("Found %d NVIDIA GPU(s)\n", count)
	// Iterate through each GPU device
	for i := 0; i < count; i++ {
		err := m.initializeDevice(i)
		if err != nil {
			klog.Errorf("Failed to initialize device %d: %v", i, err)
		}
	}
}

// initializeDevice initializes a GPU device with the given index.
func (m *NvidiaManager) initializeDevice(index int) error {
	device, err := newDevice(index, m.devManager)
	if err != nil {
		return fmt.Errorf("failed to create device object for device %d: %v", index, err)
	}
	// Add device to manager
	m.devices[device.Uuid] = device
	// Add monitor
	ch := make(chan *CommonDevice)
	m.healthMonitor[device.Uuid] = ch
	// Print device info
	device.PrintInfo()
	return nil
}

// newDevice creates a new Device object with the given ID.
func newDevice(id int, devManager nvml.Interface) (*CommonDevice, error) {
	nvmlDevice, err := devManager.DeviceGetHandleByIndex(id)
	if err != nvml.SUCCESS {
		return nil, fmt.Errorf("failed to get CommonDevice information for CommonDevice %d: %v", id, err)
	}
	UUID, err := nvmlDevice.GetUUID()
	Memory, err := nvmlDevice.GetMemoryInfo()
	device := &CommonDevice{
		Id:        id,
		Brand:     "NVIDIA",
		Uuid:      UUID,
		Memory:    int(Memory.Free),
		Allocated: false,
		Health:    true,
	}

	return device, nil
}

// StopHealthCheck Stop the healthcheck process
func (m *NvidiaManager) StopHealthCheck() {
	if m.healthCheckTicker != nil {
		m.healthCheckTicker.Stop()
	}
}

func (m *NvidiaManager) startHealthCheck(stopCh <-chan interface{}) {
	// Start health check ticker
	m.healthCheckTicker = time.NewTicker(5 * time.Second)
	defer m.healthCheckTicker.Stop()
	// Set dcgm
	dcgm.Init(dcgm.Embedded)
	defer dcgm.Shutdown()

	for {
		select {
		case <-stopCh:
			// Received stop signal, exit the goroutine
			return
		case <-m.healthCheckTicker.C:
			// Monitor health status
			m.monitorHealth()
		}
	}
}

// monitorHealth monitors the health status of all devices.
func (m *NvidiaManager) monitorHealth() {
	for _, device := range m.devices {
		go func(device *CommonDevice) {
			UUID, err := strconv.ParseUint(device.GetUUlD(), 10, 64)
			health, err := dcgm.HealthCheckByGpuId(uint(UUID))
			if err != nil {
				klog.Errorf("Failed to get health status for device %s: %v", UUID, err)
				return
			}
			klog.Infof("Device %s health status: %s\n", UUID, health)
			oldStatus := device.Health
			// Update device health
			if health.Status == "Healthy" {
				device.Health = true
			} else {
				device.Health = false
			}
			// Notify health change
			if ch, ok := m.healthMonitor[device.Uuid]; ok && oldStatus != device.Health {
				ch <- device
			}
		}(device)
	}
}

// startInfoCheck periodically scans the device status and updates it.
func (m *NvidiaManager) startInfoCheck(stopCh <-chan interface{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-stopCh:
			return // Received stop signal, exit the goroutine
		case <-ticker.C:
			// Perform device status check
			for _, device := range m.devices {
				uuid := device.Uuid
				nvmlDevice, _ := m.devManager.DeviceGetHandleByUUID(uuid)
				Memory, _ := nvmlDevice.GetMemoryInfo()
				device.Memory = int(Memory.Free)
			}
		}
	}
}

// Devices returns a slice of Device objects.
func (m *NvidiaManager) Devices() []*CommonDevice {
	var devices []*CommonDevice
	for _, device := range m.devices {
		devices = append(devices, device)
	}
	return devices
}

// QueryDevices implements the QueryDevices gRPC method.
func (m *NvidiaManager) QueryDevices(ctx context.Context, req *pb.QueryDevicesRequest) (*pb.QueryDevicesResponse, error) {
	// 构建返回的设备列表
	var devices []*pb.CommonDevice
	for _, device := range m.devices {
		pbDevice := &pb.CommonDevice{
			Brand:     device.Brand,
			Id:        int32(device.Id),
			Uuid:      device.Uuid,
			Index:     int32(device.Index),
			Memory:    int32(device.Memory),
			Health:    device.Health,
			Allocated: device.Allocated,
		}
		devices = append(devices, pbDevice)
	}
	// 构建响应
	resp := &pb.QueryDevicesResponse{
		Devices: devices,
	}
	return resp, nil
}

// startGRPCServer starts the gRPC server.
func (m *NvidiaManager) startQueryDeviceServe() {
	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer()
	// 注册服务
	pb.RegisterQueryDevicesServeServer(grpcServer, m)
	// 监听 Unix 域套接字
	listener, err := net.Listen("unix", "/tmp/query_device_socket")
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()
	// 启动 gRPC 服务
	klog.Info("Starting gRPC server...")
	if err := grpcServer.Serve(listener); err != nil {
		klog.Fatalf("Failed to serve: %v", err)
	}
}

func (m *NvidiaManager) Stop() error {
	return nil
}

func (m *NvidiaManager) Name() string {
	return "NvidiaManager"
}

func (m *NvidiaManager) Vendor() string {
	return "Nvidia"
}
