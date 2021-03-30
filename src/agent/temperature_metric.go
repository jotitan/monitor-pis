package agent

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
)

type temperatureMetric struct{}

func NewTemperatureMetric()Metric{
	return temperatureMetric{}
}

func (tm temperatureMetric)GetValue()(float32,string,error){
	cmd := exec.Command("/opt/vc/bin/vcgencmd","measure_temp")
	data,err := cmd.Output()
	if err != nil {
		return 0,"",err
	}
	r,_ := regexp.Compile("temp=([0-9]+\\.[0-9]+)")
	results := r.FindAllStringSubmatch(string(data),1)
	if len(results) != 1 {
		return 0,"",errors.New("no results")
	}
	if value,err := strconv.ParseFloat(results[0][1],32);err == nil {
		return float32(value),"temperature",nil
	}else{
		return 0,"",err
	}
}
