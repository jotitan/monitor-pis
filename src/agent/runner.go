package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jotitan/monitor-pis/config"
	"github.com/jotitan/monitor-pis/model"
	"log"
	"net/http"
	"time"
)

type AgentRunner struct {
	duration time.Duration
	monitorUrl string
	metrics []Metric
	name string
}

func NewAgentRunner(config config.AgentConfig)(*AgentRunner,error){
	// Check if monitoring server is up
	if !checkUrl(config.MonitorUrl) {
		return nil,errors.New("monitoring service is not available : " + config.MonitorUrl)
	}
	frequency,err := config.GetFrequency()
	if err != nil {
		return nil,err
	}
	metrics,err := buildMetrics(config.Metrics)
	if err != nil {
		return nil,err
	}
	return &AgentRunner{duration:frequency,monitorUrl: config.MonitorUrl,metrics:metrics,name:config.Name},nil
}

func buildMetrics(names []string)([]Metric,error){
	metrics := make([]Metric,0,len(names))
	for _,name := range names {
		metric,err := getMetric(name)
		if err != nil {
			return nil,err
		}
		metrics = append(metrics,metric)
	}
	return metrics, nil
}

func getMetric(name string)(Metric,error){
	switch name {
	case "temperature":return NewTemperatureMetric(),nil
	case "memory":return NewMemoryMetric(),nil
	case "disk":return NewDiskMetric(),nil
	case "cpu":return NewCpuMetric(),nil
	default:return nil,errors.New("impossible to find metric called " + name)
	}
}

func checkUrl(url string)bool{
	resp,err := http.Get(url + "/health")
	return err == nil && resp.StatusCode == 200
}

func (ar AgentRunner)Run(){
	log.Println("Start agent with",len(ar.metrics),"metric(s)")
	for {
		<- time.NewTicker(ar.duration).C
		ar.sendMetrics(ar.computeMetrics())
	}
	log.Println("Error, agent stop")
}

func (ar AgentRunner)sendMetrics(metrics map[string]model.MetricPoint){
	request := model.MetricsRequest{Metrics: metrics,Name:ar.name}
	data,_ := json.Marshal(request)
	buffer := bytes.NewBuffer(data)
	http.Post(ar.monitorUrl + "/metric","application/json",buffer)
}

func (ar AgentRunner)computeMetrics()map[string]model.MetricPoint {
	results := make(map[string]model.MetricPoint)
	for _,metric := range ar.metrics {
		if value, name, err := metric.GetValue() ; err == nil {
			results[name] = model.MetricPoint{Value: value,Timestamp: time.Now().Unix()}
		}
	}
	return results
}

