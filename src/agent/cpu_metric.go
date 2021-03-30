package agent

import (
	"os/exec"
	"strconv"
)

type cpuMetric struct{}

func NewCpuMetric()Metric{
	return cpuMetric{}
}

func (tm cpuMetric)GetValue()(float32,string,error){
	c := "mpstat | grep all | awk '{print $12}'"
	cmd := exec.Command("bash","-c",c)
	data,err := cmd.Output()

	if err != nil {
		return 0,"",err
	}
	value,err := strconv.ParseFloat(string(data),32)
	if err == nil {
		return 100 - float32(value),"cpu",nil
	}
	return 0,"",err
}
