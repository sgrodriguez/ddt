package ddt

import (
	"encoding/json"
	"errors"
)

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
		children := keyParentOf[n.ParentID]
		keyParentOf[n.ParentID] = append(children, n)
	}
	setChildrenToParentNodes(t.Root, keyParentOf)
	return nil
}

func getAllNodes(root *Node) []*Node {
	var res []*Node
	queue := []*Node{root}
	for len(queue) != 0 {
		top := queue[0]
		queue = queue[1:]
		res = append(res, top)
		for _, c := range top.Children {
			queue = append(queue, c)
		}
	}
	return res
}

func setChildrenToParentNodes(root *Node, keyParentOf map[int][]*Node) {
	queue := []*Node{root}
	for len(queue) != 0 {
		top := queue[0]
		queue = queue[1:]
		for _, n := range keyParentOf[top.ID] {
			top.Children = append(top.Children, n)
			queue = append(queue, n)
		}
	}
}

func addPreprocessFn(t *Tree, n *Node) error {
	if !n.PreProcessFn.Empty() {
		preProcessFn, ok := t.Functions[n.PreProcessFn.Name]
		if !ok {
			return errors.New("function name not found")
		}
		n.PreProcessFn.Function = preProcessFn.Function
	}
	return nil
}
