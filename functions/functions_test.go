package functions

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStructMethod int

func (t *testStructMethod) ReturnBool() bool {
	return true
}

func (t *testStructMethod) ReturnInt(a int) int {
	return 18 + a
}

func (t *testStructMethod) ReturnString(a, b string) string {
	return a + b
}

func (t *testStructMethod) ReturnError() error {
	return errors.New("functions: test error")
}

func (t *testStructMethod) unexportedMethod() bool {
	return true
}

func TestStructMethod(t *testing.T) {
	t.Parallel()
	a := testStructMethod(1)
	// ReturnBool
	val, err := StructMethod(&a, "ReturnBool")
	assert.NoError(t, err)
	valBool := val.(bool)
	assert.Equal(t, true, valBool)
	// ReturnInt
	val, err = StructMethod(&a, "ReturnInt", 2)
	assert.NoError(t, err)
	valInt := val.(int)
	assert.Equal(t, 20, valInt)
	// ReturnString
	val, err = StructMethod(&a, "ReturnString", "Santiago", " Rodriguez")
	assert.NoError(t, err)
	valString := val.(string)
	assert.Equal(t, "Santiago Rodriguez", valString)
	// ReturnError
	_, err = StructMethod(&a, "ReturnError")
	assert.Error(t, err)
	assert.Equal(t, "functions: test error", err.Error(), "expect same error msg")

	testErrors := map[string]struct {
		FnArgs []interface{}
	}{
		"empty fn name":            {FnArgs: []interface{}{}},
		"invalid fn name type":     {FnArgs: []interface{}{true}},
		"unknown fn":               {FnArgs: []interface{}{"Unknown"}},
		"not exported fn":          {FnArgs: []interface{}{"unexportedMethod"}},
		"res value cant interface": {FnArgs: []interface{}{}},
	}

	for name, tcs := range testErrors {
		t.Run(name, func(t *testing.T) {
			_, err := StructMethod(&a, tcs.FnArgs...)
			assert.Error(t, err)
		})
	}
}

func TestStructMethodInvalidParameters(t *testing.T) {
	t.Parallel()
	a := testStructMethod(1)
	testPanics := map[string]struct {
		FnArgs []interface{}
	}{
		"func with few args":      {FnArgs: []interface{}{"ReturnInt"}},
		"func with too many args": {FnArgs: []interface{}{"ReturnInt", 1, true}},
		"func with nil values":    {FnArgs: []interface{}{"ReturnInt", nil}},
		"invalid func args types": {FnArgs: []interface{}{"ReturnInt", true}},
		"arg with zero value ":    {FnArgs: []interface{}{"ReturnInt", reflect.Value{}}},
	}

	for name, tcs := range testPanics {
		t.Run(name, func(t *testing.T) {
			assert.Panics(t, func() { StructMethod(&a, tcs.FnArgs...) }, "Should panic")
		})
	}
}

func TestStructProperty(t *testing.T) {
	s := &struct {
		propOne bool
		PropTwo string
		Age     int64
	}{
		true,
		"hi",
		25,
	}
	p2, err := StructProperty(s, "PropTwo")
	assert.NoError(t, err)
	p2Stringed := p2.(string)
	assert.Equal(t, s.PropTwo, p2Stringed)

	_, err = StructProperty(s, "propOne")
	assert.Error(t, err)

	age, err := StructProperty(s, "Age")
	assert.NoError(t, err)
	ageInt := age.(int64)
	assert.Equal(t, s.Age, ageInt)

	_, err = StructProperty(s, "unknownField")
	assert.Error(t, err)

	_, err = StructProperty(s, 123)
	assert.Error(t, err)

	_, err = StructProperty(s)
	assert.Error(t, err)
}

func TestPreProcessFnEmpty(t *testing.T) {
	t.Parallel()
	p := PreProcessFn{}
	assert.True(t, p.Empty(), "preprocessFnEmpty empty name")
	p.Name = "Hi"
	assert.False(t, p.Empty(), "preprocessFnEmpty not empty name")
}
