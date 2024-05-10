package DevicePlugin

import (
	"DevicePluginFramework/Cache"
	"DevicePluginFramework/DeviceManager"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"log"
	"net"
	"os"
	"path"
	"time"
)

const (
	resourceName           = "github.com/AiDevice"
	serverSock             = pluginapi.DevicePluginPath + "AiDevice.sock"
	envDisableHealthChecks = "DP_DISABLE_HEALTHCHECKS"
	allHealthChecks        = "xids"
)

type devGenerateFunc func([]*DeviceManager.Device) []*pluginapi.Device
type devAllocateFunc func(*pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error)

// AiDevicePluginServer implements the Kubernetes device plugin API
type AiDevicePluginServer struct {
	devs   []*pluginapi.Device
	socket string

	stop   chan interface{}
	health chan *pluginapi.Device

	server     *grpc.Server
	cache      Cache.DeviceCache
	generateFn devGenerateFunc
	allocateFn devAllocateFunc
}

func NewAiDevicePluginServer(cache Cache.DeviceCache) *AiDevicePluginServer {
	return &AiDevicePluginServer{
		devs:   getDevices(cache),
		socket: serverSock,
		stop:   make(chan interface{}),
		health: make(chan *pluginapi.Device),
	}
}

func (m *AiDevicePluginServer) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{}, nil
}

func (m *AiDevicePluginServer) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}

func (m *AiDevicePluginServer) GetPreferredAllocation(context.Context, *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return &pluginapi.PreferredAllocationResponse{}, nil
}

// Serve starts the gRPC server and register the device plugin to Kubelet
func (m *AiDevicePluginServer) serve() error {
	err := m.start()
	if err != nil {
		log.Printf("Could not start device plugin: %s", err)
		return err
	}
	log.Println("Starting to serve on", m.socket)
	err = m.register(pluginapi.KubeletSocket, resourceName)
	if err != nil {
		log.Printf("Could not register device plugin: %s", err)
		m.Stop()
		return err
	}
	log.Println("Registered device plugin with Kubelet")
	return nil
}

// Start starts the gRPC server of the device plugin
func (m *AiDevicePluginServer) start() error {
	err := m.cleanup() // 清理旧socket
	if err != nil {
		return err
	}
	sock, err := net.Listen("unix", m.socket) // 创建监听器（会创建新socket）
	if err != nil {
		return err
	}
	m.server = grpc.NewServer([]grpc.ServerOption{}...) // 创建gRPC服务
	pluginapi.RegisterDevicePluginServer(m.server, m)   // 把m注册为服务器端点
	// 启动服务器，开始监听socket
	go m.server.Serve(sock)
	// 建立到套接字的gRPC连接
	conn, err := grpc.Dial(
		"unix://"+m.socket,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return err
	}
	defer conn.Close()
	go m.healthcheck()
	return nil
}

// Stop stops the gRPC server
func (m *AiDevicePluginServer) Stop() error {
	if m.server == nil {
		return nil
	}
	m.server.Stop()
	m.server = nil
	close(m.stop)
	return m.cleanup()
}

// Register registers the device plugin for the given resourceName with Kubelet.
func (m *AiDevicePluginServer) register(kubeletEndpoint, resourceName string) error {
	conn, err := grpc.Dial(
		"unix://"+kubeletEndpoint,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pluginapi.NewRegistrationClient(conn)
	reqt := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     path.Base(m.socket),
		ResourceName: resourceName,
	}

	_, err = client.Register(context.Background(), reqt)
	if err != nil {
		return err
	}
	return nil
}

// ListAndWatch lists devices and update that list according to the health status
func (m *AiDevicePluginServer) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	s.Send(&pluginapi.ListAndWatchResponse{Devices: m.devs})

	for {
		select {
		case <-m.stop:
			return nil
		case d := <-m.health:
			d.Health = pluginapi.Unhealthy
			s.Send(&pluginapi.ListAndWatchResponse{Devices: m.devs})
		}
	}
}

// Allocate which return list of devices.
func (m *AiDevicePluginServer) Allocate(ctx context.Context, reqs *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	devs := m.devs
	var responses pluginapi.AllocateResponse

	for _, req := range reqs.ContainerRequests {
		for _, id := range req.DevicesIDs {
			if !deviceExists(devs, id) {
				return nil, fmt.Errorf("invalid allocation request: unknown device: %s", id)
			}

			response := new(pluginapi.ContainerAllocateResponse)
			response.Devices = []*pluginapi.DeviceSpec{
				{
					ContainerPath: "/dev/Ai",
					HostPath:      "/dev/Ai",
					Permissions:   "rwm",
				},
			}

			responses.ContainerResponses = append(responses.ContainerResponses, response)
		}
	}

	return &responses, nil
}

// 清理本插件的套接字文件
func (m *AiDevicePluginServer) cleanup() error {
	if err := os.Remove(m.socket); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func (m *AiDevicePluginServer) healthcheck() {
	for range m.stop {
		return
	}
}
