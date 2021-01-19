package value

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValue(t *testing.T) {
	testCases := map[string]struct {
		valType       Type
		concreteValue interface{}
		expectedValue Value
	}{
		"int64 valid value": {
			valType:       Int64,
			concreteValue: int64(10),
			expectedValue: Value{Value: int64(10), Type: Int64},
		},
		"bool valid value": {
			valType:       Bool,
			concreteValue: true,
			expectedValue: Value{Value: true, Type: Bool},
		},
		"int valid value": {
			valType:       Int,
			concreteValue: 10,
			expectedValue: Value{Value: 10, Type: Int},
		},
		"float64 valid value": {
			valType:       Float64,
			concreteValue: 4.20,
			expectedValue: Value{Value: 4.20, Type: Float64},
		},
		"string valid value": {
			valType:       String,
			concreteValue: "xd",
			expectedValue: Value{Value: "xd", Type: String},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			val, err := NewValue(tc.valType, tc.concreteValue)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedValue.Value, val.Value)
			assert.Equal(t, tc.expectedValue.Type, val.Type)
		})
	}
	invalidValueTestCases := map[string]struct {
		valType       Type
		concreteValue interface{}
	}{
		"int64 invalid int concrete value": {
			valType:       Int64,
			concreteValue: 10,
		},
		"int64 invalid bool concrete value": {
			valType:       Int64,
			concreteValue: true,
		},
		"int invalid int64 concrete value": {
			valType:       Int,
			concreteValue: int64(10),
		},
		"int invalid bool concrete value": {
			valType:       Int,
			concreteValue: false,
		},
		"bool invalid concrete value": {
			valType:       Bool,
			concreteValue: 0,
		},
		"float64 invalid concrete value": {
			valType:       Float64,
			concreteValue: false,
		},
		"string invalid concrete value": {
			valType:       String,
			concreteValue: 10,
		},
	}
	for name, tc := range invalidValueTestCases {
		t.Run(name, func(t *testing.T) {
			_, err := NewValue(tc.valType, tc.concreteValue)
			require.Error(t, err)
		})
	}
}

func TestValueMarshal(t *testing.T) {
	tests := map[string]struct {
		input    Value
		expected string
	}{
		"marshal int64":        {input: Value{Type: Int64, Value: 12}, expected: `{"Value":12,"Type":"int64"}`},
		"marshal string":       {input: Value{Type: String, Value: "hi all"}, expected: `{"Value":"hi all","Type":"string"}`},
		"marshal string empty": {input: Value{Type: String, Value: ""}, expected: `{"Value":"","Type":"string"}`},
		"marshal int":          {input: Value{Type: Int, Value: 42}, expected: `{"Value":42,"Type":"int"}`},
		"marshal bool":         {input: Value{Type: Bool, Value: true}, expected: `{"Value":true,"Type":"bool"}`},
		"marshal float64":      {input: Value{Type: Float64, Value: 3.14}, expected: `{"Value":3.14,"Type":"float64"}`},
	}
	for name, tst := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := json.Marshal(&tst.input)
			assert.NoError(t, err, "expected no error in marshaling")
			assert.JSONEq(t, tst.expected, string(got), "expected value marshaled differs %s", name)
		})
	}
}

func TestValueUnmarshal(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected Value
	}{
		"unmarshal int64":        {expected: Value{Type: Int64, Value: int64(12)}, input: `{"Value":12,"Type":"int64"}`},
		"unmarshal string":       {expected: Value{Type: String, Value: "hi all"}, input: `{"Value":"hi all","Type":"string"}`},
		"unmarshal string empty": {expected: Value{Type: String, Value: ""}, input: `{"Value":"","Type":"string"}`},
		"unmarshal int":          {expected: Value{Type: Int, Value: 42}, input: `{"Value":42,"Type":"int"}`},
		"unmarshal bool":         {expected: Value{Type: Bool, Value: true}, input: `{"Value":true,"Type":"bool"}`},
		"unmarshal float64":      {expected: Value{Type: Float64, Value: 3.14}, input: `{"Value":3.14,"Type":"float64"}`},
	}
	for name, tst := range tests {
		t.Run(name, func(t *testing.T) {
			var v Value
			err := json.Unmarshal([]byte(tst.input), &v)
			assert.NoError(t, err, "expected no error in unmarshalling")
			assert.Equal(t, tst.expected, v, "expected value unmarshalled differs %s", name)
		})
	}
	testsError := map[string]struct {
		input string
	}{
		"unmarshal int64 string value":   {input: `{"Value":"12213","Type":"int64"}`},
		"unmarshal string bool value":    {input: `{"Value":true,"Type":"string"}`},
		"unmarshal bool int value":       {input: `{"Value":12,"Type":"bool"}`},
		"unmarshal string float value":   {input: `{"Value":123.12,"Type":"string"}`},
		"unmarshal int float value":      {input: `{"Value":123.12,"Type":"int"}`},
		"unmarshal float64 string value": {input: `{"Value":"float64","Type":"float64"}`},
		"unmarshal unknown type ":        {input: `{"Value":123.12,"Type":"unknownType"}`},
		"unmarshal invalid type json ":   {input: `{"Value":"","Type":[1,2,3]}`},
	}
	for name, tst := range testsError {
		t.Run(name, func(t *testing.T) {
			var v Value
			err := json.Unmarshal([]byte(tst.input), &v)
			assert.Error(t, err, "expected error in unmarshalling")
		})
	}
}

func TestGetValueInterfaces(t *testing.T) {
	var values []*Value
	valuesInterfaces := GetValueInterfaces(values)
	assert.Equal(t, 0, len(valuesInterfaces), "expect empty slice")
	values = append(values, &Value{Value: 1, Type: Int})
	values = append(values, &Value{Value: "asd", Type: String})
	values = append(values, &Value{Value: true, Type: Bool})
	valuesInterfaces = GetValueInterfaces(values)
	assert.Equal(t, 3, len(valuesInterfaces), "expect len of slice equal 3")
	for _, v := range values {
		assert.Contains(t, valuesInterfaces, v.Value, "expect to contain value")
	}
}
