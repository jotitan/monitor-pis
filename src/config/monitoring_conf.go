package config

import (
	"encoding/json"
	"github.com/jotitan/monitor-pis/heartbeat"
	"io/ioutil"
	"log"
)

type HeartbeatConfig struct{
	Name      string `json:"name"`
	Url       string `json:"url"`
	Frequency string `json:"frequency"`
}

type MonitoringConfig struct {
	// url to send metrics
	Port           string `json:"port"`
	Folder         string `json:"folder"`
	AutoFlushLimit int    `json:"auto_flush,omitempty"`
	RetentionDays int    `json:"retention_days,omitempty"`
	HeartBeats     []HeartbeatConfig `json:"heartbeats"`
	Resources string `json:"resources"`
}

func NewMonitoringConfig(path string) (*MonitoringConfig,error) {
	c := &MonitoringConfig{AutoFlushLimit: 5,RetentionDays: 30}
	if data,err := ioutil.ReadFile(path) ; err == nil {
		return c,json.Unmarshal(data,c)
	}else{
		return nil,err
	}
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