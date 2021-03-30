package agent

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
)

type diskMetric struct{}

func NewDiskMetric()Metric{
	return diskMetric{}
}

func (tm diskMetric)GetValue()(float32,string,error){
	c := "df |grep \"/dev/root\" | awk {'print $5'}"
	//cmd := exec.Command("df -ah","| grep \"/dev/root\"")

	cmd := exec.Command("bash","-c",c)
	data,err := cmd.Output()

	if err != nil {
		return 0,"",err
	}
	r,_ := regexp.Compile("([0-9]+)")
	results := r.FindAllStringSubmatch(string(data),1)
	if len(results) != 1 {
		return 0,"",errors.New("no results")
	}
	value,err := strconv.ParseInt(results[0][1],10,32)
	if err == nil {
		return float32(value),"disk",nil
	}
	return 0,"",err
}
