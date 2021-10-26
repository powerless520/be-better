package main

import (
	"be-better/core"
	"fmt"
	"runtime/debug"
)

func main() {

	defer func() {
		fmt.Println("main defer caller")
		if err := recover(); err != nil {
			fmt.Printf("recover err:%+v\n"+string(debug.Stack()), err)
		}
	}()

	core.GlobalInit()
}
