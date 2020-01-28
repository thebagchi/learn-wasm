package utils

import (
	"syscall/js"
)

func Keys(obj js.Value) []string {
	if !obj.Truthy() {
		return nil
	}
	var (
		keys  = js.Global().Get("Object").Call("keys", obj)
		slice = make([]string, keys.Length())
	)
	for i := 0; i < keys.Length(); i++ {
		slice[i] = keys.Index(i).String()
	}
	return slice
}
