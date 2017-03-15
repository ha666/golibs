package golibs

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
)

func ZlibZipBytes(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	compressor, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return nil, err
	}
	compressor.Write(input)
	compressor.Close()
	return buf.Bytes(), nil
}

func ZlibUnzipBytes(input []byte) ([]byte, error) {
	b := bytes.NewReader(input)
	r, err := zlib.NewReader(b)
	defer r.Close()
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(r)
	return data, nil
}

// 混淆[]byte
func ConfusedTwo(sourceBytes []byte) []byte {
	var confusedBytes []byte = make([]byte, len(sourceBytes))
	idx := 0
	for index := 0; index < len(sourceBytes); index++ {
		if index%2 == 0 {
			confusedBytes[idx] = byte(255 - sourceBytes[index])
			idx++
		}
	}
	for index := 0; index < len(sourceBytes); index++ {
		if index%2 == 1 {
			confusedBytes[idx] = byte(255 - sourceBytes[index])
			idx++
		}
	}
	return confusedBytes
}

// 反混淆[]byte
func UnConfusedTwo(confusedBytes []byte) []byte {
	loopCount := len(confusedBytes)
	count := loopCount / 2
	if loopCount%2 == 1 {
		count = loopCount/2 + 1
	}
	beforeBytes := confusedBytes[0:count]
	afterBytes := confusedBytes[count:]
	var unConfusedBytes []byte = make([]byte, loopCount)
	beforeIndex := 0
	afterIndex := 0
	for index := 0; index < loopCount; index++ {
		if index%2 == 0 {
			if beforeIndex >= len(beforeBytes) {
				continue
			}
			unConfusedBytes[index] = 255 - beforeBytes[beforeIndex]
			beforeIndex++
		} else {
			if afterIndex >= len(afterBytes) {
				continue
			}
			unConfusedBytes[index] = 255 - afterBytes[afterIndex]
			afterIndex++
		}
	}
	return unConfusedBytes
}
