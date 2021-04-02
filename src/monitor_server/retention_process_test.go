package monitor_server

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRetention(t *testing.T){
	folder,_ := ioutil.TempDir("","retention_folder")

	currentDate,_ := time.Parse("20060102","20200918")
	oneWeekDate,_ := time.Parse("20060102","20200911")
	twoWeekDate,_ := time.Parse("20060102","20200904")
	threeWeekDate,_ := time.Parse("20060102","20200828")
	oneMonthDate,_ := time.Parse("20060102","20200818")

	createFile(folder,"instance1",oneWeekDate)
	createFile(folder,"instance1",twoWeekDate)
	createFile(folder,"instance2",twoWeekDate)
	createFile(folder,"instance1",threeWeekDate)
	createFile(folder,"instance3",threeWeekDate)
	createFile(folder,"instance1",oneMonthDate)
	createFile(folder,"instance2",oneMonthDate)
	createFile(folder,"instance3",oneMonthDate)
	createFile(folder,"instance4",oneMonthDate)

	rp := NewRetentionProcess(17,folder)
	rp.clean(currentDate)

	if files := listFiles(folder) ; len(files) != 3 {
		t.Error(fmt.Sprintf("Must found only 3 files but found %d files",len(files)))
	}
}

func listFiles(folder string)[]string {
	if dir,err := os.Open(folder) ; err == nil {
		files,_ := dir.Readdirnames(-1)
		results := make([]string,len(files))
		for i,file := range files {
			results[i] = filepath.Join(folder,file)
		}
		return results
	}
	return []string{}
}

func createFile(folder,instance string, date time.Time)(string,error){
	filename := filepath.Join(folder,fmt.Sprintf("metric_%s_%s.met",instance,date.Format("20060102")))
	if f,err := os.Create(filename); err != nil {
		return "",err
	}else{
		return filename,f.Close()
	}
}
