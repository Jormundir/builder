package main

import (
	"builder/dispatcher"
	"fmt"
)

func main() {
	dispatcher, err := dispatcher.NewDispatcher()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	dispatcher.Dispatch()
}
