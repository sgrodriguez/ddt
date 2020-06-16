package ddt

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/sgrodriguez/ddt/comparators"
	"github.com/sgrodriguez/ddt/value"
)

func TestNodeMarshaller(t *testing.T) {
	node2 := &Node{
		Childes:        []*Node{},
		ID:             1,
		ParentID:       0,
		PreProcessArgs: []*value.Value{{Type: value.String, Value: "arg"}, {Type: value.Int64, Value: 123}},
		ValueToCompare: &value.Value{Type: value.String, Value: "marta"},
		Result:         &value.Value{Type: value.String, Value: "pedro"},
		Comparer:       &comparators.Greater{Equal: true},
	}
	node1 := &Node{
		Childes:        []*Node{node2},
		ID:             0,
		ParentID:       -1,
		PreProcessArgs: []*value.Value{{Type: value.String, Value: "arg"}, {Type: value.Int, Value: 12}},
		ValueToCompare: &value.Value{Type: value.String, Value: "lucia"},
		Result:         &value.Value{Type: value.String, Value: "lucia"},
		Comparer:       &comparators.Greater{Equal: false},
	}
	tr := &Tree{
		Root: node1,
		Name: "test tree",
	}
	val, err := json.Marshal(tr)
	AssertNotError(t, err)
	fmt.Println(string(val))
	var treeFromJSON *Tree
	err = json.Unmarshal(val, &treeFromJSON)
	AssertNotError(t, err)
	val, err = json.Marshal(treeFromJSON)
	AssertNotError(t, err)
	fmt.Println(string(val))
}

// AssertNotError fails the test if err is not null
func AssertNotError(t *testing.T, err error) {
	if err != nil {
		t.Errorf(`Unexpected error %s`, err)
		panic(err.Error())
	}
}

// AssertEqual fails the test if v1 and v2 are not equal
func AssertEqual(t *testing.T, label string, v1, v2 interface{}) {
	if !reflect.DeepEqual(v1, v2) {
		t.Errorf("%s: %T %#v is not equal to %T %#v\n", label, v1, v1, v2, v2)
		panic(label)
	}
}
