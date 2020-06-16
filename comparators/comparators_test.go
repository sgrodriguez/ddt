package comparators

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualCompare(t *testing.T) {
	tests := map[string]struct {
		inputA   interface{}
		inputB   interface{}
		expected bool
	}{
		"equal string":       {"hi", "hi", true},
		"equal empty string": {"", "", true},
		"not equal string":   {"hi", "", false},
		"equal bool":         {true, true, true},
		"not equal bool":     {true, false, false},
		"equal int":          {1, 1, true},
		"not equal int":      {2, 0, false},
		"equal int64":        {int64(1), int64(1), true},
		"not equal int64":    {int64(0), int64(1), false},
		"different types":    {"hi", 1, false},
		"equal float64":      {3.14, 3.14, true},
		"not equal float64":  {3.14, 2.71, false},
	}
	eq := Equal{}
	for name, tcs := range tests {
		t.Run(name, func(t *testing.T) {
			got := eq.Compare(tcs.inputA, tcs.inputB)
			assert.Equal(t, tcs.expected, got)
		})
	}
}

func TestGreaterCompare(t *testing.T) {
	type testGreater struct {
		inputA   interface{}
		inputB   interface{}
		expected bool
	}
	testsGreater := map[string]testGreater{
		"greater int":          {2, 1, true},
		"not greater int":      {2, 20, false},
		"greater int64":        {int64(2), int64(1), true},
		"not greater int64":    {int64(0), int64(1), false},
		"different types":      {"hi", 1, false},
		"greater float64":      {4.14, 3.14, true},
		"not grater float64":   {1.14, 2.71, false},
		"not comparable types": {"asd", "asd", false},
	}
	gt := Greater{}
	for name, tcs := range testsGreater {
		t.Run(name, func(t *testing.T) {
			got := gt.Compare(tcs.inputA, tcs.inputB)
			assert.Equal(t, tcs.expected, got)
		})
	}
	testsGreaterOrEqual := map[string]testGreater{
		"greater or equal int":     {2, 2, true},
		"greater or equal int64":   {int64(1), int64(1), true},
		"greater or equal float64": {3.14, 3.14, true},
	}
	gt.Equal = true
	for name, tcs := range testsGreaterOrEqual {
		t.Run(name, func(t *testing.T) {
			got := gt.Compare(tcs.inputA, tcs.inputB)
			assert.Equal(t, tcs.expected, got)
		})
	}

}

func TestLesserCompare(t *testing.T) {
	type testLesser struct {
		inputA   interface{}
		inputB   interface{}
		expected bool
	}
	testsLesser := map[string]testLesser{
		"lesser int":           {2, 1, false},
		"not lesser int":       {2, 20, true},
		"lesser int64":         {int64(2), int64(1), false},
		"not lesser int64":     {int64(0), int64(1), true},
		"different types":      {"hi", 1, false},
		"lesser float64":       {4.14, 3.14, false},
		"not lesser float64":   {1.14, 2.71, true},
		"not comparable types": {"asd", "asd", false},
	}
	lt := Lesser{}
	for name, tcs := range testsLesser {
		t.Run(name, func(t *testing.T) {
			got := lt.Compare(tcs.inputA, tcs.inputB)
			assert.Equal(t, tcs.expected, got)
		})
	}
	testsLesserOrEqual := map[string]testLesser{
		"lesser or equal int":     {2, 2, true},
		"lesser or equal int64":   {int64(1), int64(1), true},
		"lesser or equal float64": {3.14, 3.14, true},
	}
	lt.Equal = true
	for name, tcs := range testsLesserOrEqual {
		t.Run(name, func(t *testing.T) {
			got := lt.Compare(tcs.inputA, tcs.inputB)
			assert.Equal(t, tcs.expected, got)
		})
	}
}

func TestMarshalComparator(t *testing.T) {
	tests := map[string]struct {
		input    Comparer
		expected string
	}{
		"marshal equal":            {input: &Equal{}, expected: `{"type":"eq"}`},
		"marshal greater":          {input: &Greater{}, expected: `{"equal":false, "type":"gt"}`},
		"marshal greater or equal": {input: &Greater{true}, expected: `{"equal":true, "type":"gt"}`},
		"marshal lesser":           {input: &Lesser{}, expected: `{"equal":false, "type":"lt"}`},
		"marshal lesser or equal":  {input: &Lesser{true}, expected: `{"equal":true, "type":"lt"}`},
	}
	for name, tst := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := json.Marshal(&tst.input)
			assert.NoError(t, err, "expected no error in marshaling")
			assert.JSONEq(t, tst.expected, string(got), "expected comparator marshaled not differs %s", name)
		})
	}
}

func TestUnmarshallComparator(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected Comparer
	}{
		"unmarshal equal":            {expected: &Equal{}, input: `{"type":"eq"}`},
		"unmarshal greater":          {expected: &Greater{}, input: `{"equal":false, "type":"gt"}`},
		"unmarshal greater or equal": {expected: &Greater{true}, input: `{"equal":true, "type":"gt"}`},
		"unmarshal lesser":           {expected: &Lesser{}, input: `{"equal":false, "type":"lt"}`},
		"unmarshal lesser or equal":  {expected: &Lesser{true}, input: `{"equal":true, "type":"lt"}`},
	}
	for name, tst := range tests {
		t.Run(name, func(t *testing.T) {
			comp, err := CreateComparatorFromJSON(json.RawMessage(tst.input))
			assert.NoError(t, err, "expected no error in unmarshalling")
			assert.Equal(t, tst.expected, comp, "expected value unmarshalled differs %s", name)
		})
	}
	testsError := map[string]struct {
		input string
	}{
		"unmarshal unknown type":      {input: `{"type":"asd"}`},
		"unmarshal invalid json type": {input: `{"equal":"asd"}`},
	}
	for name, tst := range testsError {
		t.Run(name, func(t *testing.T) {
			_, err := CreateComparatorFromJSON(json.RawMessage(tst.input))
			assert.Error(t, err, "expected error in unmarshalling")
		})
	}
}
