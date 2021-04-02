package monitor_server

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// clean file after x days

type RetentationProcess struct {
	// After days, remove file
	days int
	folder string
}

func NewRetentionProcess(days int, folder string)*RetentationProcess{
	log.Println("Retention process :",days,"days")
	return &RetentationProcess{days,folder}
}

func (rp RetentationProcess)Launch(){
	// Launch one time per day
	ticker := time.NewTicker(24*time.Hour).C
	for {
		rp.cleanToday()
		<- ticker
	}
}

func (rp RetentationProcess) clean(currentDate time.Time){
	dateBefore := currentDate.Add(time.Duration(-24*rp.days)*time.Hour)
	if dir,err := os.Open(rp.folder) ; err == nil {
		defer dir.Close()
		files,_ := dir.Readdirnames(-1)
		for _,file := range files {
			if d,err := extractDate(file) ; err == nil {
				if dateBefore.Sub(d) > 0 {
					// Delete the file
					log.Println("Remove",file)
					os.Remove(filepath.Join(rp.folder,file))
				}
			}
		}
	}
}

func extractDate(name string)(time.Time,error){
	// Pattern "metric_instance_date.met
	r,_ := regexp.Compile("metric_(?:.*)_([0-9]{8}).met")
	if results := r.FindAllStringSubmatch(name,-1) ; len(results) == 1{
		if date,err := time.Parse("20060102",results[0][1]) ; err == nil {
			return date,nil
		}
	}
	return time.Now(),errors.New("impossible to parse")
}

func (rp RetentationProcess)cleanToday(){
	log.Println("Launch clean old files")
	rp.clean(time.Now())
}
