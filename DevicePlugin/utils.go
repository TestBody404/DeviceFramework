package DevicePlugin

import (
	"DevicePluginFramework/Cache"
	"github.com/fsnotify/fsnotify"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"log"
	"os"
	"os/signal"
)

func getDevices(cache Cache.DeviceCache) []*pluginapi.Device {
	//hostname, _ := os.Hostname()
	devs := []*pluginapi.Device{}
	//for i := 0; i < number; i++ {
	//	devs = append(devs, &pluginapi.Device{
	//		ID:     fmt.Sprintf("AiDevice-%s-%d", hostname, i),
	//		Health: pluginapi.Healthy,
	//	})
	//}
	return devs
}

func deviceExists(devs []*pluginapi.Device, id string) bool {
	for _, d := range devs {
		if d.ID == id {
			return true
		}
	}
	return false
}

func newFSWatcher(files ...string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("new watcher meet error:%s\n", err.Error())
		return nil, err
	}

	for _, f := range files {
		err = watcher.Add(f)
		if err != nil {
			log.Printf("add watcher meet error:%s\n", err.Error())
			watcher.Close()
			return nil, err
		}
	}

	return watcher, nil
}

func newOSWatcher(sigs ...os.Signal) chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, sigs...)
	return sigChan
}
