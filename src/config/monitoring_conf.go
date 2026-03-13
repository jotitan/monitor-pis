package config

import (
	"encoding/json"
	"os"
)

type HeartbeatConfig struct {
	Name      string      `json:"name"`
	Url       string      `json:"url"`
	Frequency string      `json:"frequency"`
	Alert     AlertConfig `json:"alert"`
}

type AlertConfig struct {
	Threashold  int    `json:"threashold"`
	Email       string `json:"email"`
	AlertWhenUp bool   `json:"alert_when_up"`
}

type EmailConfig struct {
	HostSMTP     string `json:"host"`
	LoginSMTP    string `json:"login"`
	PasswordSMTP string `json:"password"`
	PortSMTP     string `json:"port"`
	EmailSender  string `json:"sender"`
}

type MonitoringConfig struct {
	// url to send metrics
	Port              string            `json:"port"`
	EmailSenderConfig EmailConfig       `json:"email"`
	Folder            string            `json:"folder"`
	AutoFlushLimit    int               `json:"auto_flush,omitempty"`
	RetentionDays     int               `json:"retention_days,omitempty"`
	HeartBeats        []HeartbeatConfig `json:"heartbeats"`
	Resources         string            `json:"resources"`
}

func NewMonitoringConfig(path string) (*MonitoringConfig, error) {
	c := &MonitoringConfig{AutoFlushLimit: 5, RetentionDays: 30}
	if data, err := os.ReadFile(path); err == nil {
		return c, json.Unmarshal(data, c)
	} else {
		return nil, err
	}
}
