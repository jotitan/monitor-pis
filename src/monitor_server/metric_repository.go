package monitor_server

import (
	"errors"
	"fmt"
	"github.com/jotitan/monitor-pis/model"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	nbPointsByBlock = 1440
)

// File structure headerFileMetrics - bloc metric 1440 points - pointer next block - bloc metric 1440 points...
// Header : instance_name(50)-metric_header(20 metrics):name(20)-pointer_start-current_pointer

type headerMetric struct {
	name string
	firstBlockPosition int64
	currentBlockPosition int64	// current
}

// Data//| block => nb value (4) | next block position(8)| [timestamp(8) | value (4)]
type headerFileMetrics struct {
	name string	// max 48 char
	headerMetrics []*headerMetric
}

func newHeader(name string)*headerFileMetrics {
	return &headerFileMetrics{name: name,headerMetrics:make([]*headerMetric,0,20)}
}

// structure size name(1) |name (49) | nbMetric(1) | 20[size metric name(1)|name(19)|first block (8) | current block (8)]
func (h *headerFileMetrics)fromBytes(data []byte){
	sizeName := int(data[0])
	h.name = string(data[1:sizeName+1])
	nbMetrics := int(data[50])
	position:=51
	h.headerMetrics = make([]*headerMetric,nbMetrics)
	for i := 0 ; i < nbMetrics ; i++ {
		sizeMetricName := int(data[position])
		position+=1
		h.headerMetrics[i] = &headerMetric{name: string(data[position:position+sizeMetricName])}
		position+=29
		h.headerMetrics[i].firstBlockPosition = getBytesAsInt64(data[position:position+8])
		h.headerMetrics[i].currentBlockPosition = getBytesAsInt64(data[position+8:position+16])
		position+=16
	}
}

func(h headerFileMetrics)sizeHeader()int{
	return 971
}

func (h headerFileMetrics)toBytes()[]byte {
	// Size headerFileMetrics : name(50) + 1 + 20 x (46) = 971
	data := make([]byte,h.sizeHeader())
	data[0] = byte(len(h.name))
	writeStringToBytesWithPad(data,h.name,1,49)
	data[50] = byte(len(h.headerMetrics))
	for i,hm := range h.headerMetrics {
		position := 51+i*46
		data[position] = byte(len(hm.name))
		writeStringToBytesWithPad(data,hm.name,position+1,29)
		writeBytesToBytes(getInt64AsBytes(hm.firstBlockPosition),data,position+30)
		writeBytesToBytes(getInt64AsBytes(hm.currentBlockPosition),data,position+38)
	}
	return data
}

func (h headerFileMetrics)hasMetric(name string)bool{
	return h.getMetric(name) != nil
}

func (h headerFileMetrics)getMetric(name string)*headerMetric{
	for _,hm := range h.headerMetrics {
		if strings.EqualFold(name,hm.name){
			return hm
		}
	}
	return nil
}

func (h *headerFileMetrics)createMetric(name string, position int64) (*headerMetric,error){
	if len(h.headerMetrics) >=20 {
		return nil,errors.New("impossible to add more than 20 metrics")
	}
	metric := &headerMetric{name:name,currentBlockPosition: position,firstBlockPosition: position}
	h.headerMetrics = append(h.headerMetrics,metric)
	return metric,nil
}

func (h *headerFileMetrics)updateMetric(name string, firstBlock, currentClock int64){
	met := h.getMetric(name)
	met.firstBlockPosition = firstBlock
	met.currentBlockPosition = currentClock
}

type blockHeader struct{
	nbPoints int32
	nextBlock int64
	positionInFile int64
	pointsByBlock int
}

func createNewBlockHeader(position int64, pointsByBlock int)*blockHeader{
	return &blockHeader{nbPoints:0,nextBlock:0,positionInFile: position,pointsByBlock:pointsByBlock}
}

func newBlockHeader(f *os.File,position int64,pointsByBlock int)*blockHeader{
	data := readData(f,position,12)
	return &blockHeader{nbPoints:getBytesAsInt32(data[0:4]),nextBlock:getBytesAsInt64(data[4:12]),positionInFile: position,pointsByBlock:pointsByBlock}
}

func (bh blockHeader)availableSpace()int{
	return bh.pointsByBlock - int(bh.nbPoints)
}

func (bh blockHeader)readPoints(f * os.File)[]model.MetricPoint{
	points := make([]model.MetricPoint,bh.nbPoints)
	data := readData(f,bh.positionInFile+12,int(bh.nbPoints*12))
	for i := 0 ; i < int(bh.nbPoints) ; i++ {
		points[i] = model.MetricPoint{
			Timestamp: getBytesAsInt64(data[i*12:i*12+8]),
			Value:getBytesAsFloat32(data[i*12+8:i*12+12]),
		}
	}
	return points
}

// Position : position in file + headerFileMetrics size + points
func (bh blockHeader)getPositionInBlock()int64{
	return bh.positionInFile + 12 + int64(bh.nbPoints)*12
}

func (bh blockHeader)getSizeBlock()int{
	return 12 + bh.pointsByBlock*12
}

func (bh *blockHeader)updateNbPoints(nbPoints int){
	bh.nbPoints+=int32(nbPoints)
}

func (bh *blockHeader) flushHeader(f *os.File){
	f.WriteAt(getInt32AsBytes(bh.nbPoints),bh.positionInFile)
	f.WriteAt(getInt64AsBytes(bh.nextBlock),bh.positionInFile+4)
}

// Store metrics of an instance
type MetricInstanceRepository struct {
	Name string
	Folder string
	// Store last values, other values are in file
	lastMetrics map[string][]model.MetricPoint
	// Limit to flush auto metrics in file
	AutoFlushLimit int
	PointsByBlock int
	// Store the headerFileMetrics
	head *headerFileMetrics
	mutex *sync.Mutex
	metricsName map[string]struct{}
}

func NewMetricInstanceRepository(folder,instanceName string,nbPointsByBlock, autoFlushLimit int)*MetricInstanceRepository{
	mir :=  &MetricInstanceRepository{
		Folder:folder,
		Name:instanceName,
		lastMetrics: make(map[string][]model.MetricPoint),
		PointsByBlock: nbPointsByBlock,
		AutoFlushLimit: autoFlushLimit,
		mutex:&sync.Mutex{},
		metricsName: make(map[string]struct{}),
	}

	return mir
}

func (r * MetricInstanceRepository)readMetricsNamesFromHeader(){
	filename := r.getFilename(time.Now())
	// Open file, load headerFileMetrics
	f,err := os.Open(filename)
	defer f.Close()
	if err == nil {
		r.head = r.readHeader(f)
		for _,m := range r.head.headerMetrics {
			r.metricsName[m.name] = struct{}{}
		}
	}
}

// Return metrics, mix header (already flush) and new in memory
func (r MetricInstanceRepository) getMetricsName() []string {
	set := make(map[string]struct{})
	for name := range r.metricsName {
		set[name] = struct{}{}
	}
	if r.head != nil {
		for _, hm := range r.head.headerMetrics {
			set[hm.name] = struct{}{}
		}
	}
	list := make([]string,0,len(set))
	for name := range set {
		list = append(list,name)
	}
	return list
}

func (r * MetricInstanceRepository)Flush()error{
	// Block read/write
	r.mutex.Lock()
	defer r.mutex.Unlock()
	filename := r.getFilename(time.Now())
	// Open file, load headerFileMetrics
	f,err := os.OpenFile(filename,os.O_CREATE|os.O_RDWR,os.ModePerm)
	defer f.Close()
	if err != nil {
		return err
	}
	// Go to end
	if position,_ := f.Seek(0,2) ; position == 0 {
		// New file, write headerFileMetrics
		r.head = newHeader(r.Name)
		// Reserve space
		f.WriteAt(make([]byte,r.head.sizeHeader()),0)
		r.flushPoints(f)
	}else{
		if r.head == nil {
			r.head = r.readHeader(f)
		}
		r.flushPoints(f)
	}
	r.lastMetrics = make(map[string][]model.MetricPoint)
	return nil
}

func (r MetricInstanceRepository)getPointsAsBytes(points []model.MetricPoint)[]byte{
	data := make([]byte,12*len(points))
	for i,point := range points {
		writeBytesToBytes(getInt64AsBytes(point.Timestamp),data,i*12)
		writeBytesToBytes(getFloat32AsBytes(point.Value),data,i*12+8)
	}
	return data
}

func (r * MetricInstanceRepository)flushPoints(f *os.File){
	for nameMetric,points := range r.lastMetrics {
		var hm *headerMetric
		var bh *blockHeader
		if r.head.hasMetric(nameMetric) {
			// compute position
			hm = r.head.getMetric(nameMetric)
			// get block to write, read
			bh = newBlockHeader(f, hm.currentBlockPosition,r.PointsByBlock)
		}else {
			hm,_ = r.head.createMetric(nameMetric,r.getEndPosition(f))
			bh = createNewBlockHeader(r.getEndPosition(f),r.PointsByBlock)
			// Reserve space in file
			f.WriteAt(make([]byte,bh.getSizeBlock()),bh.positionInFile)
		}
		r.writePointsInBlock(f,bh,points,hm)
	}
	f.WriteAt(r.head.toBytes(),0)
	f.Close()
}

func (r MetricInstanceRepository)writePointsInBlock(f *os.File,bh *blockHeader,points []model.MetricPoint, hm * headerMetric){
	size := bh.availableSpace()
	pointsToWrite := points
	if size < len(points){
		pointsToWrite = points[0:size]
	}
	// Write points
	f.WriteAt(r.getPointsAsBytes(pointsToWrite),bh.getPositionInBlock())
	bh.updateNbPoints(len(pointsToWrite))
	if size < len(points){
		// Write more
		nextBlock := createNewBlockHeader(r.getEndPosition(f),r.PointsByBlock)
		f.WriteAt(make([]byte,bh.getSizeBlock()),nextBlock.positionInFile)
		bh.nextBlock = nextBlock.positionInFile
		hm.currentBlockPosition = nextBlock.positionInFile
		r.writePointsInBlock(f,nextBlock,points[size:],hm)
	}
	bh.flushHeader(f)
}

func (r MetricInstanceRepository)getEndPosition(f *os.File)int64{
	position,_ := f.Seek(0,2)
	return position
}

func (r MetricInstanceRepository)readHeader(f *os.File)*headerFileMetrics {
	h := newHeader("")
	h.fromBytes(readData(f,0,h.sizeHeader()))
	return h
}

func (r MetricInstanceRepository)getFilename(date time.Time)string{
	return filepath.Join(r.Folder,fmt.Sprintf("metric_%s_%s.met",r.Name,date.Format("20060102")))
}

func (r * MetricInstanceRepository)Search(metricName,date string)[]model.MetricPoint {
	if points,exist := r.lastMetrics[metricName] ; exist {
		// Read in file and append current points
		return append(r.readMetricsFromFile(metricName,date),points...)
	}
	return r.readMetricsFromFile(metricName,date)
}

func (r *MetricInstanceRepository)readMetricsFromFile(metricName,date string)[]model.MetricPoint{
	readDate := time.Now()
	if !strings.EqualFold("",date){
		if d,err := time.Parse("2006-01-02",date) ; err == nil {
			readDate = d
		}
	}
	f,_ := os.Open(r.getFilename(readDate))
	defer f.Close()
	if r.head == nil {
		r.head = r.readHeader(f)
	}
	hm := r.head.getMetric(metricName)
	if hm == nil {
		return []model.MetricPoint{}
	}
	return r.readBlock(f,hm.firstBlockPosition)
}

func (r *MetricInstanceRepository) readLastMetrics()map[string]float32{
	f,_ := os.Open(r.getFilename(time.Now()))
	defer f.Close()
	if r.head == nil {
		r.head = r.readHeader(f)
	}
	lasts := make(map[string]float32)
	for _,name := range r.getMetricsName() {
		// Check if last value exist
		if values,exist := r.lastMetrics[name] ; exist && len(values) > 0{
			lasts[name] = values[len(values)-1].Value
		}else{
			// Read file
			hm := r.head.getMetric(name)
			lasts[name] = r.readLastValueBlock(f,hm.currentBlockPosition)
		}
	}
	return lasts
}

func (r MetricInstanceRepository)readBlock(f *os.File,position int64)[]model.MetricPoint{
	bh := newBlockHeader(f,position,r.PointsByBlock)
	points := bh.readPoints(f)
	if bh.nextBlock != 0 {
		points = append(points,r.readBlock(f,bh.nextBlock)...)
	}
	return points
}

func (r MetricInstanceRepository)readLastValueBlock(f *os.File,position int64)float32{
	bh := newBlockHeader(f,position,r.PointsByBlock)
	lastPosition := bh.positionInFile + int64(12*bh.nbPoints +8)
	data := readData(f,lastPosition,4)
	return getBytesAsFloat32(data)
}

func (r * MetricInstanceRepository)Append(metricName string,points []model.MetricPoint){
	metricsPoint,exist := r.lastMetrics[metricName]
	if !exist {
		metricsPoint = make([]model.MetricPoint,0)
		r.updateMetrics(metricName,metricsPoint)
		r.metricsName[metricName] = struct{}{}
	}
	r.updateMetrics(metricName,append(metricsPoint,points...))
	r.checkAutoFlush()
}

func (r * MetricInstanceRepository)updateMetrics(metricName string, points []model.MetricPoint){
	// Use a locker (same as file)
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.lastMetrics[metricName] = points
}

func (r * MetricInstanceRepository)checkAutoFlush(){
	count := 0
	for _,points := range r.lastMetrics {
		count+=len(points)
	}
	if count >= r.AutoFlushLimit {
		r.Flush()
	}
}
