package cdt

import (
	"time"

	"github.com/liasece/gocoder"
)

// Error func
func Error() gocoder.Type {
	return gocoder.NewTypeName("error")
}

// Context func
func Context() gocoder.Type {
	return gocoder.NewTypeDetail("context", "Context")
}

// I func
func I(i interface{}) gocoder.Type {
	return gocoder.NewTypeI(i)
}

// Time func
func Time() gocoder.Type {
	return gocoder.NewTypeI(time.Time{})
}

// TimePtr func
func TimePtr() gocoder.Type {
	v := time.Time{}
	return gocoder.NewTypeI(&v)
}

// String func
func String() gocoder.Type {
	return gocoder.NewTypeI(string(""))
}

// StringPtr func
func StringPtr() gocoder.Type {
	v := string("")
	return gocoder.NewTypeI(&v)
}

// StringSlice func
func StringSlice() gocoder.Type {
	return gocoder.NewTypeI([]string{})
}

// StringSlicePtr func
func StringSlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]string{})
}

// Bool func
func Bool() gocoder.Type {
	return gocoder.NewTypeI(bool(false))
}

// BoolPtr func
func BoolPtr() gocoder.Type {
	v := bool(false)
	return gocoder.NewTypeI(&v)
}

// BoolSlice func
func BoolSlice() gocoder.Type {
	return gocoder.NewTypeI([]bool{})
}

// BoolSlicePtr func
func BoolSlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]bool{})
}

// Int func
func Int() gocoder.Type {
	return gocoder.NewTypeI(int(0))
}

// IntPtr func
func IntPtr() gocoder.Type {
	v := int(0)
	return gocoder.NewTypeI(&v)
}

// IntSlice func
func IntSlice() gocoder.Type {
	return gocoder.NewTypeI([]int{})
}

// IntSlicePtr func
func IntSlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]int{})
}

// Int8 func
func Int8() gocoder.Type {
	return gocoder.NewTypeI(int8(0))
}

// Int8Ptr func
func Int8Ptr() gocoder.Type {
	v := int8(0)
	return gocoder.NewTypeI(&v)
}

// Int8Slice func
func Int8Slice() gocoder.Type {
	return gocoder.NewTypeI([]int8{})
}

// Int8SlicePtr func
func Int8SlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]int8{})
}

// Int16 func
func Int16() gocoder.Type {
	return gocoder.NewTypeI(int16(0))
}

// Int16Ptr func
func Int16Ptr() gocoder.Type {
	v := int16(0)
	return gocoder.NewTypeI(&v)
}

// Int16Slice func
func Int16Slice() gocoder.Type {
	return gocoder.NewTypeI([]int16{})
}

// Int16SlicePtr func
func Int16SlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]int16{})
}

// Int32 func
func Int32() gocoder.Type {
	return gocoder.NewTypeI(int32(0))
}

// Int32Ptr func
func Int32Ptr() gocoder.Type {
	v := int32(0)
	return gocoder.NewTypeI(&v)
}

// Int32Slice func
func Int32Slice() gocoder.Type {
	return gocoder.NewTypeI([]int32{})
}

// Int32SlicePtr func
func Int32SlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]int32{})
}

// Int64 func
func Int64() gocoder.Type {
	return gocoder.NewTypeI(int64(0))
}

// Int64Ptr func
func Int64Ptr() gocoder.Type {
	v := int64(0)
	return gocoder.NewTypeI(&v)
}

// Int64Slice func
func Int64Slice() gocoder.Type {
	return gocoder.NewTypeI([]int64{})
}

// Int64SlicePtr func
func Int64SlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]int64{})
}

// Uint func
func Uint() gocoder.Type {
	return gocoder.NewTypeI(uint(0))
}

// UintPtr func
func UintPtr() gocoder.Type {
	v := uint(0)
	return gocoder.NewTypeI(&v)
}

// UintSlice func
func UintSlice() gocoder.Type {
	return gocoder.NewTypeI([]uint{})
}

// UintSlicePtr func
func UintSlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]uint{})
}

// Uint8 func
func Uint8() gocoder.Type {
	return gocoder.NewTypeI(uint8(0))
}

// Uint8Ptr func
func Uint8Ptr() gocoder.Type {
	v := uint8(0)
	return gocoder.NewTypeI(&v)
}

// Uint8Slice func
func Uint8Slice() gocoder.Type {
	return gocoder.NewTypeI([]uint8{})
}

// Uint8SlicePtr func
func Uint8SlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]uint8{})
}

// Byte func
func Byte() gocoder.Type {
	return gocoder.NewTypeI(byte(0))
}

// BytePtr func
func BytePtr() gocoder.Type {
	v := byte(0)
	return gocoder.NewTypeI(&v)
}

// ByteSlice func
func ByteSlice() gocoder.Type {
	return gocoder.NewTypeI([]byte{})
}

// ByteSlicePtr func
func ByteSlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]byte{})
}

// Uint16 func
func Uint16() gocoder.Type {
	return gocoder.NewTypeI(uint16(0))
}

// Uint16Ptr func
func Uint16Ptr() gocoder.Type {
	v := uint16(0)
	return gocoder.NewTypeI(&v)
}

// Uint16Slice func
func Uint16Slice() gocoder.Type {
	return gocoder.NewTypeI([]uint16{})
}

// Uint16SlicePtr func
func Uint16SlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]uint16{})
}

// Uint32 func
func Uint32() gocoder.Type {
	return gocoder.NewTypeI(uint32(0))
}

// Uint32Ptr func
func Uint32Ptr() gocoder.Type {
	v := uint32(0)
	return gocoder.NewTypeI(&v)
}

// Uint32Slice func
func Uint32Slice() gocoder.Type {
	return gocoder.NewTypeI([]uint32{})
}

// Uint32SlicePtr func
func Uint32SlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]uint32{})
}

// Uint64 func
func Uint64() gocoder.Type {
	return gocoder.NewTypeI(uint64(0))
}

// Uint64Ptr func
func Uint64Ptr() gocoder.Type {
	v := uint64(0)
	return gocoder.NewTypeI(&v)
}

// Uint64Slice func
func Uint64Slice() gocoder.Type {
	return gocoder.NewTypeI([]uint64{})
}

// Uint64SlicePtr func
func Uint64SlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]uint64{})
}

// Float32 func
func Float32() gocoder.Type {
	return gocoder.NewTypeI(float32(0))
}

// Float32Ptr func
func Float32Ptr() gocoder.Type {
	v := float32(0)
	return gocoder.NewTypeI(&v)
}

// Float32Slice func
func Float32Slice() gocoder.Type {
	return gocoder.NewTypeI([]float32{})
}

// Float32SlicePtr func
func Float32SlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]float32{})
}

// Float64 func
func Float64() gocoder.Type {
	return gocoder.NewTypeI(float64(0))
}

// Float64Ptr func
func Float64Ptr() gocoder.Type {
	v := float64(0)
	return gocoder.NewTypeI(&v)
}

// Float64Slice func
func Float64Slice() gocoder.Type {
	return gocoder.NewTypeI([]float64{})
}

// Float64SlicePtr func
func Float64SlicePtr() gocoder.Type {
	return gocoder.NewTypeI(&[]float64{})
}
