package monitor_server

import (
	"fmt"
	"github.com/jotitan/monitor-pis/config"
	"github.com/jotitan/monitor-pis/model"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var defaultTime,_ = time.Parse("2006-01-02","2018-07-26")

func TestFilename(t *testing.T){
	r := NewMetricInstanceRepository("myfolder","instance1",1440,10)
	if value := r.getFilename(defaultTime) ; !strings.EqualFold(value,filepath.Join("myfolder","metric_instance1_20180726.met")){
		t.Error("Impossible filename " + value)
	}
}

func TestWriteFile(t *testing.T){
	dir,_ := ioutil.TempDir("","test")

	r := NewMetricInstanceRepository(dir,"instance1",5,5)
	r.Append("metric1",[]model.MetricPoint{
		model.NewMetricPoint(10,1.5),
		model.NewMetricPoint(11,1.3),
		model.NewMetricPoint(12,1.8),
		model.NewMetricPoint(14,2.1),
		model.NewMetricPoint(15,1.2),
	})
	r.Append("metric2",[]model.MetricPoint{
		model.NewMetricPoint(10,80),
		model.NewMetricPoint(15,74),
		model.NewMetricPoint(20,67),
	})
	r.Flush()

	// Update metrics of file
	r.Append("metric1",[]model.MetricPoint{
		model.NewMetricPoint(17,3.5),
		model.NewMetricPoint(18,0.5),
		model.NewMetricPoint(20,3.3),
	})
	r.Append("metric3",[]model.MetricPoint{
		model.NewMetricPoint(17,65),
		model.NewMetricPoint(18,66),
		model.NewMetricPoint(20,67),
		model.NewMetricPoint(21,71),
		model.NewMetricPoint(24,76),
	})

	r.Flush()

	// Read file to search inside
	points := r.Search("metric1","")
	if nb := len(points) ; nb != 8 {
		t.Error(fmt.Sprintf("Must find 8 values but find %d",nb))
	}

	r.Append("metric1",[]model.MetricPoint{
		model.NewMetricPoint(24,1.35),
		model.NewMetricPoint(26,1.5),
		model.NewMetricPoint(27,1.3),
		model.NewMetricPoint(27,-1.1),
	})
	points = r.Search("metric1","")
	if nb := len(points) ; nb != 12 {
		t.Error(fmt.Sprintf("Must find 12 values but find %d",nb))
	}
	if points[11].Value != -1.1 {
		t.Error("Bad value")
	}

	r2 := NewMetricInstanceRepository(dir,"instance1",5,10)
	points = r2.Search("metric1","")
	if nb := len(points) ; nb != 8 {
		t.Error(fmt.Sprintf("Must find 8 values but find %d",nb))
	}

	if values := r2.readLastMetrics() ; values["metric1"] != 3.3{
		t.Error("Bad value, must found -1.1",values)
	}
	r2.Append("metric1",[]model.MetricPoint{
		model.NewMetricPoint(30,14.5),
	})
	if values := r2.readLastMetrics() ; values["metric1"] != 14.5{
		t.Error("Bad value, must found 14.5",values)
	}

}


func TestSaveMetrics(t *testing.T){
	dir,_ := ioutil.TempDir("save_metric","test")
	mr := NewMetricRepository(config.MonitoringConfig{Folder:dir})
	metrics := []model.MetricPoint{
		model.NewMetricPoint(1000,12),
		model.NewMetricPoint(1001,45),
	}
	mr.AppendMetrics("instance1", "name",metrics)
	if names := mr.GetMetricsName("instance1") ; len(names) != 1 {
		t.Error("Must found 1 metric name")
	}
}

func TestHeader(t *testing.T){
	head := newHeader("my series")
	head.createMetric("temperature",0)
	head.createMetric("memory",0)
	head.createMetric("cpu",0)
	head.createMetric("heartbeat",0)

	head.updateMetric("temperature",1000,1000)
	head.updateMetric("cpu",5000,10000)

	data := head.toBytes()

	deserHead := headerFileMetrics{}
	deserHead.fromBytes(data)
	if !strings.EqualFold(head.name,deserHead.name) {
		t.Error("Name must be same but find " + deserHead.name)
	}
	if len(head.headerMetrics) != len(deserHead.headerMetrics){
		t.Error(fmt.Sprintf("Must find same columns number but found %d",len(deserHead.headerMetrics)))
	}
	if value := deserHead.getMetric("temperature").firstBlockPosition; value != 1000 {
		t.Error(fmt.Sprintf("Temperature must have first bloc with value 1000 but found %d",value))
	}
	if value := deserHead.getMetric("cpu").currentBlockPosition; value != 10000 {
		t.Error(fmt.Sprintf("cpu must have current bloc with value 10000 but found %d",value))
	}
	if value := deserHead.getMetric("heartbeat").currentBlockPosition; value != 0 {
		t.Error(fmt.Sprintf("heartbeat must have current bloc with value 0 but found %d",value))
	}
	if deserHead.getMetric("toto") != nil {
		t.Error("toto must be null")
	}
}