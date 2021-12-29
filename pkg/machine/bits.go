package machine

func signExtend5(v uint16) int16 {
	if v&0b1_0000 == 0 {
		return int16(v)
	}

	return int16(v | 0b1111_1111_1110_0000)
}

func signExtend6(v uint16) int16 {
	if v&0b10_0000 == 0 {
		return int16(v)
	}

	return int16(v | 0b1111_1111_1100_0000)
}

func signExtend9(v uint16) int16 {
	if v&0b1_0000_0000 == 0 {
		return int16(v)
	}

	return int16(v | 0b1111_1110_0000_0000)
}

func signExtend11(v uint16) int16 {
	if v&0b100_0000_0000 == 0 {
		return int16(v)
	}

	return int16(v | 0b1111_1000_0000_0000)
}

func addOffsetU16(base uint16, offset int16) uint16 {
	return uint16(int(base) + int(offset))
}
