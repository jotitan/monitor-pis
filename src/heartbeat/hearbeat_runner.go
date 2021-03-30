package heartbeat

import (
	"crypto/tls"
	"github.com/jotitan/monitor-pis/model"
	"log"
	"net/http"
	"strings"
	"time"
)

type Heartbeat struct {
	Name string
	formatName string
	url string
	frequency time.Duration
}

func NewHeartBeat(name, url string, frequency time.Duration)Heartbeat{
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	formatName := strings.ToLower(strings.ReplaceAll(name," ","_"))
	return Heartbeat{Name:name,formatName:formatName,url:url,frequency: frequency}
}

func (hb Heartbeat)Start(appendMethod func(string,[]model.MetricPoint)){
	ticker := time.NewTicker(hb.frequency).C
	go func(){
		for{
			<- ticker
			point := model.MetricPoint{Value:hb.testAsFloat(),Timestamp: time.Now().Unix()}
			appendMethod(hb.formatName,[]model.MetricPoint{point})
		}
	}()
}

func (hb Heartbeat)testAsFloat()float32{
	if hb.test() {
		return 1
	}
	return 0
}

func (hb Heartbeat)test()bool{
	log.Println("Call",hb.url)
	resp,err := http.Get(hb.url)
	defer func(){
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	return err == nil && resp.StatusCode == 200
}