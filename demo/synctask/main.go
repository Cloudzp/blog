package main

import (
	"fmt"
	"time"
)

func main() {
	var ss []string
	go func() {
		ss = modifyArray(ss)
		fmt.Printf("%v", ss)
	}()
	time.Sleep(time.Second)

	fmt.Printf("%v", ss)
}

func modifyArray(ss []string) []string {
	ss = append(ss, "1", "2")
	return ss
}
