package monitor_server

import (
	"encoding/json"
	"github.com/jotitan/monitor-pis/config"
	"github.com/jotitan/monitor-pis/model"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	heartbeatInstanceName = "heartbeat"
)

// Manage all metrics store
type MetricRepository struct {
	folder    string
	instances map[string]*MetricInstanceRepository
	autoFlushLimit int
}

func NewMetricRepository(conf config.MonitoringConfig)*MetricRepository{
	// Load instances from file
	mr := &MetricRepository{folder:conf.Folder,instances:make(map[string]*MetricInstanceRepository),autoFlushLimit: conf.AutoFlushLimit}
	mr.loadInstancesNames()
	mr.launchHeartBeats(conf)

	return mr
}

func (mr * MetricRepository)launchHeartBeats(conf config.MonitoringConfig){
	for _,heartbeat :=range conf.GetHeartbeats() {
		log.Println("Run",heartbeat.Name)
		heartbeat.Start(mr.getInstance(heartbeatInstanceName,true).Append)
	}
}

func (mr * MetricRepository)loadInstancesNames(){
	if data, err := ioutil.ReadFile(mr.getInstancesFilename()) ; err == nil {
		instances := make([]string,0)
		json.Unmarshal(data,&instances)
		for _,instance := range instances {
			mr.instances[instance] = NewMetricInstanceRepository(mr.folder,instance,nbPointsByBlock,mr.autoFlushLimit)
		}
		log.Println("Load instances",len(mr.instances))
	}
}

func (mr * MetricRepository)getInstancesFilename()string{
	return filepath.Join(mr.folder,"instances.json")
}

func (mr * MetricRepository)saveInstancesNames(){
	if f, err := os.OpenFile(mr.getInstancesFilename(),os.O_CREATE|os.O_RDWR|os.O_TRUNC,os.ModePerm); err == nil {
		defer f.Close()
		data,_ := json.Marshal(mr.GetInstancesName())
		f.Write(data)
	}
}

func (mr * MetricRepository)getInstance(instanceName string, createIfNoExist bool)*MetricInstanceRepository{
	instance,exist := mr.instances[instanceName]
	if !exist && createIfNoExist {
		instance = NewMetricInstanceRepository(mr.folder,instanceName,nbPointsByBlock,mr.autoFlushLimit)
		mr.instances[instanceName] = instance
		mr.saveInstancesNames()
		// Save metrics
	}
	return instance
}

func (mr *MetricRepository) AppendMetrics(instanceName string, metricName string, points []model.MetricPoint) {
	instance := mr.getInstance(instanceName,true)
	instance.Append(metricName,points)
}

func (mr *MetricRepository) AppendManyMetrics(instanceName string, points map[string]model.MetricPoint) {
	instance := mr.getInstance(instanceName,true)
	log.Println("Push metric",instanceName)
	for metricName,value := range points {
		instance.Append(metricName,[]model.MetricPoint{value})
	}
}

func (mr *MetricRepository)Search(instanceName, metricName,date string)[]model.MetricPoint{
	if instance := mr.getInstance(instanceName, false) ; instance != nil {
		return instance.Search(metricName,date)
	}
	return []model.MetricPoint{}
}

func (mr MetricRepository) GetMetricsName(instanceName string)[]string {
	instance,exist := mr.instances[instanceName]
	if !exist {
		return []string{}
	}
	return instance.getMetricsName()
}

func (mr MetricRepository) GetLastHeartbeats()map[string]float32 {
	instance := mr.getInstance(heartbeatInstanceName,false)
	return instance.readLastMetrics()
}

func (mr MetricRepository) GetLastMetrics()map[string]map[string]float32 {
	instancesValues := make(map[string]map[string]float32)
	for _,instanceName := range mr.GetInstancesName() {
		if !strings.EqualFold(instanceName, heartbeatInstanceName) {
			instance := mr.getInstance(instanceName,false)
			instancesValues[instanceName] = instance.readLastMetrics()
		}
	}
	return instancesValues
}

func (mr MetricRepository) GetInstancesName()[]string {
	names := make([]string,0,len(mr.instances))
	for name := range mr.instances {
		names = append(names,name)
	}
	return names
}

func (mr *MetricRepository) Flush() {
	log.Println("Flush all instances")
	for _,instance := range mr.instances{
		instance.Flush()
	}
}
