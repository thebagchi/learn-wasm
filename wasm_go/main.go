package main

import (
	"fmt"
	"github.com/thebagchi/learn-wasm/wasm_go/utils"
	"syscall/js"
)

func main() {
	{
		keys := utils.Keys(js.Global())
		fmt.Println("Len: ", len(keys))
		if len(keys) > 0 {
			for _, key := range keys {
				fmt.Println(key)
			}
		}
	}
	err := utils.Bind(&utils.Global{}, js.Global())
	if nil != err {
		fmt.Println("Error: ", err)
	}
	fmt.Println("Hello World")
}
