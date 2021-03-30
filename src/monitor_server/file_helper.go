package monitor_server

import (
	"encoding/binary"
	"math"
	"os"
)

func readData(f *os.File,position int64, length int)[]byte{
	data := make([]byte,length)
	f.ReadAt(data,position)
	return data
}

func getBytesAsInt64(tab []byte)int64{
	return int64(binary.LittleEndian.Uint64(tab))
}

func getBytesAsInt32(tab []byte)int32{
	return int32(binary.LittleEndian.Uint32(tab))
}

func getFloat32AsBytes(value float32)[]byte{
	return getInt32AsBytes(int32(math.Float32bits(value)))
}

func getBytesAsFloat32(tab []byte)float32{
	return math.Float32frombits(uint32(getBytesAsInt32(tab)))
}

func getInt64AsBytes(value int64)[]byte{
	return []byte{byte(value),byte(value >> 8),byte(value >> 16),byte(value >> 24),
		byte(value >> 32),byte(value >> 40),byte(value >> 48),byte(value >> 56)}
}

func getInt32AsBytes(value int32)[]byte{
	return []byte{byte(value),byte(value >> 8),byte(value >> 16),byte(value >> 24)}
}

func writeStringToBytesWithPad(data []byte, value string, position, totalLength int){
	writeBytesToBytes([]byte(value),data,position)
	if totalLength > len(value) {
		writeBytesToBytes(make([]byte, totalLength - len(value)),data,position+len(value))
	}
}

func writeBytesToBytes(source,target []byte, position int){
	for pos,b := range source{
		target[position+pos] = b
	}
}
