package monitor_server

import (
	"encoding/json"
	"github.com/jotitan/monitor-pis/config"
	"github.com/jotitan/monitor-pis/model"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type MonitoringServer struct{
	port string
	repository *MetricRepository
}

func NewMonitoringServer(conf config.MonitoringConfig)MonitoringServer{
	return MonitoringServer{conf.Port,NewMetricRepository(conf)}
}

func (ms MonitoringServer)Start(){
	log.Println("Start on",ms.port)
	s := ms.createServer()
	http.ListenAndServe(":" + ms.port,s)
}

func (ms MonitoringServer)createServer()*http.ServeMux{
	s := http.NewServeMux()
	s.HandleFunc("/health",ms.health)
	s.HandleFunc("/metric_names", ms.metricNames)
	s.HandleFunc("/metric", ms.metric)
	s.HandleFunc("/flush", ms.flush)
	s.HandleFunc("/search", ms.search)
	s.HandleFunc("/instances_names", ms.instancesNames)
	s.HandleFunc("/instances", ms.instances)
	s.HandleFunc("/heartbeats", ms.heartbeats)

	return s
}

func (ms MonitoringServer)health(w http.ResponseWriter, r * http.Request){
	addCors(w)
	w.WriteHeader(200)
}

func (ms MonitoringServer)instancesNames(w http.ResponseWriter, r * http.Request){
	addCors(w)
	data,_ := json.Marshal(ms.repository.GetInstancesName())
	w.Header().Set("Content-type","application/json")
	w.Write(data)
}

func (ms MonitoringServer)instances(w http.ResponseWriter, r * http.Request){
	addCors(w)
	data,_ := json.Marshal(ms.repository.GetLastMetrics())
	w.Header().Set("Content-type","application/json")
	w.Write(data)
}

func (ms MonitoringServer)heartbeats(w http.ResponseWriter, r * http.Request){
	addCors(w)
	heartbeats := ms.repository.GetLastHeartbeats()
	results := make([]interface{},0,len(heartbeats))
	for name, value := range heartbeats {
		results = append(results,struct{
			Name string
			IsUp bool
		}{name,value==1})
	}
	data,_ := json.Marshal(results)
	w.Header().Set("Content-type","application/json")
	w.Write(data)
}

func (ms MonitoringServer)search(w http.ResponseWriter, r * http.Request){
	addCors(w)
	instance := r.FormValue("instance")
	metric := r.FormValue("metric")
	if strings.EqualFold("",metric){
		ms.searchAll(w,instance)
	}else {
		points := ms.repository.Search(instance, metric)
		if data, err := json.Marshal(points); err == nil {
			log.Println("Search", instance, metric, ". Found", len(points))
			w.Header().Set("Content-type", "application/json")
			w.Write(data)
		}
	}
}

func (ms MonitoringServer)searchAll(w http.ResponseWriter, instanceName string){
	instance := ms.repository.getInstance(instanceName, false)
	metrics := make(map[string][]model.MetricPoint)
	for _,metric := range instance.getMetricsName(){
		metrics[metric] = ms.repository.Search(instanceName, metric)
	}
	log.Println("Search all metrics", instanceName)
	if data, err := json.Marshal(metrics); err == nil {
		w.Header().Set("Content-type", "application/json")
		w.Write(data)
	}
}

func (ms MonitoringServer)metricNames(w http.ResponseWriter, r * http.Request){
	addCors(w)
	instance := r.FormValue("instance")
	names := ms.repository.GetMetricsName(instance)
	if data,err := json.Marshal(names) ; err == nil {
		log.Println("Load metrics",instance)
		w.Header().Set("Content-type","application/json")
		w.Write(data)
	}
}

func (ms MonitoringServer)flush(w http.ResponseWriter, r * http.Request){
	addCors(w)
	if strings.EqualFold(r.Method,http.MethodPost) {
		ms.repository.Flush()
	}
}

func (ms MonitoringServer)metric(w http.ResponseWriter, r * http.Request){
	addCors(w)
	if strings.EqualFold(r.Method,http.MethodPost) {
		ms.pushMetric(w,r)
	}
}

func (ms MonitoringServer) pushMetric(w http.ResponseWriter, r * http.Request){
	addCors(w)
	data,_ := ioutil.ReadAll(r.Body)
	metrics := model.MetricsRequest{}
	json.Unmarshal(data,&metrics)
	ms.repository.AppendManyMetrics(metrics.Name,metrics.Metrics)
}

func addCors(w http.ResponseWriter){
	w.Header().Set("Access-Control-Allow-Origin","*")
}