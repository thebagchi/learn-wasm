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
	fmt.Println("Hello World")
}
