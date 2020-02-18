package utils

import (
	"fmt"
	"reflect"
	"strings"
	"syscall/js"
)

// Adapted from https://github.com/nlepage/golang-wasm/

type JSObject interface {
	JSValue() js.Value
	SetValue(js.Value)
}

type Node struct {
	js.Value
}

func (obj Node) JSValue() js.Value {
	return obj.Value
}

func (obj *Node) SetValue(value js.Value) {
	obj.Value = value
}

type Global struct {
	Node
	Window func() *Window `wasm:"window"`
}

func (obj Global) JSValue() js.Value {
	return obj.Value
}

func (obj *Global) SetValue(value js.Value) {
	obj.Value = value
}

type Window struct {
	Node
	Document func() *Document  `wasm:"document"`
	Location func() *Location  `wasm:"location"`
	Alert    func(interface{}) `wasm:"alert()"`
}

func (obj Window) JSValue() js.Value {
	return obj.Value
}

func (obj *Window) SetValue(value js.Value) {
	obj.Value = value
}

type Location struct {
	Node
	HRef     func() string `wasm:"href"`
	Origin   func() string `wasm:"origin"`
	Protocol func() string `wasm:"protocol"`
	Host     func() string `wasm:"host"`
	Hostname func() string `wasm:"hostname"`
	Port     func() string `wasm:"port"`
	Pathname func() string `wasm:"pathname"`
	Search   func() string `wasm:"search"`
	Hash     func() string `wasm:"hash"`
	Assign   func(string)  `wasm:"assign()"`
	Replace  func(string)  `wasm:"replace()"`
	Reload   func(bool)    `wasm:"reload()"`
}

func (obj Location) JSValue() js.Value {
	return obj.Value
}

func (obj *Location) SetValue(value js.Value) {
	obj.Value = value
}

type Document struct {
	Node
	Body          func() *HtmlElement       `wasm:"body"`
	CreateElement func(string) *HtmlElement `wasm:"createElement()"`
}

func (obj Document) JSValue() js.Value {
	return obj.Value
}

func (obj *Document) SetValue(value js.Value) {
	obj.Value = value
}

type HtmlElement struct {
	Node
	InnerHtml    func() string      `wasm:"innerHTML"`
	SetInnerHtml func(string)       `wasm:"innerHTML"`
	AppendChild  func(*HtmlElement) `wasm:"appendChild()"`
}

func (obj HtmlElement) JSValue() js.Value {
	return obj.Value
}

func (obj *HtmlElement) SetValue(value js.Value) {
	obj.Value = value
}

type HtmlCollection struct {
	Node
	Length    func() int             `wasm:"length"`
	Item      func(int) *HtmlElement `wasm:"item"`
	NamedItem func(int) *HtmlElement `wasm:"namedItem"`
}

func (obj HtmlCollection) JSValue() js.Value {
	return obj.Value
}

func (obj *HtmlCollection) SetValue(value js.Value) {
	obj.Value = value
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

func input(args []reflect.Value) []interface{} {
	input := make([]interface{}, 0)
	for _, v := range args {
		input = append(input, v.Interface())
	}
	return input
}

func BindFunction(tag string, fptr reflect.Type, object js.Value, value interface{}, index int) error {
	var (
		inputs  = fptr.NumIn()
		outputs = fptr.NumOut()
	)
	if tag, truth := IsFunction(tag); truth {
		if outputs > 1 {
			return fmt.Errorf("property %s accessed with more than 1 output parameters, expected function with 0/1 outputs", tag)
		}
		if outputs == 0 {
			parameters := make([]reflect.Type, 0)
			for i := 0; i < inputs; i++ {
				parameters = append(parameters, fptr.In(i))
			}
			function := reflect.FuncOf(parameters, []reflect.Type{}, fptr.IsVariadic())
			reflect.ValueOf(value).Elem().Field(index).Set(
				reflect.MakeFunc(function, func(args []reflect.Value) []reflect.Value {
					object.Call(tag, input(args)...)
					return []reflect.Value{}
				}),
			)
		}
		if outputs == 1 {
			parameters := make([]reflect.Type, 0)
			for i := 0; i < inputs; i++ {
				parameters = append(parameters, fptr.In(i))
			}
			returns := fptr.Out(0)
			switch returns.Kind() {
			case reflect.String:
				function := reflect.FuncOf(parameters, []reflect.Type{returns}, fptr.IsVariadic())
				reflect.ValueOf(value).Elem().Field(index).Set(
					reflect.MakeFunc(function, func(args []reflect.Value) []reflect.Value {
						return []reflect.Value{
							reflect.ValueOf(object.Call(tag, input(args)...).String()),
						}
					}),
				)
			case reflect.Int:
				function := reflect.FuncOf(parameters, []reflect.Type{returns}, fptr.IsVariadic())
				reflect.ValueOf(value).Elem().Field(index).Set(
					reflect.MakeFunc(function, func(args []reflect.Value) []reflect.Value {
						return []reflect.Value{
							reflect.ValueOf(object.Call(tag, input(args)...).Int()),
						}
					}),
				)
			case reflect.Float64:
				function := reflect.FuncOf(parameters, []reflect.Type{returns}, fptr.IsVariadic())
				reflect.ValueOf(value).Elem().Field(index).Set(
					reflect.MakeFunc(function, func(args []reflect.Value) []reflect.Value {
						return []reflect.Value{
							reflect.ValueOf(object.Call(tag, input(args)...).Float()),
						}
					}),
				)
			case reflect.Bool:
				function := reflect.FuncOf(parameters, []reflect.Type{returns}, fptr.IsVariadic())
				reflect.ValueOf(value).Elem().Field(index).Set(
					reflect.MakeFunc(function, func(args []reflect.Value) []reflect.Value {
						return []reflect.Value{
							reflect.ValueOf(object.Call(tag, input(args)...).Bool()),
						}
					}),
				)
			case reflect.Ptr:
				if kind := returns.Elem().Kind(); kind != reflect.Struct {
					return fmt.Errorf("property %s, unsupported kind %s", tag, returns.Kind())
				}
				if !returns.Elem().Implements(reflect.TypeOf((*js.Wrapper)(nil)).Elem()) {
					return fmt.Errorf("property %s, unsupported kind %s, must implement js.Wrapper", tag, returns.Kind())
				}
				function := reflect.FuncOf(parameters, []reflect.Type{returns}, fptr.IsVariadic())
				reflect.ValueOf(value).Elem().Field(index).Set(
					reflect.MakeFunc(function, func(args []reflect.Value) []reflect.Value {
						v := reflect.New(returns.Elem()).Interface()
						err := Bind(v, object.Call(tag, input(args)...))
						if nil == err {
							return []reflect.Value{reflect.ValueOf(v)}
						} else {
							fmt.Println("Error: ", err)
						}
						return []reflect.Value{reflect.Zero(returns)}
					}),
				)
			case reflect.Struct:
				function := reflect.FuncOf(parameters, []reflect.Type{returns}, fptr.IsVariadic())
				reflect.ValueOf(value).Elem().Field(index).Set(
					reflect.MakeFunc(function, func(args []reflect.Value) []reflect.Value {
						v := reflect.New(returns).Interface()
						err := Bind(v, object.Call(tag, input(args)...))
						if nil == err {
							return []reflect.Value{reflect.ValueOf(v)}
						} else {
							fmt.Println("Error: ", err)
						}
						return []reflect.Value{reflect.Zero(returns)}
					}),
				)
			default:
				return fmt.Errorf("function %s, returns unsupported kind %s", tag, returns.Kind())
			}
		}
	}
	if tag, truth := IsProperty(tag); truth {
		if inputs > 1 {
			return fmt.Errorf("property %s accessed with more than 1 input parameters, expected function with 0/1 inputs", tag)
		}
		if outputs > 1 {
			return fmt.Errorf("property %s accessed with no output parameters, expected function with 0/1 outputs", tag)
		}
		if outputs == 0 {
			if inputs != 1 {
				return fmt.Errorf("property %s accessed with zero or more than 1 input parameters, expected function with 1 inputs", tag)
			}
			parameters := make([]reflect.Type, 0)
			for i := 0; i < inputs; i++ {
				parameters = append(parameters, fptr.In(i))
			}
			function := reflect.FuncOf(parameters, []reflect.Type{}, false)
			reflect.ValueOf(value).Elem().Field(index).Set(
				reflect.MakeFunc(function, func(args []reflect.Value) []reflect.Value {
					object.Set(tag, input(args)[0])
					return []reflect.Value{}
				}),
			)
		}
		if outputs > 0 {
			if inputs != 0 {
				return fmt.Errorf("property %s accessed with input parameters, expected function with no inputs", tag)
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
				if !returns.Elem().Implements(reflect.TypeOf((*js.Wrapper)(nil)).Elem()) {
					return fmt.Errorf("property %s, unsupported kind %s, must implement js.Wrapper", tag, returns.Kind())
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
	}
	return nil
}

func MatchInputs(method reflect.Method, inputs []reflect.Type) error {
	for i, value := range inputs {
		if method.Type.In(i) != value {
			return fmt.Errorf("type mismatch, expected %s got %s", method.Type.In(i), value)
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
	for _, method := range Functions(v) {
		if strings.Compare("SetValue", method.Name) == 0 {
			if method.Type.NumIn() == 2 && method.Type.NumOut() == 0 {
				inputs := []reflect.Type{
					reflect.TypeOf(v),
					reflect.TypeOf(object),
				}
				if err := MatchInputs(method, inputs); nil == err {
					value := reflect.ValueOf(v).MethodByName(method.Name)
					value.Call([]reflect.Value{
						reflect.ValueOf(object),
					})
					return nil
				}
			}
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

func Functions(v interface{}) []reflect.Method {
	var (
		elem    = reflect.ValueOf(v).Type()
		count   = elem.NumMethod()
		methods = make([]reflect.Method, 0)
	)
	for i := 0; i < count; i++ {
		methods = append(methods, elem.Method(i))
	}
	return methods
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
