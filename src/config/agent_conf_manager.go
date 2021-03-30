package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"time"
)

type AgentConfig struct {
	// url to send metrics
	MonitorUrl string `json:"monitor_url"`
	// can be temperature, cpu, memory
	Metrics []string `json:"metrics"`
	Frequency string `json:"frequency"`
	Name string `json:"instance_name"`
}

func NewConfig(path string) AgentConfig {
	c := AgentConfig{}
	if data,err := ioutil.ReadFile(path) ; err == nil {
		json.Unmarshal(data,&c)
	}else{
		log.Println("Error",err)
	}
	return c
}

func (ac AgentConfig)GetFrequency()(time.Duration,error){
	return parseFrequency(ac.Frequency)
}

func parseFrequency(frequency string)(time.Duration, error) {
	r,_ := regexp.Compile("([0-9]+) ?([smh])")
	results := r.FindAllStringSubmatch(frequency,1)
	if len(results) == 1 {
		if value, err := strconv.ParseInt(results[0][1],10,32) ; err == nil {
			switch results[0][2]{
			case "s":return time.Duration(value)*time.Second,nil
			case "m":return time.Duration(value)*time.Minute,nil
			case "h":return time.Duration(value)*time.Hour,nil
			}
		}
	}
	return 0,errors.New("no duration define")
}