package utils

import "strconv"

func MustAtoi(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return n
}

func MustAtoui(s string) uint64 {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return n
}

func MustAtof(s string) float64 {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return n
}

func MustAtob(s string) bool {
	n, err := strconv.ParseBool(s)
	if err != nil {
		panic(err)
	}
	return n
}
