package agent

import (
	"github.com/cloudfoundry/gosigar"
)

type memoryMetric struct{}

func NewMemoryMetric()Metric{
	return memoryMetric{}
}

func (tm memoryMetric)GetValue()(float32,string,error){
	m := sigar.Mem{}
	m.Get()
	return float32(m.Used*100/m.Total),"memory",nil
}
