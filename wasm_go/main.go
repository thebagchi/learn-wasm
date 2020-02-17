package main

import (
	"fmt"
	"github.com/thebagchi/learn-wasm/wasm_go/utils"
	"runtime"
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
			window.Alert("Hello World")
			keys := utils.Keys(window.JSValue())
			fmt.Println("Len: ", len(keys))
			if len(keys) > 0 {
				for _, key := range keys {
					// fmt.Println(key)
					_ = key
				}
			}
			{
				keys := utils.Keys(js.Global().Get("document"))
				fmt.Println("Len: ", len(keys))
			}
			{
				location := window.Location()
				keys := utils.Keys(location.JSValue())
				fmt.Println("Len: ", len(keys))
				if len(keys) > 0 {
					for _, key := range keys {
						fmt.Println(key)
						_ = key
					}
				}
			}
			document := window.Document()
			if nil != document {
				keys := utils.Keys(document.JSValue())
				fmt.Println("Len: ", len(keys))
				if len(keys) > 0 {
					for _, key := range keys {
						// fmt.Println(key)
						_ = key
					}
				}
				{
					keys := utils.Keys(js.Global().Get("document"))
					fmt.Println("Len: ", len(keys))
				}
				body := document.Body()
				if nil != body {
					element := document.CreateElement("p")
					if nil != element {
						element.SetInnerHtml("Hello from Go!")
						fmt.Println(element.InnerHtml())
						body.AppendChild(element)
					}
				}
			} else {
				fmt.Println("Document is nil")
			}
		} else {
			fmt.Println("Window is nil")
		}
	}
	fmt.Println("Hello World")
	runtime.GC()

	done := make(chan bool, 0)
	<-done
}
