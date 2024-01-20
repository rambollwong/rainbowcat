package util

import (
	"bytes"
	"encoding/binary"
)

// Uint64ToBytes converts a uint64 value to a byte slice.
func Uint64ToBytes(n uint64) []byte {
	x := n
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// BytesToUint64 converts a byte slice to a uint64 value.
func BytesToUint64(b []byte) uint64 {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint64
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

// Uint32ToBytes converts a uint32 value to a byte slice.
func Uint32ToBytes(n uint32) []byte {
	x := n
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// BytesToUint32 converts a byte slice to a uint32 value.
func BytesToUint32(b []byte) uint32 {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint32
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

// UintToBytes converts a uint value to a byte slice.
func UintToBytes(n uint) []byte {
	x := uint64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// BytesToUint converts a byte slice to a uint value.
func BytesToUint(b []byte) uint {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint64
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
	return uint(x)
}

// Int64ToBytes converts an int64 value to a byte slice.
func Int64ToBytes(n int64) []byte {
	x := n
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// BytesToInt64 converts a byte slice to an int64 value.
func BytesToInt64(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int64
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

// Int32ToBytes converts an int32 value to a byte slice.
func Int32ToBytes(n int32) []byte {
	x := n
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// BytesToInt32 converts a byte slice to an int32 value.
func BytesToInt32(b []byte) int32 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

// IntToBytes converts an int value to a byte slice.
func IntToBytes(n int) []byte {
	x := int64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// BytesToInt converts a byte slice to an int value.
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int64
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}
