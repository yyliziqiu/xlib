package xutil

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
)

func CompressBytesByGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	wtr := gzip.NewWriter(&buf)
	defer wtr.Close()

	_, err := wtr.Write(data)
	if err != nil {
		return nil, err
	}

	// 在返回结果之前，必须先调用 close()，
	// 否则 writer 缓冲区中的内容不会全部返回
	// 导致解压时报 io.ErrUnexpectedEOF
	err = wtr.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func UnCompressBytesByGzip(data []byte) ([]byte, error) {
	rdr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer rdr.Close()

	result, err := io.ReadAll(rdr)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func CompressBytesByZlib(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	wtr := zlib.NewWriter(&buf)
	defer wtr.Close()

	_, err := wtr.Write(data)
	if err != nil {
		return nil, err
	}

	// 在返回结果之前，必须先调用 close()，
	// 否则 writer 缓冲区中的内容不会全部返回
	// 导致解压时报 io.ErrUnexpectedEOF
	err = wtr.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func UnCompressBytesByZlib(data []byte) ([]byte, error) {
	rdr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer rdr.Close()

	result, err := io.ReadAll(rdr)
	if err != nil {
		return nil, err
	}

	return result, nil
}
