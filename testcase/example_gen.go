// Code generated by generator, DO NOT EDIT.
package test

import (
	"math"
	"unsafe"
)

func (s *Example) MarshalBinary(b []byte) int {
	m := 0
	// Pad
	for i := 0; i < int(5); i++ {
		b[m+i] = 0
	}
	m += int(5)
	// I8f
	b[m] = (uint8)(int8(s.I8f))
	m++
	// I16f
	_ = b[m+1] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte((uint16)(int16(s.I16f)) >> 8)
	b[m+1] = byte((uint16)(int16(s.I16f)) >> 0)
	m += 2
	// I32f
	_ = b[m+3] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte((uint32)(int32(s.I32f)) >> 24)
	b[m+1] = byte((uint32)(int32(s.I32f)) >> 16)
	b[m+2] = byte((uint32)(int32(s.I32f)) >> 8)
	b[m+3] = byte((uint32)(int32(s.I32f)) >> 0)
	m += 4
	// I64f
	_ = b[m+7] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte((uint64)(int64(s.I64f)) >> 56)
	b[m+1] = byte((uint64)(int64(s.I64f)) >> 48)
	b[m+2] = byte((uint64)(int64(s.I64f)) >> 40)
	b[m+3] = byte((uint64)(int64(s.I64f)) >> 32)
	b[m+4] = byte((uint64)(int64(s.I64f)) >> 24)
	b[m+5] = byte((uint64)(int64(s.I64f)) >> 16)
	b[m+6] = byte((uint64)(int64(s.I64f)) >> 8)
	b[m+7] = byte((uint64)(int64(s.I64f)) >> 0)
	m += 8
	// U8f
	b[m] = uint8(s.U8f)
	m++
	// U16f
	_ = b[m+1] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte(uint16(s.U16f) >> 0)
	b[m+1] = byte(uint16(s.U16f) >> 8)
	m += 2
	// U32f
	_ = b[m+3] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte(uint32(s.U32f) >> 0)
	b[m+1] = byte(uint32(s.U32f) >> 8)
	b[m+2] = byte(uint32(s.U32f) >> 16)
	b[m+3] = byte(uint32(s.U32f) >> 24)
	m += 4
	// U64f
	_ = b[m+7] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte(uint64(s.U64f) >> 0)
	b[m+1] = byte(uint64(s.U64f) >> 8)
	b[m+2] = byte(uint64(s.U64f) >> 16)
	b[m+3] = byte(uint64(s.U64f) >> 24)
	b[m+4] = byte(uint64(s.U64f) >> 32)
	b[m+5] = byte(uint64(s.U64f) >> 40)
	b[m+6] = byte(uint64(s.U64f) >> 48)
	b[m+7] = byte(uint64(s.U64f) >> 56)
	m += 8
	// Boolf
	b[m] = uint8(s.Boolf & 1)
	m++
	// Byte4f
	for i := 0; i < int(4); i++ {
		b[m] = byte(s.Byte4f[i])
		m++
	}
	// I8
	b[m] = (uint8)(int8(s.I8))
	m++
	// I16
	_ = b[m+1] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte((uint16)(int16(s.I16)) >> 8)
	b[m+1] = byte((uint16)(int16(s.I16)) >> 0)
	m += 2
	// I32
	_ = b[m+3] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte((uint32)(int32(s.I32)) >> 24)
	b[m+1] = byte((uint32)(int32(s.I32)) >> 16)
	b[m+2] = byte((uint32)(int32(s.I32)) >> 8)
	b[m+3] = byte((uint32)(int32(s.I32)) >> 0)
	m += 4
	// I64
	_ = b[m+7] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte((uint64)(int64(s.I64)) >> 56)
	b[m+1] = byte((uint64)(int64(s.I64)) >> 48)
	b[m+2] = byte((uint64)(int64(s.I64)) >> 40)
	b[m+3] = byte((uint64)(int64(s.I64)) >> 32)
	b[m+4] = byte((uint64)(int64(s.I64)) >> 24)
	b[m+5] = byte((uint64)(int64(s.I64)) >> 16)
	b[m+6] = byte((uint64)(int64(s.I64)) >> 8)
	b[m+7] = byte((uint64)(int64(s.I64)) >> 0)
	m += 8
	// U8
	b[m] = uint8(s.U8)
	m++
	// U16
	_ = b[m+1] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte(uint16(s.U16) >> 0)
	b[m+1] = byte(uint16(s.U16) >> 8)
	m += 2
	// U32
	_ = b[m+3] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte(uint32(s.U32) >> 0)
	b[m+1] = byte(uint32(s.U32) >> 8)
	b[m+2] = byte(uint32(s.U32) >> 16)
	b[m+3] = byte(uint32(s.U32) >> 24)
	m += 4
	// U64
	_ = b[m+7] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte(uint64(s.U64) >> 0)
	b[m+1] = byte(uint64(s.U64) >> 8)
	b[m+2] = byte(uint64(s.U64) >> 16)
	b[m+3] = byte(uint64(s.U64) >> 24)
	b[m+4] = byte(uint64(s.U64) >> 32)
	b[m+5] = byte(uint64(s.U64) >> 40)
	b[m+6] = byte(uint64(s.U64) >> 48)
	b[m+7] = byte(uint64(s.U64) >> 56)
	m += 8
	// BoolT
	b[m] = uint8(*(*uint8)(unsafe.Pointer(&s.BoolT)) & 1 & 1)
	m++
	// BoolF
	b[m] = uint8(*(*uint8)(unsafe.Pointer(&s.BoolF)) & 1 & 1)
	m++
	// Byte4
	for i := 0; i < int(4); i++ {
		b[m] = byte(s.Byte4[i])
		m++
	}
	// Float1
	{
		tmp := math.Float32bits(float32(float32(s.Float1)))
		_ = b[m+3] // bounds check hint to compiler; see golang.org/issue/14808
		b[m+0] = byte(tmp >> 24)
		b[m+1] = byte(tmp >> 16)
		b[m+2] = byte(tmp >> 8)
		b[m+3] = byte(tmp >> 0)
		m += 4
	}
	// Float2
	{
		tmp := math.Float64bits(float64(float64(s.Float2)))
		_ = b[m+7] // bounds check hint to compiler; see golang.org/issue/14808
		b[m+0] = byte(tmp >> 56)
		b[m+1] = byte(tmp >> 48)
		b[m+2] = byte(tmp >> 40)
		b[m+3] = byte(tmp >> 32)
		b[m+4] = byte(tmp >> 24)
		b[m+5] = byte(tmp >> 16)
		b[m+6] = byte(tmp >> 8)
		b[m+7] = byte(tmp >> 0)
		m += 8
	}
	// I32f2
	_ = b[m+3] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte((uint32)(int32(s.I32f2)) >> 24)
	b[m+1] = byte((uint32)(int32(s.I32f2)) >> 16)
	b[m+2] = byte((uint32)(int32(s.I32f2)) >> 8)
	b[m+3] = byte((uint32)(int32(s.I32f2)) >> 0)
	m += 4
	// U32f2
	_ = b[m+3] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte(uint32(s.U32f2) >> 24)
	b[m+1] = byte(uint32(s.U32f2) >> 16)
	b[m+2] = byte(uint32(s.U32f2) >> 8)
	b[m+3] = byte(uint32(s.U32f2) >> 0)
	m += 4
	// I32f3
	_ = b[m+7] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte((uint64)(int64(s.I32f3)) >> 56)
	b[m+1] = byte((uint64)(int64(s.I32f3)) >> 48)
	b[m+2] = byte((uint64)(int64(s.I32f3)) >> 40)
	b[m+3] = byte((uint64)(int64(s.I32f3)) >> 32)
	b[m+4] = byte((uint64)(int64(s.I32f3)) >> 24)
	b[m+5] = byte((uint64)(int64(s.I32f3)) >> 16)
	b[m+6] = byte((uint64)(int64(s.I32f3)) >> 8)
	b[m+7] = byte((uint64)(int64(s.I32f3)) >> 0)
	m += 8
	// Size1
	s.Size1 = int(len(s.Str))
	_ = b[m+3] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte((uint32)(int(s.Size1)) >> 0)
	b[m+1] = byte((uint32)(int(s.Size1)) >> 8)
	b[m+2] = byte((uint32)(int(s.Size1)) >> 16)
	b[m+3] = byte((uint32)(int(s.Size1)) >> 24)
	m += 4
	// Str
	for i := copy(b[m:m+int(s.Size1)], []byte(s.Str)); i < int(s.Size1); i++ {
		b[m+i] = 0
	}
	m += int(s.Size1)
	// Strb
	for i := copy(b[m:m+int(4)], []byte(s.Strb)); i < int(4); i++ {
		b[m+i] = 0
	}
	m += int(4)
	// Size2
	s.Size2 = int(len(s.Str2))
	b[m] = uint8(s.Size2)
	m++
	// Str2
	for i := copy(b[m:m+int(s.Size2)], []byte(s.Str2)); i < int(s.Size2); i++ {
		b[m+i] = 0
	}
	m += int(s.Size2)
	// Size3
	s.Size3 = int(len(s.Bstr))
	b[m] = uint8(s.Size3)
	m++
	// Bstr
	for i := 0; i < int(s.Size3); i++ {
		b[m] = byte(s.Bstr[i])
		m++
	}
	// Size4
	_ = b[m+3] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte((uint32)(int(s.Size4)) >> 0)
	b[m+1] = byte((uint32)(int(s.Size4)) >> 8)
	b[m+2] = byte((uint32)(int(s.Size4)) >> 16)
	b[m+3] = byte((uint32)(int(s.Size4)) >> 24)
	m += 4
	// Str4a
	for i := copy(b[m:m+int(s.Size4)], []byte(s.Str4a)); i < int(s.Size4); i++ {
		b[m+i] = 0
	}
	m += int(s.Size4)
	// Str4b
	for i := copy(b[m:m+int(s.Size4)], []byte(s.Str4b)); i < int(s.Size4); i++ {
		b[m+i] = 0
	}
	m += int(s.Size4)
	// Size5
	b[m] = uint8(s.Size5)
	m++
	// Bstr2
	for i := 0; i < int(s.Size5); i++ {
		b[m] = byte(s.Bstr2[i])
		m++
	}
	// Nested
	m += s.Nested.MarshalBinary(b[m:])
	// NestedP
	m += (*s.NestedP).MarshalBinary(b[m:])
	// TestP64
	_ = b[m+7] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte((uint64)(int64((*s.TestP64))) >> 56)
	b[m+1] = byte((uint64)(int64((*s.TestP64))) >> 48)
	b[m+2] = byte((uint64)(int64((*s.TestP64))) >> 40)
	b[m+3] = byte((uint64)(int64((*s.TestP64))) >> 32)
	b[m+4] = byte((uint64)(int64((*s.TestP64))) >> 24)
	b[m+5] = byte((uint64)(int64((*s.TestP64))) >> 16)
	b[m+6] = byte((uint64)(int64((*s.TestP64))) >> 8)
	b[m+7] = byte((uint64)(int64((*s.TestP64))) >> 0)
	m += 8
	// NestedSize
	s.NestedSize = int(len(s.NestedA))
	_ = b[m+3] // bounds check hint to compiler; see golang.org/issue/14808
	b[m+0] = byte((uint32)(int(s.NestedSize)) >> 24)
	b[m+1] = byte((uint32)(int(s.NestedSize)) >> 16)
	b[m+2] = byte((uint32)(int(s.NestedSize)) >> 8)
	b[m+3] = byte((uint32)(int(s.NestedSize)) >> 0)
	m += 4
	// NestedA
	for i := 0; i < int(s.NestedSize); i++ {
		m += s.NestedA[i].MarshalBinary(b[m:])
	}
	// CustomTypeSize
	s.CustomTypeSize = Int3(len(s.CustomTypeSizeArr))
	m += s.CustomTypeSize.MarshalBinary(b[m:])
	// CustomTypeSizeArr
	for i := 0; i < int(s.CustomTypeSize); i++ {
		b[m] = byte(s.CustomTypeSizeArr[i])
		m++
	}
	return m
}
func (s *Example) UnmarshalBinary(b []byte) int {
	m := 0
	// Pad
	m += int(5)
	// I8f
	if m+1 > len(b) {
		return m
	}
	s.I8f = int(int8(b[m]))
	m += 1
	// I16f
	if m+2 > len(b) {
		return m
	}
	s.I16f = int(int16(uint16(b[m+0])<<8 | uint16(b[m+1])<<0))
	m += 2
	// I32f
	if m+4 > len(b) {
		return m
	}
	s.I32f = int(int32(uint32(b[m+0])<<24 | uint32(b[m+1])<<16 | uint32(b[m+2])<<8 | uint32(b[m+3])<<0))
	m += 4
	// I64f
	if m+8 > len(b) {
		return m
	}
	s.I64f = int(int64(uint64(b[m+0])<<56 | uint64(b[m+1])<<48 | uint64(b[m+2])<<40 | uint64(b[m+3])<<32 | uint64(b[m+4])<<24 | uint64(b[m+5])<<16 | uint64(b[m+6])<<8 | uint64(b[m+7])<<0))
	m += 8
	// U8f
	if m+1 > len(b) {
		return m
	}
	s.U8f = int(b[m])
	m += 1
	// U16f
	if m+2 > len(b) {
		return m
	}
	s.U16f = int(uint16(b[m+0])<<0 | uint16(b[m+1])<<8)
	m += 2
	// U32f
	if m+4 > len(b) {
		return m
	}
	s.U32f = int(uint32(b[m+0])<<0 | uint32(b[m+1])<<8 | uint32(b[m+2])<<16 | uint32(b[m+3])<<24)
	m += 4
	// U64f
	if m+8 > len(b) {
		return m
	}
	s.U64f = int(uint64(b[m+0])<<0 | uint64(b[m+1])<<8 | uint64(b[m+2])<<16 | uint64(b[m+3])<<24 | uint64(b[m+4])<<32 | uint64(b[m+5])<<40 | uint64(b[m+6])<<48 | uint64(b[m+7])<<56)
	m += 8
	// Boolf
	if m+1 > len(b) {
		return m
	}
	s.Boolf = int(b[m]) & 1
	m += 1
	// Byte4f
	if len(s.Byte4f) < int(4) {
		s.Byte4f = make([]byte, int(4))
	}
	for i := 0; i < int(4); i++ {
		if m+1 > len(b) {
			return m
		}
		s.Byte4f[i] = byte(b[m])
		m += 1
	}
	// I8
	if m+1 > len(b) {
		return m
	}
	s.I8 = int8(int8(b[m]))
	m += 1
	// I16
	if m+2 > len(b) {
		return m
	}
	s.I16 = int16(int16(uint16(b[m+0])<<8 | uint16(b[m+1])<<0))
	m += 2
	// I32
	if m+4 > len(b) {
		return m
	}
	s.I32 = int32(int32(uint32(b[m+0])<<24 | uint32(b[m+1])<<16 | uint32(b[m+2])<<8 | uint32(b[m+3])<<0))
	m += 4
	// I64
	if m+8 > len(b) {
		return m
	}
	s.I64 = int64(int64(uint64(b[m+0])<<56 | uint64(b[m+1])<<48 | uint64(b[m+2])<<40 | uint64(b[m+3])<<32 | uint64(b[m+4])<<24 | uint64(b[m+5])<<16 | uint64(b[m+6])<<8 | uint64(b[m+7])<<0))
	m += 8
	// U8
	if m+1 > len(b) {
		return m
	}
	s.U8 = uint8(b[m])
	m += 1
	// U16
	if m+2 > len(b) {
		return m
	}
	s.U16 = uint16(uint16(b[m+0])<<0 | uint16(b[m+1])<<8)
	m += 2
	// U32
	if m+4 > len(b) {
		return m
	}
	s.U32 = uint32(uint32(b[m+0])<<0 | uint32(b[m+1])<<8 | uint32(b[m+2])<<16 | uint32(b[m+3])<<24)
	m += 4
	// U64
	if m+8 > len(b) {
		return m
	}
	s.U64 = uint64(uint64(b[m+0])<<0 | uint64(b[m+1])<<8 | uint64(b[m+2])<<16 | uint64(b[m+3])<<24 | uint64(b[m+4])<<32 | uint64(b[m+5])<<40 | uint64(b[m+6])<<48 | uint64(b[m+7])<<56)
	m += 8
	// BoolT
	if m+1 > len(b) {
		return m
	}
	s.BoolT = uint8(b[m]) > 0
	m += 1
	// BoolF
	if m+1 > len(b) {
		return m
	}
	s.BoolF = uint8(b[m]) > 0
	m += 1
	// Byte4
	for i := 0; i < int(4); i++ {
		if m+1 > len(b) {
			return m
		}
		s.Byte4[i] = byte(b[m])
		m += 1
	}
	// Float1
	if m+4 > len(b) {
		return m
	}
	s.Float1 = float32(math.Float32frombits(uint32(b[m+0])<<24 | uint32(b[m+1])<<16 | uint32(b[m+2])<<8 | uint32(b[m+3])<<0))
	m += 4
	// Float2
	if m+8 > len(b) {
		return m
	}
	s.Float2 = float64(math.Float64frombits(uint64(b[m+0])<<56 | uint64(b[m+1])<<48 | uint64(b[m+2])<<40 | uint64(b[m+3])<<32 | uint64(b[m+4])<<24 | uint64(b[m+5])<<16 | uint64(b[m+6])<<8 | uint64(b[m+7])<<0))
	m += 8
	// I32f2
	if m+4 > len(b) {
		return m
	}
	s.I32f2 = int64(int32(uint32(b[m+0])<<24 | uint32(b[m+1])<<16 | uint32(b[m+2])<<8 | uint32(b[m+3])<<0))
	m += 4
	// U32f2
	if m+4 > len(b) {
		return m
	}
	s.U32f2 = int64(uint32(b[m+0])<<24 | uint32(b[m+1])<<16 | uint32(b[m+2])<<8 | uint32(b[m+3])<<0)
	m += 4
	// I32f3
	if m+8 > len(b) {
		return m
	}
	s.I32f3 = int32(int64(uint64(b[m+0])<<56 | uint64(b[m+1])<<48 | uint64(b[m+2])<<40 | uint64(b[m+3])<<32 | uint64(b[m+4])<<24 | uint64(b[m+5])<<16 | uint64(b[m+6])<<8 | uint64(b[m+7])<<0))
	m += 8
	// Size1
	if m+4 > len(b) {
		return m
	}
	s.Size1 = int(int32(uint32(b[m+0])<<0 | uint32(b[m+1])<<8 | uint32(b[m+2])<<16 | uint32(b[m+3])<<24))
	m += 4
	// Str
	if m+int(s.Size1) > len(b) {
		return m
	}
	s.Str = string(b[m : m+int(s.Size1)])
	m += int(s.Size1)
	// Strb
	if m+int(4) > len(b) {
		return m
	}
	s.Strb = string(b[m : m+int(4)])
	m += int(4)
	// Size2
	if m+1 > len(b) {
		return m
	}
	s.Size2 = int(b[m])
	m += 1
	// Str2
	if m+int(s.Size2) > len(b) {
		return m
	}
	s.Str2 = string(b[m : m+int(s.Size2)])
	m += int(s.Size2)
	// Size3
	if m+1 > len(b) {
		return m
	}
	s.Size3 = int(b[m])
	m += 1
	// Bstr
	if len(s.Bstr) < int(s.Size3) {
		s.Bstr = make([]byte, int(s.Size3))
	}
	for i := 0; i < int(s.Size3); i++ {
		if m+1 > len(b) {
			return m
		}
		s.Bstr[i] = byte(b[m])
		m += 1
	}
	// Size4
	if m+4 > len(b) {
		return m
	}
	s.Size4 = int(int32(uint32(b[m+0])<<0 | uint32(b[m+1])<<8 | uint32(b[m+2])<<16 | uint32(b[m+3])<<24))
	m += 4
	// Str4a
	if m+int(s.Size4) > len(b) {
		return m
	}
	s.Str4a = string(b[m : m+int(s.Size4)])
	m += int(s.Size4)
	// Str4b
	if m+int(s.Size4) > len(b) {
		return m
	}
	s.Str4b = string(b[m : m+int(s.Size4)])
	m += int(s.Size4)
	// Size5
	if m+1 > len(b) {
		return m
	}
	s.Size5 = int(b[m])
	m += 1
	// Bstr2
	if len(s.Bstr2) < int(s.Size5) {
		s.Bstr2 = make([]byte, int(s.Size5))
	}
	for i := 0; i < int(s.Size5); i++ {
		if m+1 > len(b) {
			return m
		}
		s.Bstr2[i] = byte(b[m])
		m += 1
	}
	// Nested
	m += s.Nested.UnmarshalBinary(b[m:])
	// NestedP
	if s.NestedP == nil {
		s.NestedP = new(Nested)
	}
	m += (*s.NestedP).UnmarshalBinary(b[m:])
	// TestP64
	if s.TestP64 == nil {
		s.TestP64 = new(int)
	}
	if m+8 > len(b) {
		return m
	}
	(*s.TestP64) = int(int64(uint64(b[m+0])<<56 | uint64(b[m+1])<<48 | uint64(b[m+2])<<40 | uint64(b[m+3])<<32 | uint64(b[m+4])<<24 | uint64(b[m+5])<<16 | uint64(b[m+6])<<8 | uint64(b[m+7])<<0))
	m += 8
	// NestedSize
	if m+4 > len(b) {
		return m
	}
	s.NestedSize = int(int32(uint32(b[m+0])<<24 | uint32(b[m+1])<<16 | uint32(b[m+2])<<8 | uint32(b[m+3])<<0))
	m += 4
	// NestedA
	if len(s.NestedA) < int(s.NestedSize) {
		s.NestedA = make([]Nested, int(s.NestedSize))
	}
	for i := 0; i < int(s.NestedSize); i++ {
		m += s.NestedA[i].UnmarshalBinary(b[m:])
	}
	// CustomTypeSize
	m += s.CustomTypeSize.UnmarshalBinary(b[m:])
	// CustomTypeSizeArr
	if len(s.CustomTypeSizeArr) < int(s.CustomTypeSize) {
		s.CustomTypeSizeArr = make([]byte, int(s.CustomTypeSize))
	}
	for i := 0; i < int(s.CustomTypeSize); i++ {
		if m+1 > len(b) {
			return m
		}
		s.CustomTypeSizeArr[i] = byte(b[m])
		m += 1
	}
	return m
}
func (s *Example) SizeOf() int {
	m := 0
	m += int(5)
	m += 1
	m += 2
	m += 4
	m += 8
	m += 1
	m += 2
	m += 4
	m += 8
	m += 1
	for i := 0; i < int(4); i++ {
		m += 1
	}
	m += 1
	m += 2
	m += 4
	m += 8
	m += 1
	m += 2
	m += 4
	m += 8
	m += 1
	m += 1
	for i := 0; i < int(4); i++ {
		m += 1
	}
	m += 4
	m += 8
	m += 4
	m += 4
	m += 8
	s.Size1 = int(len(s.Str))
	m += 4

	m += int(s.Size1)

	m += int(4)
	s.Size2 = int(len(s.Str2))
	m += 1

	m += int(s.Size2)
	s.Size3 = int(len(s.Bstr))
	m += 1
	for i := 0; i < int(s.Size3); i++ {
		m += 1
	}
	m += 4

	m += int(s.Size4)

	m += int(s.Size4)
	m += 1
	for i := 0; i < int(s.Size5); i++ {
		m += 1
	}
	m += s.Nested.SizeOf()
	m += (*s.NestedP).SizeOf()
	m += 8
	s.NestedSize = int(len(s.NestedA))
	m += 4
	for i := 0; i < int(s.NestedSize); i++ {
		m += s.NestedA[i].SizeOf()
	}
	s.CustomTypeSize = Int3(len(s.CustomTypeSizeArr))
	m += s.CustomTypeSize.SizeOf()
	for i := 0; i < int(s.CustomTypeSize); i++ {
		m += 1
	}
	return m
}
func (s *Nested) MarshalBinary(b []byte) int {
	m := 0
	// Test2
	b[m] = (uint8)(int8(s.Test2))
	m++
	return m
}
func (s *Nested) UnmarshalBinary(b []byte) int {
	m := 0
	// Test2
	if m+1 > len(b) {
		return m
	}
	s.Test2 = int(int8(b[m]))
	m += 1
	return m
}
func (s *Nested) SizeOf() int {
	m := 0
	m += 1
	return m
}
