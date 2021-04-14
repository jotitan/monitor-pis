package monitor_server

import (
	"encoding/json"
	"github.com/jotitan/monitor-pis/config"
	"github.com/jotitan/monitor-pis/model"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type MonitoringServer struct {
	port       string
	repository *MetricRepository
	resources  string
}

func NewMonitoringServer(conf config.MonitoringConfig)MonitoringServer{
	return MonitoringServer{port:conf.Port,repository: NewMetricRepository(conf),resources: conf.Resources}
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
	s.HandleFunc("/", ms.defaultHandle)

	return s
}

func (ms MonitoringServer)health(w http.ResponseWriter, r * http.Request){
	setCommonHeader(w)
	w.WriteHeader(200)
}

func (ms MonitoringServer)instancesNames(w http.ResponseWriter, r * http.Request){
	setCommonHeader(w)
	data,_ := json.Marshal(ms.repository.GetInstancesName())
	w.Write(data)
}

func (ms MonitoringServer)instances(w http.ResponseWriter, r * http.Request){
	setCommonHeader(w)
	data,_ := json.Marshal(ms.repository.GetLastMetrics())
	w.Write(data)
}

func (ms MonitoringServer)heartbeats(w http.ResponseWriter, r * http.Request){
	setCommonHeader(w)
	heartbeats := ms.repository.GetLastHeartbeats()
	results := make([]interface{},0,len(heartbeats))
	for name, value := range heartbeats {
		results = append(results,struct{
			Name string
			IsUp bool
		}{name,value==1})
	}
	data,_ := json.Marshal(results)
	w.Write(data)
}

func (ms MonitoringServer)defaultHandle(w http.ResponseWriter,r * http.Request){
	http.ServeFile(w, r, filepath.Join(ms.resources, r.RequestURI[1:]))
}

func extractSearchParam(r * http.Request)(string,string,string){
	return r.FormValue("instance"),
		r.FormValue("metric"),
		r.FormValue("date")
}

func (ms MonitoringServer)search(w http.ResponseWriter, r * http.Request){
	setCommonHeader(w)
	instance,metric,date := extractSearchParam(r)
	if strings.EqualFold("",metric){
		ms.searchAll(w,instance,date)
	}else {
		points := ms.repository.Search(instance, metric,date)
		if data, err := json.Marshal(points); err == nil {
			log.Println("Search", instance, metric, ". Found", len(points))
			w.Write(data)
		}
	}
}

func (ms MonitoringServer)searchAll(w http.ResponseWriter, instanceName,date string){
	instance := ms.repository.getInstance(instanceName, false)
	setCommonHeader(w)
	metrics := make(map[string][]model.MetricPoint)
	if instance != nil {
		for _,metric := range instance.getMetricsName(){
			metrics[metric] = ms.repository.Search(instanceName, metric,date)
		}
	}
	log.Println("Search all metrics", instanceName)
	if data, err := json.Marshal(metrics); err == nil {
		w.Write(data)
	}
}

func (ms MonitoringServer)metricNames(w http.ResponseWriter, r * http.Request){
	setCommonHeader(w)
	instance := r.FormValue("instance")
	names := ms.repository.GetMetricsName(instance)
	if data,err := json.Marshal(names) ; err == nil {
		log.Println("Load metrics",instance)
		w.Write(data)
	}
}

func (ms MonitoringServer)flush(w http.ResponseWriter, r * http.Request){
	setCommonHeader(w)
	if strings.EqualFold(r.Method,http.MethodPost) {
		ms.repository.Flush()
	}
}

func (ms MonitoringServer)metric(w http.ResponseWriter, r * http.Request){
	setCommonHeader(w)
	if strings.EqualFold(r.Method,http.MethodPost) {
		ms.pushMetric(w,r)
	}
}

func (ms MonitoringServer) pushMetric(w http.ResponseWriter, r * http.Request){
	setCommonHeader(w)
	data,_ := ioutil.ReadAll(r.Body)
	metrics := model.MetricsRequest{}
	json.Unmarshal(data,&metrics)
	ms.repository.AppendManyMetrics(metrics.Name,metrics.Metrics)
}

func setCommonHeader(w http.ResponseWriter){
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Content-type", "application/json")
}
