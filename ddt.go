package ddt

import (
	"encoding/json"
	"errors"

	"github.com/sgrodriguez/ddt/functions"
)

// DefaultFns default functions
var DefaultFns = []functions.PreProcessFn{
	{Function: functions.StructMethod, Name: "StructMethod"},
	{Function: functions.StructProperty, Name: "StructProperty"},
}

// Tree Type
type Tree struct {
	Root      *Node                             `json:"-"`
	Functions map[string]functions.PreProcessFn `json:"-"`
	Name      string                            `json:"name"`
}

// NewTree creates a tree
func NewTree(name string, fn []functions.PreProcessFn) *Tree {
	return &Tree{Name: name, Functions: addNewPreProcessFn(fn), Root: &Node{ID: 0, Result: nil, ParentID: -1}}
}

// ResolveTree resolves a tree
func ResolveTree(t *Tree, s interface{}) (interface{}, error) {
	return t.Root.NextNode(s)
}

// GetFunctionsNames get functions names
func (t *Tree) GetFunctionsNames() []string {
	fnNames := make([]string, len(t.Functions))
	i := 0
	for k := range t.Functions {
		fnNames[i] = k
		i++
	}
	return fnNames
}

// MarshalJSON ...
func (t *Tree) MarshalJSON() ([]byte, error) {
	type TreeAlias Tree
	allNodes := getAllNodes(t.Root)
	return json.Marshal(&struct {
		Nodes []*Node `json:"nodes"`
		*TreeAlias
	}{
		Nodes:     allNodes,
		TreeAlias: (*TreeAlias)(t),
	})
}

// UnmarshalJSON ...
func (t *Tree) UnmarshalJSON(data []byte) error {
	type TreeAlias Tree
	auxTree := &struct {
		Nodes []*Node `json:"nodes"`
		*TreeAlias
	}{
		Nodes:     []*Node{},
		TreeAlias: (*TreeAlias)(t),
	}
	err := json.Unmarshal(data, auxTree)
	if err != nil {
		return err
	}
	keyParentOf := map[int][]*Node{}
	for _, n := range auxTree.Nodes {
		if err := addPreprocessFn(t, n); err != nil {
			return err
		}
		if isRoot(n) {
			t.Root = n
			continue
		}
		val, found := keyParentOf[n.ParentID]
		if !found {
			keyParentOf[n.ParentID] = []*Node{n}
		} else {
			keyParentOf[n.ParentID] = append(val, n)
		}
	}
	setChildesToParentNodes(t.Root, keyParentOf)
	return nil
}

func isRoot(n *Node) bool {
	return n.ID == 0
}

func addNewPreProcessFn(newPreProcessFn []functions.PreProcessFn) map[string]functions.PreProcessFn {
	newPreProcessFn = append(newPreProcessFn, DefaultFns...)
	res := map[string]functions.PreProcessFn{}
	for _, p := range newPreProcessFn {
		res[p.Name] = p
	}
	return res
}

func addPreprocessFn(t *Tree, n *Node) error {
	if !n.PreProcessFn.Empty() {
		preProcessFn, ok := t.Functions[n.PreProcessFn.Name]
		if !ok {
			return errors.New("unmarshalling tree failed: function name not found")
		}
		n.PreProcessFn.Function = preProcessFn.Function
	}
	return nil
}

func setChildesToParentNodes(root *Node, keyParentOf map[int][]*Node) {
	queue := []*Node{root}
	for len(queue) != 0 {
		top := queue[0]
		queue = queue[1:]
		for _, n := range keyParentOf[top.ID] {
			top.Childes = append(top.Childes, n)
			queue = append(queue, n)
		}
	}
}

func getAllNodes(root *Node) []*Node {
	var res []*Node
	queue := []*Node{root}
	for len(queue) != 0 {
		top := queue[0]
		queue = queue[1:]
		res = append(res, top)
		for _, c := range top.Childes {
			queue = append(queue, c)
		}
	}
	return res
}
