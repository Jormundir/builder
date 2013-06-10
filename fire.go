package main

import (
	"fire/dispatcher"
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
