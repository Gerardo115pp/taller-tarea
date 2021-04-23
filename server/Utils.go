package main

import (
	"os"
	"strconv"
)

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func stringToInt(s string) int {
	new_int, err := strconv.Atoi(s)
	if err == nil {
		return new_int
	} else {
		panic(err)
	}
}

func stringToUint64(s string) uint64 {
	new_uint64, err := strconv.ParseUint(s, 0, 64)
	if err == nil {
		return new_uint64
	} else {
		panic(err)
	}
}
