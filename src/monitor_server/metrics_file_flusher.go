package monitor_server

import (
	"github.com/jotitan/monitor-pis/model"
	"os"
	"sync"
)

type FileFlusher struct {
	pointsByBlock int
}

func (ff FileFlusher) flushPoints(header *headerFileMetrics, locker *sync.RWMutex,f *os.File, metrics map[string][]model.MetricPoint) {
	locker.RLock()
	defer locker.RUnlock()
	for nameMetric,points := range metrics {
		var hm *headerMetric
		var bh *blockHeader
		if header.hasMetric(nameMetric) {
			// compute position
			hm = header.getMetric(nameMetric)
			// get block to write, read
			bh = newBlockHeader(f, hm.currentBlockPosition,ff.pointsByBlock)
		}else {
			hm,_ = header.createMetric(nameMetric,getEndPosition(f))
			bh = createNewBlockHeader(getEndPosition(f),ff.pointsByBlock)
			// Reserve space in file
			f.WriteAt(make([]byte,bh.getSizeBlock()),bh.positionInFile)
		}
		ff.writePointsInBlock(f,bh,points,hm)
	}
	f.WriteAt(header.toBytes(),0)
	f.Close()
}

func (ff FileFlusher)writePointsInBlock(f *os.File,bh *blockHeader,points []model.MetricPoint, hm * headerMetric){
	size := bh.availableSpace()
	pointsToWrite := points
	if size < len(points){
		pointsToWrite = points[0:size]
	}
	// Write points
	f.WriteAt(getPointsAsBytes(pointsToWrite),bh.getPositionInBlock())
	bh.updateNbPoints(len(pointsToWrite))
	if size < len(points){
		// Write more
		nextBlock := createNewBlockHeader(getEndPosition(f),ff.pointsByBlock)
		f.WriteAt(make([]byte,bh.getSizeBlock()),nextBlock.positionInFile)
		bh.nextBlock = nextBlock.positionInFile
		hm.currentBlockPosition = nextBlock.positionInFile
		ff.writePointsInBlock(f,nextBlock,points[size:],hm)
	}
	bh.flushHeader(f)
}

func getPointsAsBytes(points []model.MetricPoint)[]byte{
	data := make([]byte,12*len(points))
	for i,point := range points {
		writeBytesToBytes(getInt64AsBytes(point.Timestamp),data,i*12)
		writeBytesToBytes(getFloat32AsBytes(point.Value),data,i*12+8)
	}
	return data
}

func getEndPosition(f *os.File)int64{
	position,_ := f.Seek(0,2)
	return position
}
