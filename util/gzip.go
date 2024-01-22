package util

import (
	"bytes"
	"compress/gzip"
	"io"
)

// GZipCompressBytes compresses a byte slice using gzip compression.
// It returns the compressed byte slice and any error encountered during the compression process.
func GZipCompressBytes(data []byte) ([]byte, error) {
	var input bytes.Buffer
	g, err := gzip.NewWriterLevel(&input, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}
	_, err = g.Write(data)
	if err != nil {
		return nil, err
	}
	err = g.Close()
	if err != nil {
		return nil, err
	}
	return input.Bytes(), nil
}

// GZipDecompressBytes decompresses a byte slice using gzip decompression.
// It returns the decompressed byte slice and any error encountered during the decompression process.
func GZipDecompressBytes(data []byte) ([]byte, error) {
	var out bytes.Buffer
	var in bytes.Buffer
	in.Write(data)
	r, err := gzip.NewReader(&in)
	if err != nil {
		return nil, err
	}
	// nolint:gosec
	_, err = io.Copy(&out, r)
	if err != nil {
		return nil, err
	}
	err = r.Close()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
