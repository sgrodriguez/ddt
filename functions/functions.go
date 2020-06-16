package functions

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

// StructMethod ...
func StructMethod(str interface{}, args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("structMethod InvalidDecisionTreeFuncArg")
	}
	methodName, ok := args[0].(string)
	if !ok {
		return nil, errors.New("structMethod InvalidDecisionTreeFuncArg")
	}
	inputs := make([]reflect.Value, len(args)-1)
	for i, j := range args {
		if i == 0 {
			continue
		}
		inputs[i-1] = reflect.ValueOf(j)
	}
	method := reflect.ValueOf(str).MethodByName(methodName)
	if !method.IsValid() {
		return nil, errors.New("structMethod InvalidDecisionTreeFuncArg")
	}
	res := method.Call(inputs)
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	for _, r := range res {
		if r.Type().Implements(errorInterface) {
			if !r.IsNil() && r.CanInterface() {
				return nil, r.Interface().(error)
			}
		}
	}
	if !res[0].CanInterface() {
		return nil, errors.New("structMethod can not access the value")
	}
	return res[0].Interface(), nil
}

// StructProperty ...
func StructProperty(str interface{}, args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("structProperty InvalidDecisionTreeFuncArg")
	}
	prop, ok := args[0].(string)
	if !ok {
		return nil, errors.New("structProperty InvalidDecisionTreeFuncArg")
	}
	r := reflect.ValueOf(str)
	val := reflect.Indirect(r).FieldByName(prop)
	if !val.IsValid() || !val.CanInterface() {
		return nil, errors.New("structProperty can not access the value")
	}
	return val.Interface(), nil
}
