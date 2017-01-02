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