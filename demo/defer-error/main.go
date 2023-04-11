package main

import "fmt"

func main() {

	doSomething()

}

func doSomething() {
	var err error
	defer func() {
		if err != nil {
			fmt.Println("error: %+v", err) // do something
		}
	}()

	err = fmt.Errorf("error")
}
