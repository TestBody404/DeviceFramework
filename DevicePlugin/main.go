package DevicePlugin

import (
	"DevicePluginFramework/Cache"
	"k8s.io/klog/v2"
	"log"
	"os"
	"syscall"

	"github.com/fsnotify/fsnotify"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

//var (
//	mountsAllowed = 5
//)

func StartDevicePluginServer(cache Cache.DeviceCache) {
	//flag.IntVar(&mountsAllowed, "mounts_allowed", 5, "maximum times the device can be mounted")
	//flag.Parse()
	log.Println("Starting")
	defer func() { log.Println("Stopped:") }()
	devicePlugin := NewAiDevicePluginServer(cache)
	devicePlugin.serve()
	go checkKubeletStatus(&devicePlugin, cache)
}

func checkKubeletStatus(devicePlugin **AiDevicePluginServer, cache Cache.DeviceCache) {
	klog.Info("Starting file watcher.")
	watcher, err := newFSWatcher(pluginapi.DevicePluginPath)
	if err != nil {
		klog.Errorf("Failed to create FS watcher: %v", err)
		os.Exit(1)
	}
	defer watcher.Close()

	klog.Info("Starting OS watcher.")
	sigs := newOSWatcher(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for restart := true; ; {
		if restart {
			// restart device plugin server
			if *devicePlugin != nil {
				(*devicePlugin).Stop()
			}
			*devicePlugin = NewAiDevicePluginServer(cache)
			if err := (*devicePlugin).serve(); err != nil {
				klog.Errorf("Could not contact Kubelet, retrying. Did you enable the device plugin feature gate? Error: %v", err)
			} else {
				restart = false
			}
		}
		select {
		case event := <-watcher.Events:
			// restart if kubeletSocket created
			if event.Name == pluginapi.KubeletSocket && event.Op&fsnotify.Create == fsnotify.Create {
				klog.Infof("inotify: %s created, restarting.", pluginapi.KubeletSocket)
				restart = true
			}
		case err := <-watcher.Errors:
			// handle file system error
			klog.Errorf("inotify: %v", err)
		case s := <-sigs:
			switch s {
			// restart if os signup
			case syscall.SIGHUP:
				klog.Info("Received SIGHUP, restarting.")
				restart = true
			default:
				klog.Infof("Received signal \"%v\", shutting down.", s)
				(*devicePlugin).Stop()
				return
			}
		}
	}
}
