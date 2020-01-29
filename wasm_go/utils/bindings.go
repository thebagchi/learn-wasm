package utils

import (
	"fmt"
	"reflect"
	"syscall/js"
)

type Global struct {
	js.Value
	Window func() *Window `wasm:"window"`
}

type Window struct {
	js.Value
}

type Document struct {
	js.Value
}

func Bind(v interface{}, object js.Value) error {
	if err := IsStructurePtr(v); nil != err {
		return err
	}
	if err := BindValue(v, object); nil != err {
		return err
	}
	value := reflect.ValueOf(v).Elem()
	for i, field := range Members(v) {
		tag := field.Tag.Get("wasm")
		if len(tag) == 0 {
			continue
		}
		value := value.Field(i)
		if kind := value.Type().Kind(); kind != reflect.Func {
			return fmt.Errorf("field of type %s found, func expected", kind)
		}
	}
	return nil
}

func BindFunction() {

}

func BindValue(v interface{}, object js.Value) error {
	value := reflect.ValueOf(v).Elem()
	for i, field := range Members(v) {
		if reflect.TypeOf(object) == field.Type {
			value := value.Field(i)
			value.Set(reflect.ValueOf(object))
			return nil
		}
	}
	return fmt.Errorf("field of type %s not found", reflect.TypeOf(object))
}

func IsStructurePtr(v interface{}) error {
	t := reflect.TypeOf(v)
	if kind := t.Kind(); kind != reflect.Ptr {
		return fmt.Errorf("%s received, ptr to struct expected", kind)
	}
	if kind := t.Elem().Kind(); kind != reflect.Struct {
		return fmt.Errorf("ptr to %s received, ptr to struct expected", kind)
	}
	return nil
}

func Members(v interface{}) []reflect.StructField {
	var (
		elem   = reflect.TypeOf(v).Elem()
		count  = elem.NumField()
		fields = make([]reflect.StructField, count)
	)
	for i := 0; i < count; i++ {
		fields[i] = elem.Field(i)
	}
	return fields
}
