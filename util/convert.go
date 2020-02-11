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

func IsPow2(u uint64) (int,bool){
	count := 0
	max :=0
	var bin uint64 = 1
	for ;bin<4294967295;max++ {
		if (bin & u)!=0 {
			count++
		}
		bin = bin << 1
	}
	return max,count==1
}