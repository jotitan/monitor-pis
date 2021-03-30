package config

import (
	"encoding/json"
	"github.com/jotitan/monitor-pis/heartbeat"
	"io/ioutil"
	"log"
)

type MonitoringConfig struct {
	// url to send metrics
	Port string `json:"port"`
	Folder string `json:"folder"`
	HeartBeats []struct {
		Name string `json:"name"`
		Url string `json:"url"`
		Frequency string `json:"frequency"`
	} `json:"heartbeats"`
}

func NewMonitoringConfig(path string) (*MonitoringConfig,error) {
	c := &MonitoringConfig{}
	if data,err := ioutil.ReadFile(path) ; err == nil {
		json.Unmarshal(data,c)
	}else{
		return nil,err
	}
	return c,nil
}

func (mc MonitoringConfig)GetHeartbeats()[]heartbeat.Heartbeat {
	heartbeats := make([]heartbeat.Heartbeat,0,len(mc.HeartBeats))
	for _,hb := range mc.HeartBeats {
		if frequency, err := parseFrequency(hb.Frequency) ; err == nil {
			heartbeats = append(heartbeats, heartbeat.NewHeartBeat(hb.Name, hb.Url,frequency))
		}else{
			log.Println("Error",err)
		}
	}
	return heartbeats
}