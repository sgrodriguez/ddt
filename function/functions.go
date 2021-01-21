package function

import (
	"errors"
	"reflect"
)

// PreProcessFn ..
type PreProcessFn struct {
	Function func(str interface{}, args ...interface{}) (interface{}, error)
	Name     string
}

// Empty ..
func (p PreProcessFn) Empty() bool {
	return p.Name == ""
}

// CallStructMethod ... explain better what is
func CallStructMethod(str interface{}, args ...interface{}) (interface{}, error) {
	if len(args) < 1 || str == nil {
		return nil, errors.New("callStructMethod invalid methods args")
	}
	methodName, ok := args[0].(string)
	if !ok {
		return nil, errors.New("callStructMethod invalid methods args")
	}
	inputs := make([]reflect.Value, len(args)-1)
	// avoid first value of args
	for i, j := range args[1:] {
		inputs[i] = reflect.ValueOf(j)
	}
	method := reflect.ValueOf(str).MethodByName(methodName)
	if !method.IsValid() {
		return nil, errors.New("callStructMethod invalid methods args")
	}
	res := method.Call(inputs)
	if len(res) == 0 {
		return nil, errors.New("callStructMethod empty result from method call")
	}
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	for _, r := range res {
		if r.Type().Implements(errorInterface) {
			if !r.IsNil() && r.CanInterface() {
				return nil, r.Interface().(error)
			}
		}
	}
	if !res[0].CanInterface() {
		return nil, errors.New("callStructMethod can not interface result value")
	}
	return res[0].Interface(), nil
}

// GetStructAttribute ...  explain better what is
func GetStructAttribute(str interface{}, args ...interface{}) (interface{}, error) {
	if len(args) < 1 || str == nil {
		return nil, errors.New("getStructAttribute invalid args")
	}
	prop, ok := args[0].(string)
	if !ok {
		return nil, errors.New("getStructAttribute invalid args")
	}
	r := reflect.ValueOf(str)
	val := reflect.Indirect(r).FieldByName(prop)
	if !val.IsValid() || !val.CanInterface() {
		return nil, errors.New("getStructAttribute invalid struct attribute")
	}
	return val.Interface(), nil
}
