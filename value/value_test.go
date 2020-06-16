package value

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
