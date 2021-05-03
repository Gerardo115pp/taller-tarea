package main

import (
	"fmt"
	"os"
	"strconv"
)

func clamp(x int, min int, max int) int {
	if x > max {
		return x
	} else if x < min {
		return min
	}
	return x
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func logFatal(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
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
