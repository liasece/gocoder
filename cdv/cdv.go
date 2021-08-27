package cdv

import (
	"time"

	"github.com/liasece/gocoder"
)

// I func
func I(i interface{}) gocoder.Value {
	return gocoder.NewValueNameI("", i)
}

// Func func
func Func(name string, argTypes []gocoder.Type, returns []gocoder.Type) gocoder.Value {
	return gocoder.NewValueFunc(name, nil, argTypes, returns)
}

// NameI func
func NameI(name string, i interface{}) gocoder.Value {
	return gocoder.NewValueNameI(name, i)
}

// NameType func
func NameType(name string, t gocoder.Type) gocoder.Value {
	return gocoder.NewValue(name, t)
}

// Nil func
func Nil() gocoder.Value {
	return gocoder.NewValueNil()
}

// None func
func None() gocoder.Value {
	return gocoder.NewValueNone()
}

// Time func
func Time(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, time.Time{})
}

// TimePtr func
func TimePtr(name string) gocoder.Value {
	v := time.Time{}
	return gocoder.NewValueNameI(name, &v)
}

// String func
func String(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, string(""))
}

// StringPtr func
func StringPtr(name string) gocoder.Value {
	v := string("")
	return gocoder.NewValueNameI(name, &v)
}

// StringSlice func
func StringSlice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []string{})
}

// StringSlicePtr func
func StringSlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]string{})
}

// Bool func
func Bool(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, bool(false))
}

// BoolPtr func
func BoolPtr(name string) gocoder.Value {
	v := bool(false)
	return gocoder.NewValueNameI(name, &v)
}

// BoolSlice func
func BoolSlice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []bool{})
}

// BoolSlicePtr func
func BoolSlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]bool{})
}

// Int func
func Int(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, int(0))
}

// IntPtr func
func IntPtr(name string) gocoder.Value {
	v := int(0)
	return gocoder.NewValueNameI(name, &v)
}

// IntSlice func
func IntSlice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []int{})
}

// IntSlicePtr func
func IntSlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]int{})
}

// Int8 func
func Int8(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, int8(0))
}

// Int8Ptr func
func Int8Ptr(name string) gocoder.Value {
	v := int8(0)
	return gocoder.NewValueNameI(name, &v)
}

// Int8Slice func
func Int8Slice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []int8{})
}

// Int8SlicePtr func
func Int8SlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]int8{})
}

// Int16 func
func Int16(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, int16(0))
}

// Int16Ptr func
func Int16Ptr(name string) gocoder.Value {
	v := int16(0)
	return gocoder.NewValueNameI(name, &v)
}

// Int16Slice func
func Int16Slice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []int16{})
}

// Int16SlicePtr func
func Int16SlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]int16{})
}

// Int32 func
func Int32(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, int32(0))
}

// Int32Ptr func
func Int32Ptr(name string) gocoder.Value {
	v := int32(0)
	return gocoder.NewValueNameI(name, &v)
}

// Int32Slice func
func Int32Slice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []int32{})
}

// Int32SlicePtr func
func Int32SlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]int32{})
}

// Int64 func
func Int64(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, int64(0))
}

// Int64Ptr func
func Int64Ptr(name string) gocoder.Value {
	v := int64(0)
	return gocoder.NewValueNameI(name, &v)
}

// Int64Slice func
func Int64Slice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []int64{})
}

// Int64SlicePtr func
func Int64SlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]int64{})
}

// Uint func
func Uint(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, uint(0))
}

// UintPtr func
func UintPtr(name string) gocoder.Value {
	v := uint(0)
	return gocoder.NewValueNameI(name, &v)
}

// UintSlice func
func UintSlice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []uint{})
}

// UintSlicePtr func
func UintSlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]uint{})
}

// Uint8 func
func Uint8(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, uint8(0))
}

// Uint8Ptr func
func Uint8Ptr(name string) gocoder.Value {
	v := uint8(0)
	return gocoder.NewValueNameI(name, &v)
}

// Uint8Slice func
func Uint8Slice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []uint8{})
}

// Uint8SlicePtr func
func Uint8SlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]uint8{})
}

// Byte func
func Byte(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, byte(0))
}

// BytePtr func
func BytePtr(name string) gocoder.Value {
	v := byte(0)
	return gocoder.NewValueNameI(name, &v)
}

// ByteSlice func
func ByteSlice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []byte{})
}

// ByteSlicePtr func
func ByteSlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]byte{})
}

// Uint16 func
func Uint16(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, uint16(0))
}

// Uint16Ptr func
func Uint16Ptr(name string) gocoder.Value {
	v := uint16(0)
	return gocoder.NewValueNameI(name, &v)
}

// Uint16Slice func
func Uint16Slice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []uint16{})
}

// Uint16SlicePtr func
func Uint16SlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]uint16{})
}

// Uint32 func
func Uint32(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, uint32(0))
}

// Uint32Ptr func
func Uint32Ptr(name string) gocoder.Value {
	v := uint32(0)
	return gocoder.NewValueNameI(name, &v)
}

// Uint32Slice func
func Uint32Slice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []uint32{})
}

// Uint32SlicePtr func
func Uint32SlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]uint32{})
}

// Uint64 func
func Uint64(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, uint64(0))
}

// Uint64Ptr func
func Uint64Ptr(name string) gocoder.Value {
	v := uint64(0)
	return gocoder.NewValueNameI(name, &v)
}

// Uint64Slice func
func Uint64Slice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []uint64{})
}

// Uint64SlicePtr func
func Uint64SlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]uint64{})
}

// Float32 func
func Float32(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, float32(0))
}

// Float32Ptr func
func Float32Ptr(name string) gocoder.Value {
	v := float32(0)
	return gocoder.NewValueNameI(name, &v)
}

// Float32Slice func
func Float32Slice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []float32{})
}

// Float32SlicePtr func
func Float32SlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]float32{})
}

// Float64 func
func Float64(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, float64(0))
}

// Float64Ptr func
func Float64Ptr(name string) gocoder.Value {
	v := float64(0)
	return gocoder.NewValueNameI(name, &v)
}

// Float64Slice func
func Float64Slice(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, []float64{})
}

// Float64SlicePtr func
func Float64SlicePtr(name string) gocoder.Value {
	return gocoder.NewValueNameI(name, &[]float64{})
}
