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
				// fmt.Println(key)
				_ = key
			}
		}
	}
	global := utils.Global{}
	err := utils.Bind(&global, js.Global())
	if nil != err {
		fmt.Println("Error: ", err)
	} else {
		window := global.Window()
		if nil != window {
			keys := utils.Keys(window.JSValue())
			fmt.Println("Len: ", len(keys))
			if len(keys) > 0 {
				for _, key := range keys {
					fmt.Println(key)
				}
			}
			{
				keys := utils.Keys(js.Global().Get("document"))
				fmt.Println("Len: ", len(keys))
			}
			document := window.Document()
			if nil != document {
				keys := utils.Keys(window.JSValue())
				fmt.Println("Len: ", len(keys))
				if len(keys) > 0 {
					for _, key := range keys {
						fmt.Println(key)
					}
				}
				{
					keys := utils.Keys(js.Global().Get("document"))
					fmt.Println("Len: ", len(keys))
				}
			} else {
				fmt.Println("Document is nil")
			}
		} else {
			fmt.Println("Window is nil")
		}
	}
	fmt.Println("Hello World")
}
