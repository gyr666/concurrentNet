package util

import "errors"

func StandRead(i int, dis []byte, capital uint64, RP *uint64) ([]byte, error) {
	if Int2Uint64(i) > capital-*RP {
		return nil, errors.New(INDEX_OUTOF_BOUND)
	}
	bt := make([]byte, i)
	BlockCopy(dis, *RP, bt, 0, Int2Uint64(i))
	*RP = *RP + Int2Uint64(i)
	return bt, nil
}

func StandWrite(des []byte, capital uint64, WP *uint64, _b []byte) error {
	l := Int2Uint64(len(_b))
	if l > capital-*WP {
		return errors.New(INDEX_OUTOF_BOUND)
	}
	BlockCopy(_b, 0, des, *WP, l)
	*WP += l
	return nil
}

func ReadOne(s []byte ,RP *uint64) byte {
	 v:= s[*RP]
	 *RP += 1
	 return v
}

func WriteOne(s []byte,b byte ,WP *uint64)  {
	s[*WP] = b
	*WP += 1
}
