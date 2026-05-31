package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
)

// Compressor compresses data.
type Compressor func([]byte) ([]byte, error)

// GzipCompress uses gzip (better ratio, more CPU).
var GzipCompress Compressor = func(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// FlateCompress uses raw DEFLATE (lighter).
var FlateCompress Compressor = func(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, flate.BestSpeed)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// NoCompress returns the data unchanged (for tiny payloads).
var NoCompress Compressor = func(data []byte) ([]byte, error) {
	return data, nil
}

// SelectStrategy picks a compressor based on payload size.
func SelectStrategy(size int) Compressor {
	switch {
	case size < 64:
		return NoCompress
	case size < 4096:
		return FlateCompress
	default:
		return GzipCompress
	}
}

// gunzip is a helper used by tests to verify round-trips.
func gunzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
