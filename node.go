package ddt

import (
	"encoding/json"
	"errors"

	"github.com/sgrodriguez/ddt/comparators"
	"github.com/sgrodriguez/ddt/functions"
	"github.com/sgrodriguez/ddt/value"
)

// Node Type
type Node struct {
	Childes []*Node `json:"-"`

	ID             int                    `json:"id"`
	ParentID       int                    `json:"parentId"`
	PreProcessFn   functions.PreProcessFn `json:"-"`
	PreProcessArgs []*value.Value         `json:"preProcessFnArgs,omitempty"`
	Comparer       comparators.Comparer   `json:"comparer,omitempty"`
	ValueToCompare *value.Value           `json:"valueToCompare,omitempty"`
	Result         *value.Value           `json:"result,omitempty"`
}

// NextNode ...
func (n *Node) NextNode(structt interface{}) (interface{}, error) {
	if len(n.Childes) == 0 {
		return n.Result, nil
	}
	resValue, err := getValueToCompare(structt, n.PreProcessFn, n.PreProcessArgs)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Childes {
		if c.Comparer.Compare(resValue, c.ValueToCompare.Value) {
			return c.NextNode(structt)
		}
	}
	return nil, errors.New("value not found when comparing in childes nodes")
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
	n.PreProcessFn = functions.PreProcessFn{Name: nodeAlias.PreProcessFn}
	comp, err := comparators.CreateComparatorFromJSON(nodeAlias.Comparer)
	if err != nil {
		return err
	}
	if comp != nil {
		n.Comparer = comp
	}
	return nil
}

func getValueToCompare(structt interface{}, fn functions.PreProcessFn, args []*value.Value) (interface{}, error) {
	var resValue interface{}
	var err error
	if !fn.Empty() {
		resValue, err = fn.Function(structt, value.GetValueInterfaces(args)...)
		if err != nil {
			return nil, err
		}
	} else {
		if len(args) > 0 {
			resValue = args[0]
		} else {
			return nil, errors.New("argument with no pre process function not found")
		}
	}
	return resValue, nil
}
