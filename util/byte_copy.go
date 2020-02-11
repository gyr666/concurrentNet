package util

func BlockCopy(src []byte, srcOffset uint64, dst []byte, dstOffset, count uint64) bool {
	var index uint64 = 0
	for i := srcOffset; i < srcOffset+count; i++ {
		dst[Uint642Int(dstOffset+index)] = src[Uint642Int(srcOffset+index)]
		index++
	}
	return true
}
