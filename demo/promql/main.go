package main

import (
	"fmt"
	"time"
)

func main() {
	startT := time.Now()
	endT := startT.Round(time.Minute)

	fmt.Println(startT.String(), endT.String())

}
