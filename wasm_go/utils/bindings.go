package utils

import (
	"fmt"
	"reflect"
	"strings"
	"syscall/js"
)

type Global struct {
	js.Value
	Window func() *Window `wasm:"window()"`
}

func (this *Global) JSValue() js.Value {
	return this.Value
}

type Window struct {
	js.Value
}

func (this *Window) JSValue() js.Value {
	return this.Value
}

type Document struct {
	js.Value
}

func (this *Document) JSValue() js.Value {
	return this.Value
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
		if err := BindFunction(tag, value.Type(), object, v, i); nil != err {
			return err
		}
	}
	return nil
}

func BindFunction(tag string, fptr reflect.Type, object js.Value, intf interface{}, index int) error {
	var (
		inputs  = fptr.NumIn()
		outputs = fptr.NumOut()
	)
	if tag, truth := IsFunction(tag); truth {
		fmt.Println(tag, "is function")
	}
	if tag, truth := IsProperty(tag); truth {
		if inputs > 0 {
			return fmt.Errorf("property %s accessed with input parameters, expected function with 0 inputs", tag)
		}
		if outputs != 1 {
			return fmt.Errorf("property %s accessed with no output parameters, expected function with 1 outputs", tag)
		}
	}
	return nil
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

func IsProperty(tag string) (string, bool) {
	tag, cond := IsFunction(tag)
	return tag, !cond
}

func IsFunction(tag string) (string, bool) {
	if strings.HasSuffix(tag, "()") {
		return strings.TrimSuffix(tag, "()"), true
	}
	return tag, false
}
