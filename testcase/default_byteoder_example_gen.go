// Code generated by generator, DO NOT EDIT.
package test

func (s *DefaultBO) SizeOf() int {
	m := 0
	// L uint16

	// B uint16

	m += (32) / 8
	if (32)%8 != 0 {
		m++
	}

	return m
}

func (s *DefaultBO) MarshalBinary(b []byte) int {
	if len(b) < s.SizeOf() {
		return 0
	}
	m := 0
	// L uint16

	_ = b[(m + 1)] // bounds check hint to compiler; see golang.org/issue/14808
	b[(m + 0)] = byte(uint16(s.L) >> 0)
	b[(m + 1)] = byte(uint16(s.L) >> 8)

	// B uint16

	_ = b[(m + 3)] // bounds check hint to compiler; see golang.org/issue/14808
	b[(m + 2)] = byte(uint16(s.B) >> 8)
	b[(m + 3)] = byte(uint16(s.B) >> 0)

	m += (32) / 8
	if (32)%8 != 0 {
		m++
	}

	return m
}

func (s *DefaultBO) UnmarshalBinary(b []byte) int {
	m := 0
	// L uint16

	if 16+m*8+16 > len(b)*8 {
		return 0
	}

	s.L = int(uint16(b[(m+0)])<<0 | uint16(b[(m+1)])<<8)
	// B uint16

	s.B = int(uint16(b[(m+2)])<<8 | uint16(b[(m+3)])<<0)
	m += (32) / 8
	if (32)%8 != 0 {
		m++
	}

	return m
}
