package util

func Int2Uint64(i int) uint64 {
	return uint64(i)
}

func Uint642Int(i uint64) int {
	return int(i)
}

func Uint642UInt8(i uint64) uint8 {
	return uint8(i)
}

func Int2UInt8(i int) uint8 {
	return uint8(i)
}

func IsPow2(u uint64) []uint8 {
	v := make([]uint8, 0)
	i := 0
	var b uint64 = 1
	for ; i < 63; i++ {
		if (u & (b << i)) != 0 {
			if i == 0 {
				v = append(v, 0)
			} else {
				v = append(v, uint8(i)-1)
			}
		}
	}
	return v
}
