package utils

import (
	"fmt"
	"reflect"
	"strings"
	"syscall/js"
)

type Global struct {
	js.Value
	Window   func() *Window   `wasm:"window"`
	Document func() *Document `wasm:"document"`
}

func (obj *Global) JSValue() js.Value {
	return obj.Value
}

type Window struct {
	js.Value
}

func (obj *Window) JSValue() js.Value {
	return obj.Value
}

type Document struct {
	js.Value
}

func (obj *Document) JSValue() js.Value {
	return obj.Value
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

func BindFunction(tag string, fptr reflect.Type, object js.Value, value interface{}, index int) error {
	var (
		inputs  = fptr.NumIn()
		outputs = fptr.NumOut()
	)
	if tag, truth := IsFunction(tag); truth {
		// TODO: (Handle Functions)
		fmt.Println(tag, "is function")
	}
	if tag, truth := IsProperty(tag); truth {
		if inputs > 0 {
			return fmt.Errorf("property %s accessed with input parameters, expected function with 0 inputs", tag)
		}
		if outputs != 1 {
			return fmt.Errorf("property %s accessed with no output parameters, expected function with 1 outputs", tag)
		}
		returns := fptr.Out(0)
		switch returns.Kind() {
		case reflect.String:
			reflect.ValueOf(value).Elem().Field(index).Set(
				reflect.ValueOf(func() string {
					return object.Get(tag).String()
				}),
			)
		case reflect.Bool:
			reflect.ValueOf(value).Elem().Field(index).Set(
				reflect.ValueOf(func() bool {
					return object.Get(tag).Bool()
				}),
			)
		case reflect.Int:
			reflect.ValueOf(value).Elem().Field(index).Set(
				reflect.ValueOf(func() int {
					return object.Get(tag).Int()
				}),
			)
		case reflect.Float64:
			reflect.ValueOf(value).Elem().Field(index).Set(
				reflect.ValueOf(func() float64 {
					return object.Get(tag).Float()
				}),
			)
		case reflect.Map:
			// TODO: (Handling for Map)
			return fmt.Errorf("property %s, unsupported kind %s", tag, returns.Kind())
		case reflect.Array:
			// TODO: (Handling for Array)
			return fmt.Errorf("property %s, unsupported kind %s", tag, returns.Kind())
		case reflect.Ptr:
			if kind := returns.Elem().Kind(); kind != reflect.Struct {
				return fmt.Errorf("property %s, unsupported kind %s", tag, returns.Kind())
			}
			function := reflect.FuncOf([]reflect.Type{}, []reflect.Type{returns}, false)
			reflect.ValueOf(value).Elem().Field(index).Set(
				reflect.MakeFunc(function, func(args []reflect.Value) []reflect.Value {
					v := reflect.New(returns.Elem()).Interface()
					err := Bind(v, object.Get(tag))
					if nil == err {
						return []reflect.Value{reflect.ValueOf(v)}
					} else {
						fmt.Println("Error: ", err)
					}
					return []reflect.Value{reflect.Zero(returns)}
				}),
			)
		case reflect.Struct:
			function := reflect.FuncOf([]reflect.Type{}, []reflect.Type{returns}, false)
			reflect.ValueOf(value).Elem().Field(index).Set(
				reflect.MakeFunc(function, func(args []reflect.Value) []reflect.Value {
					v := reflect.New(returns).Interface()
					err := Bind(v, object.Get(tag))
					if nil == err {
						return []reflect.Value{reflect.ValueOf(v)}
					} else {
						fmt.Println("Error: ", err)
					}
					return []reflect.Value{reflect.Zero(returns)}
				}),
			)
		default:
			return fmt.Errorf("property %s, unsupported kind %s", tag, returns.Kind())
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
