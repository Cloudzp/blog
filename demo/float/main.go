package main

import "fmt"

func main() {
	x := 10e-7
	y := fmt.Sprintf("%.3f", x)

	fmt.Println(x, y)
}
