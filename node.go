package ddt

import (
	"encoding/json"
	"errors"
	"github.com/sgrodriguez/ddt/comparator"

	"github.com/sgrodriguez/ddt/function"
	"github.com/sgrodriguez/ddt/value"
)

// Comparer interface
type Comparer interface {
	Compare(a, b interface{}) bool
}

// Node Type
type Node struct {
	Childes []*Node `json:"-"`

	ID             int                   `json:"id"`
	ParentID       int                   `json:"parentId"`
	PreProcessFn   function.PreProcessFn `json:"-"`
	PreProcessArgs []*value.Value        `json:"preProcessFnArgs,omitempty"`
	Comparer       Comparer              `json:"comparer,omitempty"`
	ValueToCompare *value.Value          `json:"valueToCompare,omitempty"`
	Result         *value.Value          `json:"result,omitempty"`
}

// NextNode ...
func (n *Node) NextNode(input interface{}) (interface{}, error) {
	if len(n.Childes) == 0 {
		return n.Result.Value, nil
	}
	resValue, err := getValueToCompare(input, n.PreProcessFn, n.PreProcessArgs)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Childes {
		if c.Comparer.Compare(resValue, c.ValueToCompare.Value) {
			return c.NextNode(input)
		}
	}
	return nil, errors.New("value not found when comparing with all childes nodes")
}

func getValueToCompare(input interface{}, fn function.PreProcessFn, args []*value.Value) (interface{}, error) {
	if !fn.Empty() {
		resValue, err := fn.Function(input, value.GetValueInterfaces(args)...)
		if err != nil {
			return nil, err
		}
		return resValue, nil
	}
	// pre processing the input value not need it.
	return input, nil
}

// MarshalJSON ...
func (n *Node) MarshalJSON() ([]byte, error) {
	type NodeAlias Node
	return json.Marshal(&struct {
		PreProcessFn string `json:"preProcessFnName"`
		*NodeAlias
	}{
		PreProcessFn: n.PreProcessFn.Name,
		NodeAlias:    (*NodeAlias)(n),
	})
}

// UnmarshalJSON ...
func (n *Node) UnmarshalJSON(data []byte) error {
	type NodeAlias Node
	nodeAlias := &struct {
		PreProcessFn string          `json:"preProcessFnName"`
		Comparer     json.RawMessage `json:"comparer,omitempty"`
		*NodeAlias
	}{
		NodeAlias: (*NodeAlias)(n),
	}
	if err := json.Unmarshal(data, nodeAlias); err != nil {
		return err
	}
	n.PreProcessFn = function.PreProcessFn{Name: nodeAlias.PreProcessFn}
	if nodeAlias.Comparer != nil {
		comp, err := createComparatorFromJSON(nodeAlias.Comparer)
		if err != nil {
			return err
		}
		n.Comparer = comp
	}
	return nil
}

// CreateComparatorFromJSON ...
func createComparatorFromJSON(message json.RawMessage) (Comparer, error) {
	aux := &struct {
		Comp  string `json:"type"`
		Equal bool   `json:"equal"`
	}{}
	if err := json.Unmarshal(message, aux); err != nil {
		return nil, err
	}
	switch aux.Comp {
	case "eq":
		return &comparator.Equal{}, nil
	case "lt":
		return &comparator.Lesser{Equal: aux.Equal}, nil
	case "gt":
		return &comparator.Greater{Equal: aux.Equal}, nil
	}
	return nil, errors.New("invalid comparer")
}
