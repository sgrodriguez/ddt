package ddt

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sgrodriguez/ddt/compare"
	"github.com/sgrodriguez/ddt/function"
	"github.com/sgrodriguez/ddt/value"
)

func TestResolveTree_SimpleTree(t *testing.T) {
	leaf1 := Node{
		ID:             1,
		ParentID:       0,
		ValueToCompare: &value.Value{Value: int64(60), Type: value.Int64},
		Comparer:       &compare.Greater{},
		Result:         &value.Value{Value: "prize1", Type: value.String},
	}
	leaf11 := Node{
		ID:             3,
		ParentID:       2,
		ValueToCompare: &value.Value{Value: int64(30), Type: value.Int64},
		Comparer:       &compare.Equal{},
		Result:         &value.Value{Value: "prize2", Type: value.String},
	}
	leaf12 := Node{
		ID:             4,
		ParentID:       2,
		ValueToCompare: &value.Value{Value: int64(30), Type: value.Int64},
		Comparer:       &compare.Greater{},
		Result:         &value.Value{Value: "prize3", Type: value.String},
	}
	leaf13 := Node{
		ID:             5,
		ParentID:       2,
		ValueToCompare: &value.Value{Value: int64(30), Type: value.Int64},
		Comparer:       &compare.Lesser{},
		Result:         &value.Value{Value: "prize4", Type: value.String},
	}
	node1 := Node{
		Childes:        []*Node{&leaf11, &leaf12, &leaf13},
		ID:             2,
		ParentID:       0,
		ValueToCompare: &value.Value{Value: int64(60), Type: value.Int64},
		Comparer:       &compare.Lesser{},
	}
	root := Node{
		ID:       0,
		ParentID: -1,
		Childes:  []*Node{&node1, &leaf1},
	}
	simpleTree, err := NewTree("simpleTree", &root)
	require.NoError(t, err)
	testCases := map[string]struct {
		input    int64
		expected string
	}{
		"given 100 as input expect prize1": {input: 100, expected: "prize1"},
		"given 45 as input expect prize3":  {input: 45, expected: "prize3"},
		"given 30 as input expect prize2":  {input: 30, expected: "prize2"},
		"given 10 as input expect prize4":  {input: 10, expected: "prize4"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actual, err := ResolveTree(simpleTree, tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual.(string))
		})
	}
	_, err = ResolveTree(simpleTree, "invalid input")
	require.Error(t, err)
	_, err = ResolveTree(simpleTree, nil)
	require.Error(t, err)
	_, err = ResolveTree(simpleTree, 4.20)
	require.Error(t, err)
}

type user struct {
	Age       int
	FirstName string
	LastName  string
}

func (u *user) UnderAge() bool {
	return u.Age < 18
}

func (u *user) FullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *user) ReturnErr() error {
	return errors.New("always error")
}

func newUser(age int, firstName, lastName string) *user {
	return &user{
		Age:       age,
		FirstName: firstName,
		LastName:  lastName,
	}
}

func userTree() *Tree {
	node6 := &Node{
		ID:             6,
		ParentID:       2,
		ValueToCompare: &value.Value{Type: value.Int, Value: 30},
		Result:         &value.Value{Type: value.String, Value: "node6"},
		Comparer:       &compare.Greater{},
	}
	node5 := &Node{
		ID:             5,
		ParentID:       2,
		ValueToCompare: &value.Value{Type: value.Int, Value: 30},
		Result:         &value.Value{Type: value.String, Value: "node5"},
		Comparer:       &compare.Lesser{Equal: true},
	}
	node3 := &Node{
		ID:             3,
		ParentID:       1,
		ValueToCompare: &value.Value{Type: value.String, Value: "SANTIAGO LUCIA"},
		Result:         &value.Value{Type: value.String, Value: "node3"},
		Comparer:       &compare.Equal{},
	}
	node4 := &Node{
		ID:             4,
		ParentID:       1,
		ValueToCompare: &value.Value{Type: value.String, Value: "LUCIA SANTIAGO"},
		Result:         &value.Value{Type: value.String, Value: "node4"},
		Comparer:       &compare.Equal{},
	}
	node1 := &Node{
		ID:             1,
		ParentID:       0,
		Childes:        []*Node{node3, node4},
		ValueToCompare: &value.Value{Type: value.Bool, Value: true},
		PreProcessArgs: []*value.Value{{Type: value.String, Value: "FullName"}},
		PreProcessFn:   function.PreProcessFn{Function: function.StructMethod, Name: "StructMethod"},
		Comparer:       &compare.Equal{},
	}
	node2 := &Node{
		ID:             2,
		ParentID:       0,
		Childes:        []*Node{node5, node6},
		ValueToCompare: &value.Value{Type: value.Bool, Value: false},
		Comparer:       &compare.Equal{},
		PreProcessArgs: []*value.Value{{Type: value.String, Value: "Age"}},
		PreProcessFn:   function.PreProcessFn{Function: function.StructAttribute, Name: "StructAttribute"},
	}
	root := &Node{
		Childes:        []*Node{node1, node2},
		PreProcessArgs: []*value.Value{{Type: value.String, Value: "UnderAge"}},
		PreProcessFn:   function.PreProcessFn{Function: function.StructMethod, Name: "StructMethod"},
		ID:             0,
		ParentID:       -1,
	}
	tree, _ := NewTree("userTree", root)
	return tree
}

func TestResolveTree_UserTree(t *testing.T) {
	userTree := userTree()
	b, err := json.Marshal(userTree)
	require.NoError(t, err)

	treeFromJSON, err := NewTree("treeFromJSON", &Node{ID: 0, ParentID: -1})
	err = json.Unmarshal(b, treeFromJSON)
	require.NoError(t, err)

	testCases := map[string]struct {
		input    *user
		expected string
	}{
		"expected node3": {input: newUser(11, "SANTIAGO", "LUCIA"), expected: "node3"},
		"expected node4": {input: newUser(11, "LUCIA", "SANTIAGO"), expected: "node4"},
		"expected node5": {input: newUser(25, "LUCIA", "SANTIAGO"), expected: "node5"},
		"expected node6": {input: newUser(65, "LUCIA", "SANTIAGO"), expected: "node6"},
	}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			actual, err := ResolveTree(userTree, tc.input)
			require.NoError(t, err)
			actualFromNewTree, err := ResolveTree(treeFromJSON, tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual.(string))
			assert.Equal(t, tc.expected, actualFromNewTree.(string))
		})
	}
	//
	t.Run("handle empty/invalid user", func(t *testing.T) {
		_, err = ResolveTree(userTree, &user{})
		require.Error(t, err)
		assert.EqualError(t, err, "value not found when comparing with all childes nodes")
		_, err = ResolveTree(userTree, nil)
		require.Error(t, err)
		assert.EqualError(t, err, "structMethod invalid methods args")
		_, err = ResolveTree(userTree, 123)
		require.Error(t, err)
		assert.EqualError(t, err, "structMethod invalid methods args")

	})
	t.Run("handling method error", func(t *testing.T) {
		node1 := userTree.Root.Childes[0]
		node1.PreProcessArgs = []*value.Value{{Type: value.String, Value: "ReturnErr"}}
		_, err = ResolveTree(userTree, &user{Age: 11})
		require.Error(t, err)
		assert.EqualError(t, err, "always error")
	})
	t.Run("invalid method name", func(t *testing.T) {
		node1 := userTree.Root.Childes[0]
		node1.PreProcessArgs = []*value.Value{{Type: value.String, Value: "ASD"}}
		_, err = ResolveTree(userTree, &user{Age: 11})
		assert.EqualError(t, err, "structMethod invalid methods args")
	})
	t.Run("invalid struct attribute name", func(t *testing.T) {
		node2 := userTree.Root.Childes[1]
		node2.PreProcessArgs = []*value.Value{{Type: value.String, Value: "ASD"}}
		_, err = ResolveTree(userTree, &user{Age: 30})
		assert.EqualError(t, err, "structAttribute invalid struct attribute")
	})
}

// TODO replace json []byte literal with some lib to manipulate better (ie simplejson)
func TestModifyTreeTroughJSON(t *testing.T) {
	t.Run("modified node4 result to 100 int", func(t *testing.T) {
		ut := userTree()
		valueBeforeModified, err := ResolveTree(ut, newUser(11, "LUCIA", "SANTIAGO"))
		require.NoError(t, err)
		assert.Equal(t, "node4", valueBeforeModified.(string))
		newTree := []byte(`{"nodes":[{"preProcessFnName":"StructMethod","id":0,"parentId":-1,"preProcessFnArgs":[{"Value":"UnderAge","Type":"string"}]},{"preProcessFnName":"StructMethod","id":1,"parentId":0,"preProcessFnArgs":[{"Value":"FullName","Type":"string"}],"comparer":{"type":"eq"},"valueToCompare":{"Value":true,"Type":"bool"}},{"preProcessFnName":"StructAttribute","id":2,"parentId":0,"preProcessFnArgs":[{"Value":"Age","Type":"string"}],"comparer":{"type":"eq"},"valueToCompare":{"Value":false,"Type":"bool"}},{"preProcessFnName":"","id":3,"parentId":1,"comparer":{"type":"eq"},"valueToCompare":{"Value":"SANTIAGO LUCIA","Type":"string"},"result":{"Value":"node3","Type":"string"}},{"preProcessFnName":"","id":4,"parentId":1,"comparer":{"type":"eq"},"valueToCompare":{"Value":"LUCIA SANTIAGO","Type":"string"},"result":{"Value":100,"Type":"int"}},{"preProcessFnName":"","id":5,"parentId":2,"comparer":{"type":"lt","equal":true},"valueToCompare":{"Value":30,"Type":"int"},"result":{"Value":"node5","Type":"string"}},{"preProcessFnName":"","id":6,"parentId":2,"comparer":{"type":"gt","equal":false},"valueToCompare":{"Value":30,"Type":"int"},"result":{"Value":"node6","Type":"string"}}],"name":"userTree"}`)
		err = json.Unmarshal(newTree, ut)
		require.NoError(t, err)
		actualResult, err := ResolveTree(ut, newUser(11, "LUCIA", "SANTIAGO"))
		require.NoError(t, err)
		assert.Equal(t, 100, actualResult.(int))
	})
	t.Run("modified under age node value to compare", func(t *testing.T) {
		ut := userTree()
		resultBeforeModified, err := ResolveTree(ut, newUser(33, "SANTIAGO", "LUCIA"))
		require.NoError(t, err)
		assert.Equal(t, "node6", resultBeforeModified.(string))
		newTree := []byte(`{"nodes":[{"preProcessFnName":"StructMethod","id":0,"parentId":-1,"preProcessFnArgs":[{"Value":"UnderAge","Type":"string"}]},{"preProcessFnName":"StructMethod","id":1,"parentId":0,"preProcessFnArgs":[{"Value":"FullName","Type":"string"}],"comparer":{"type":"eq"},"valueToCompare":{"Value":false,"Type":"bool"}},{"preProcessFnName":"StructAttribute","id":2,"parentId":0,"preProcessFnArgs":[{"Value":"Age","Type":"string"}],"comparer":{"type":"eq"},"valueToCompare":{"Value":true,"Type":"bool"}},{"preProcessFnName":"","id":3,"parentId":1,"comparer":{"type":"eq"},"valueToCompare":{"Value":"SANTIAGO LUCIA","Type":"string"},"result":{"Value":"node3","Type":"string"}},{"preProcessFnName":"","id":4,"parentId":1,"comparer":{"type":"eq"},"valueToCompare":{"Value":"LUCIA SANTIAGO","Type":"string"},"result":{"Value":"node4","Type":"string"}},{"preProcessFnName":"","id":5,"parentId":2,"comparer":{"type":"lt","equal":true},"valueToCompare":{"Value":30,"Type":"int"},"result":{"Value":"node5","Type":"string"}},{"preProcessFnName":"","id":6,"parentId":2,"comparer":{"type":"gt","equal":false},"valueToCompare":{"Value":30,"Type":"int"},"result":{"Value":"node6","Type":"string"}}],"name":"userTree"}`)
		err = json.Unmarshal(newTree, ut)
		require.NoError(t, err)
		actualResult, err := ResolveTree(ut, newUser(33, "SANTIAGO", "LUCIA"))
		require.NoError(t, err)
		assert.Equal(t, "node3", actualResult.(string))
	})
	t.Run("modified age compare from gt to lt", func(t *testing.T) {
		ut := userTree()
		resultBeforeModified, err := ResolveTree(ut, newUser(55, "SANTIAGO", "LUCIA"))
		require.NoError(t, err)
		assert.Equal(t, "node6", resultBeforeModified.(string))
		newTree := []byte(`{"nodes":[{"preProcessFnName":"StructMethod","id":0,"parentId":-1,"preProcessFnArgs":[{"Value":"UnderAge","Type":"string"}]},{"preProcessFnName":"StructMethod","id":1,"parentId":0,"preProcessFnArgs":[{"Value":"FullName","Type":"string"}],"comparer":{"type":"eq"},"valueToCompare":{"Value":true,"Type":"bool"}},{"preProcessFnName":"StructAttribute","id":2,"parentId":0,"preProcessFnArgs":[{"Value":"Age","Type":"string"}],"comparer":{"type":"eq"},"valueToCompare":{"Value":false,"Type":"bool"}},{"preProcessFnName":"","id":3,"parentId":1,"comparer":{"type":"eq"},"valueToCompare":{"Value":"SANTIAGO LUCIA","Type":"string"},"result":{"Value":"node3","Type":"string"}},{"preProcessFnName":"","id":4,"parentId":1,"comparer":{"type":"eq"},"valueToCompare":{"Value":"LUCIA SANTIAGO","Type":"string"},"result":{"Value":"node4","Type":"string"}},{"preProcessFnName":"","id":5,"parentId":2,"comparer":{"type":"gt","equal":false},"valueToCompare":{"Value":30,"Type":"int"},"result":{"Value":"node5","Type":"string"}},{"preProcessFnName":"","id":6,"parentId":2,"comparer":{"type":"lt","equal":true},"valueToCompare":{"Value":30,"Type":"int"},"result":{"Value":"node6","Type":"string"}}],"name":"userTree"}`)
		err = json.Unmarshal(newTree, ut)
		require.NoError(t, err)
		actualResult, err := ResolveTree(ut, newUser(55, "SANTIAGO", "LUCIA"))
		require.NoError(t, err)
		assert.Equal(t, "node5", actualResult.(string))
	})
	t.Run("modified Fullname pre process StructMethod fn to FirstName StructAttribute fn", func(t *testing.T) {
		ut := userTree()
		resultBeforeModified, err := ResolveTree(ut, newUser(11, "SANTIAGO", "LUCIA"))
		require.NoError(t, err)
		assert.Equal(t, "node3", resultBeforeModified.(string))
		resultBeforeModified, err = ResolveTree(ut, newUser(11, "LUCIA", "SANTIAGO"))
		require.NoError(t, err)
		assert.Equal(t, "node4", resultBeforeModified.(string))
		newTree := []byte(`{"nodes":[{"preProcessFnName":"StructMethod","id":0,"parentId":-1,"preProcessFnArgs":[{"Value":"UnderAge","Type":"string"}]},{"preProcessFnName":"StructAttribute","id":1,"parentId":0,"preProcessFnArgs":[{"Value":"FirstName","Type":"string"}],"comparer":{"type":"eq"},"valueToCompare":{"Value":true,"Type":"bool"}},{"preProcessFnName":"StructAttribute","id":2,"parentId":0,"preProcessFnArgs":[{"Value":"Age","Type":"string"}],"comparer":{"type":"eq"},"valueToCompare":{"Value":false,"Type":"bool"}},{"preProcessFnName":"","id":3,"parentId":1,"comparer":{"type":"eq"},"valueToCompare":{"Value":"SANTIAGO","Type":"string"},"result":{"Value":"node3","Type":"string"}},{"preProcessFnName":"","id":4,"parentId":1,"comparer":{"type":"eq"},"valueToCompare":{"Value":"LUCIA","Type":"string"},"result":{"Value":"node4","Type":"string"}},{"preProcessFnName":"","id":5,"parentId":2,"comparer":{"type":"lt","equal":true},"valueToCompare":{"Value":30,"Type":"int"},"result":{"Value":"node5","Type":"string"}},{"preProcessFnName":"","id":6,"parentId":2,"comparer":{"type":"gt","equal":false},"valueToCompare":{"Value":30,"Type":"int"},"result":{"Value":"node6","Type":"string"}}],"name":"userTree"}`)
		err = json.Unmarshal(newTree, ut)
		require.NoError(t, err)
		actualResult, err := ResolveTree(ut, newUser(11, "SANTIAGO", "xd"))
		require.NoError(t, err)
		assert.Equal(t, "node3", actualResult.(string))
		actualResult, err = ResolveTree(ut, newUser(11, "LUCIA", "xd"))
		require.NoError(t, err)
		assert.Equal(t, "node4", actualResult.(string))
	})
}

func TestNewTree_InvalidRootNode(t *testing.T) {
	roots := []*Node{
		{
			ParentID: 2,
		},
		{
			ParentID: -1,
			ID:       3,
		},
		{
			ParentID: -1,
			ID:       0,
			Result:   &value.Value{Value: 10, Type: value.Int},
		},
	}
	for _, r := range roots {
		_, err := NewTree("invalidTree", r)
		require.Error(t, err)
	}
}

func TestMarshalComparator(t *testing.T) {
	tests := map[string]struct {
		input    Comparer
		expected string
	}{
		"marshal equal":            {input: &compare.Equal{}, expected: `{"type":"eq"}`},
		"marshal greater":          {input: &compare.Greater{}, expected: `{"equal":false, "type":"gt"}`},
		"marshal greater or equal": {input: &compare.Greater{Equal: true}, expected: `{"equal":true, "type":"gt"}`},
		"marshal lesser":           {input: &compare.Lesser{}, expected: `{"equal":false, "type":"lt"}`},
		"marshal lesser or equal":  {input: &compare.Lesser{Equal: true}, expected: `{"equal":true, "type":"lt"}`},
	}
	for name, tst := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := json.Marshal(&tst.input)
			assert.NoError(t, err, "expected no error in marshaling")
			assert.JSONEq(t, tst.expected, string(got), "expected compare marshaled not differs %s", name)
		})
	}
}

func TestUnmarshallComparator(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected Comparer
	}{
		"unmarshal equal":            {expected: &compare.Equal{}, input: `{"type":"eq"}`},
		"unmarshal greater":          {expected: &compare.Greater{}, input: `{"equal":false, "type":"gt"}`},
		"unmarshal greater or equal": {expected: &compare.Greater{Equal: true}, input: `{"equal":true, "type":"gt"}`},
		"unmarshal lesser":           {expected: &compare.Lesser{}, input: `{"equal":false, "type":"lt"}`},
		"unmarshal lesser or equal":  {expected: &compare.Lesser{Equal: true}, input: `{"equal":true, "type":"lt"}`},
	}
	for name, tst := range tests {
		t.Run(name, func(t *testing.T) {
			comp, err := createComparatorFromJSON(json.RawMessage(tst.input))
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
			_, err := createComparatorFromJSON(json.RawMessage(tst.input))
			assert.Error(t, err, "expected error in unmarshalling")
		})
	}
}
