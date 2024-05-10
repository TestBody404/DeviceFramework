package Cache

import (
	"DevicePluginFramework/DeviceManager"
	"DevicePluginFramework/pb"
	"context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

const grpcAddress string = "/tmp/device_cache_socket"

// DeviceCache is a cache for storing device objects.
type DeviceCache struct {
	cache      *ConcurrencyLRUCache
	brand      string
	mu         sync.Mutex
	grpcClient pb.QueryDevicesServeClient // 设备管理模块的 gRPC 客户端
}

// NewDeviceCache creates a new instance of DeviceCache with the given maxEntries.
func NewDeviceCache(maxEntries int) *DeviceCache {
	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial gRPC server: %v", err)
	}
	client := pb.NewQueryDevicesServeClient(conn)
	return &DeviceCache{
		cache:      New(maxEntries),
		grpcClient: client,
	}
}

// StartPeriodicGrpcCall 启动周期性的 gRPC 调用请求，用于更新缓存中的数据。
func (dc *DeviceCache) StartPeriodicGrpcCall(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				dc.updateCacheFromGrpc()
			}
		}
	}()
}

// updateCacheFromGrpc 从 gRPC 调用请求中更新缓存中的数据。
func (dc *DeviceCache) updateCacheFromGrpc() {
	resp, err := dc.grpcClient.QueryDevices(context.Background(), &pb.QueryDevicesRequest{})
	if err != nil {
		log.Printf("Failed to get devices from gRPC: %v", err)
		return
	}

	// Update cache
	dc.mu.Lock()
	defer dc.mu.Unlock()

	for _, pbDevice := range resp.Devices {
		dc.brand = pbDevice.Brand
		device := buildDevice(*pbDevice)
		dc.SetDevice(device.Uuid, device, time.Duration(negInt64One)) // 设置过期时间为-1，即永不过期
	}
}

func buildDevice(pbDevice pb.CommonDevice) DeviceManager.CommonDevice {
	device := DeviceManager.CommonDevice{
		Brand:     pbDevice.Brand,
		Id:        int(pbDevice.Id),
		Uuid:      pbDevice.Uuid,
		Index:     int(pbDevice.Index),
		Memory:    int(pbDevice.Memory),
		Health:    pbDevice.Health,
		Allocated: pbDevice.Allocated,
	}
	return device
}

// SetDevice adds or updates a device in the cache.
func (dc *DeviceCache) SetDevice(key string, device interface{}, expireTime time.Duration) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return dc.cache.Set(key, device, expireTime)
}

// GetDevice retrieves a device from the cache by its key.
func (dc *DeviceCache) GetDevice(key string) (interface{}, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return dc.cache.Get(key)
}

// DeleteDevice deletes a device from the cache by its key.
func (dc *DeviceCache) DeleteDevice(key string) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.cache.Delete(key)
}

// SetDeviceIfNotExists sets a device in the cache if it does not already exist.
func (dc *DeviceCache) SetDeviceIfNotExists(key string, device interface{}, expireTime time.Duration) bool {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return dc.cache.SetIfNX(key, device, expireTime)
}
