package test

import "fmt"

func Equal(a, b uint64, s string) {
	if a == b {
		fmt.Printf("Test： %s\t\033[32;1m[PASSED]\033[0m\n", s)
	} else {
		fmt.Printf("Test： %s\t\033[33;1m[FAILD ]\033[0m %d != %d \n", s, a, b)
	}
}
