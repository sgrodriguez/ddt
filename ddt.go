package ddt

import (
	"errors"

	"github.com/sgrodriguez/ddt/function"
)

// Tree Type
type Tree struct {
	Root      *Node                            `json:"-"`
	Functions map[string]function.PreProcessFn `json:"-"`
	Name      string                           `json:"name"`
}

// NewTree creates a tree
func NewTree(name string, rootNode *Node, fn ...function.PreProcessFn) (*Tree, error) {
	if !isValidRootNode(rootNode) {
		return nil, errors.New("invalid root node")
	}
	return &Tree{Name: name, Functions: addNewPreProcessFn(fn), Root: rootNode}, nil
}

// ResolveTree resolves a tree given a input
func ResolveTree(t *Tree, input interface{}) (interface{}, error) {
	return t.Root.NextNode(input)
}

// DefaultFns default function
var DefaultFns = []function.PreProcessFn{
	{Function: function.StructMethod, Name: "StructMethod"},
	{Function: function.StructAttribute, Name: "StructAttribute"},
}

func addNewPreProcessFn(newPreProcessFn []function.PreProcessFn) map[string]function.PreProcessFn {
	newPreProcessFn = append(newPreProcessFn, DefaultFns...)
	res := map[string]function.PreProcessFn{}
	for _, p := range newPreProcessFn {
		res[p.Name] = p
	}
	return res
}

func isRoot(n *Node) bool {
	return n.ID == 0
}

func isValidRootNode(n *Node) bool {
	return n.ID == 0 && n.ParentID == -1 && n.Result == nil
}
