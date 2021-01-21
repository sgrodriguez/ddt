package function

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
	return errors.New("function: test error")
}

func (t *testStructMethod) unexportedMethod() bool {
	return true
}

func TestStructMethod(t *testing.T) {
	t.Parallel()
	a := testStructMethod(1)
	t.Run("return bool", func(t *testing.T) {
		val, err := CallStructMethod(&a, "ReturnBool")
		assert.NoError(t, err)
		valBool := val.(bool)
		assert.Equal(t, true, valBool)
	})
	t.Run("return int", func(t *testing.T) {
		val, err := CallStructMethod(&a, "ReturnInt", 2)
		assert.NoError(t, err)
		valInt := val.(int)
		assert.Equal(t, 20, valInt)
	})
	t.Run("return string", func(t *testing.T) {
		val, err := CallStructMethod(&a, "ReturnString", "Santiago", " Rodriguez")
		assert.NoError(t, err)
		valString := val.(string)
		assert.Equal(t, "Santiago Rodriguez", valString)
	})
	t.Run("return error", func(t *testing.T) {
		_, err := CallStructMethod(&a, "ReturnError")
		assert.Error(t, err)
		assert.EqualError(t, err, "function: test error", "expect same error msg")
	})
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
			_, err := CallStructMethod(&a, tcs.FnArgs...)
			assert.Error(t, err)
		})
	}
	t.Run("nil input", func(t *testing.T) {
		_, err := CallStructMethod(nil, "ReturnError")
		assert.Error(t, err)
		assert.EqualError(t, err, "callStructMethod invalid methods args")
	})
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
			assert.Panics(t, func() { CallStructMethod(&a, tcs.FnArgs...) }, "Should panic")
		})
	}
}

func TestStructAttribute(t *testing.T) {
	s := &struct {
		propOne bool
		PropTwo string
		Age     int64
	}{
		true,
		"hi",
		25,
	}

	t.Run("get struct attribute AttributeTwo", func(t *testing.T) {
		p2, err := GetStructAttribute(s, "PropTwo")
		assert.NoError(t, err)
		p2Stringed := p2.(string)
		assert.Equal(t, s.PropTwo, p2Stringed)

	})

	t.Run("get struct attribute propOne", func(t *testing.T) {
		_, err := GetStructAttribute(s, "propOne")
		assert.Error(t, err)

	})

	t.Run("get struct attribute Age", func(t *testing.T) {
		age, err := GetStructAttribute(s, "Age")
		assert.NoError(t, err)
		ageInt := age.(int64)
		assert.Equal(t, s.Age, ageInt)

	})

	t.Run("invalid args", func(t *testing.T) {
		_, err := GetStructAttribute(s, "unknownField")
		assert.Error(t, err)

		_, err = GetStructAttribute(s, 123)
		assert.Error(t, err)

		_, err = GetStructAttribute(s)
		assert.Error(t, err)

	})
}

func TestPreProcessFnEmpty(t *testing.T) {
	t.Parallel()
	p := PreProcessFn{}
	assert.True(t, p.Empty(), "preprocessFnEmpty empty name")
	p.Name = "Hi"
	assert.False(t, p.Empty(), "preprocessFnEmpty not empty name")
}
