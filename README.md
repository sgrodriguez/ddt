[![Go Report Card](https://goreportcard.com/badge/github.com/sgrodriguez/ddt)](https://goreportcard.com/report/github.com/sgrodriguez/ddt)
[![codecov](https://codecov.io/gh/sgrodriguez/go-ddt/branch/master/graph/badge.svg?token=8JU0YG71WZ)](https://codecov.io/gh/sgrodriguez/go-ddt)
[![Build Status](https://travis-ci.com/sgrodriguez/ddt.svg?branch=master)](https://travis-ci.com/sgrodriguez/ddt)
# DDT
##Dynamic decision tree
DDT allows building custom decision trees based in a set of defined rules, avoiding the task of writing ifs/switches statements.

## Examples
### Create user tree
Use a struct as input of the tree and pre-process data before comparing with the next level of nodes.
In this example we use some methods and attributes of the user struct.  
![alt text](docs/user_tree.png?raw=true)
```go
package main

import (
	"fmt"
	"github.com/sgrodriguez/ddt"
	"github.com/sgrodriguez/ddt/compare"
	"github.com/sgrodriguez/ddt/function"
	"github.com/sgrodriguez/ddt/value"
)

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

func main() {
	node6 := &ddt.Node{
		ID:             6,
		ParentID:       2,
		ValueToCompare: &value.Value{Type: value.Int, Value: 30},
		Result:         &value.Value{Type: value.String, Value: "node6"},
		Comparer:       &compare.Greater{},
	}
	node5 := &ddt.Node{
		ID:             5,
		ParentID:       2,
		ValueToCompare: &value.Value{Type: value.Int, Value: 30},
		Result:         &value.Value{Type: value.String, Value: "node5"},
		Comparer:       &compare.Lesser{Equal: true},
	}
	node3 := &ddt.Node{
		ID:             3,
		ParentID:       1,
		ValueToCompare: &value.Value{Type: value.String, Value: "SANTIAGO LUCIA"},
		Result:         &value.Value{Type: value.String, Value: "node3"},
		Comparer:       &compare.Equal{},
	}
	node4 := &ddt.Node{
		ID:             4,
		ParentID:       1,
		ValueToCompare: &value.Value{Type: value.String, Value: "LUCIA SANTIAGO"},
		Result:         &value.Value{Type: value.String, Value: "node4"},
		Comparer:       &compare.Equal{},
	}
	node1 := &ddt.Node{
		ID:             1,
		ParentID:       0,
		Children:       []*ddt.Node{node3, node4},
		ValueToCompare: &value.Value{Type: value.Bool, Value: true},
		PreProcessArgs: []*value.Value{{Type: value.String, Value: "FullName"}},
		PreProcessFn:   function.PreProcessFn{Function: function.CallStructMethod, Name: "CallStructMethod"},
		Comparer:       &compare.Equal{},
	}
	node2 := &ddt.Node{
		ID:             2,
		ParentID:       0,
		Children:       []*ddt.Node{node5, node6},
		ValueToCompare: &value.Value{Type: value.Bool, Value: false},
		Comparer:       &compare.Equal{},
		PreProcessArgs: []*value.Value{{Type: value.String, Value: "Age"}},
		PreProcessFn:   function.PreProcessFn{Function: function.GetStructAttribute, Name: "GetStructAttribute"},
	}
	root := &ddt.Node{
		Children:       []*ddt.Node{node1, node2},
		PreProcessArgs: []*value.Value{{Type: value.String, Value: "UnderAge"}},
		PreProcessFn:   function.PreProcessFn{Function: function.CallStructMethod, Name: "CallStructMethod"},
		ID:             0,
		ParentID:       -1,
	}
	userTree, err := ddt.NewTree("userTree", root)
	if err != nil {
		panic(err)
	}
	result, err := ddt.ResolveTree(userTree, &user{Age: 12, FirstName: "SANTIAGO", LastName: "LUCIA"})
	if err != nil {
		panic(err)
	}
	// result node3
	fmt.Println(result.(string))
}
```
### Create the user tree from json
```go
package main


import (
	"encoding/json"
	"fmt"
	"github.com/sgrodriguez/ddt"
)

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

func main() {
	// define empty tree
	tree, err := ddt.NewTree("newTree", &ddt.Node{ID: 0, ParentID: -1})
	if err != nil {
		panic(err)
	}
	treeFromJson := []byte(`
{
   "nodes":[
      {
         "preProcessFnName":"CallStructMethod",
         "id":0,
         "parentId":-1,
         "preProcessFnArgs":[
            {
               "Value":"UnderAge",
               "Type":"string"
            }
         ]
      },
      {
         "preProcessFnName":"CallStructMethod",
         "id":1,
         "parentId":0,
         "preProcessFnArgs":[
            {
               "Value":"FullName",
               "Type":"string"
            }
         ],
         "comparer":{
            "type":"eq"
         },
         "valueToCompare":{
            "Value":true,
            "Type":"bool"
         }
      },
      {
         "preProcessFnName":"GetStructAttribute",
         "id":2,
         "parentId":0,
         "preProcessFnArgs":[
            {
               "Value":"Age",
               "Type":"string"
            }
         ],
         "comparer":{
            "type":"eq"
         },
         "valueToCompare":{
            "Value":false,
            "Type":"bool"
         }
      },
      {
         "preProcessFnName":"",
         "id":3,
         "parentId":1,
         "comparer":{
            "type":"eq"
         },
         "valueToCompare":{
            "Value":"SANTIAGO LUCIA",
            "Type":"string"
         },
         "result":{
            "Value":"node3",
            "Type":"string"
         }
      },
      {
         "preProcessFnName":"",
         "id":4,
         "parentId":1,
         "comparer":{
            "type":"eq"
         },
         "valueToCompare":{
            "Value":"LUCIA SANTIAGO",
            "Type":"string"
         },
         "result":{
            "Value":"node4",
            "Type":"string"
         }
      },
      {
         "preProcessFnName":"",
         "id":5,
         "parentId":2,
         "comparer":{
            "type":"lt",
            "equal":true
         },
         "valueToCompare":{
            "Value":30,
            "Type":"int"
         },
         "result":{
            "Value":"node5",
            "Type":"string"
         }
      },
      {
         "preProcessFnName":"",
         "id":6,
         "parentId":2,
         "comparer":{
            "type":"gt",
            "equal":false
         },
         "valueToCompare":{
            "Value":30,
            "Type":"int"
         },
         "result":{
            "Value":"node6",
            "Type":"string"
         }
      }
   ],
   "name":"userTree"
}`)
	err = json.Unmarshal(treeFromJson, tree)
	if err != nil {
		panic(err)
	}
	result, err := ddt.ResolveTree(tree, &user{Age: 12, FirstName: "SANTIAGO", LastName: "LUCIA"})
	if err != nil {
		panic(err)
	}
	// result node3
	fmt.Println(result.(string))
	treeByte, err := json.Marshal(tree)
	fmt.Println(string(treeByte))
}
```
### Create simple tree
Create a simple tree using only basic types.
![alt text](docs/simple_tree.png?raw=true)
```go
package main

import (
	"fmt"
	"github.com/sgrodriguez/ddt"
	"github.com/sgrodriguez/ddt/compare"
	"github.com/sgrodriguez/ddt/value"
)

func main() {
	leaf1 := ddt.Node{
		ID:             1,
		ParentID:       0,
		ValueToCompare: &value.Value{Value: int64(60), Type: value.Int64},
		Comparer:       &compare.Greater{},
		Result:         &value.Value{Value: "prize1", Type: value.String},
	}
	leaf11 := ddt.Node{
		ID:             3,
		ParentID:       2,
		ValueToCompare: &value.Value{Value: int64(30), Type: value.Int64},
		Comparer:       &compare.Equal{},
		Result:         &value.Value{Value: "prize2", Type: value.String},
	}
	leaf12 := ddt.Node{
		ID:             4,
		ParentID:       2,
		ValueToCompare: &value.Value{Value: int64(30), Type: value.Int64},
		Comparer:       &compare.Greater{},
		Result:         &value.Value{Value: "prize3", Type: value.String},
	}
	leaf13 := ddt.Node{
		ID:             5,
		ParentID:       2,
		ValueToCompare: &value.Value{Value: int64(30), Type: value.Int64},
		Comparer:       &compare.Lesser{},
		Result:         &value.Value{Value: "prize4", Type: value.String},
	}
	node1 := ddt.Node{
		Children:       []*ddt.Node{&leaf11, &leaf12, &leaf13},
		ID:             2,
		ParentID:       0,
		ValueToCompare: &value.Value{Value: int64(60), Type: value.Int64},
		Comparer:       &compare.Lesser{Equal: true},
	}
	root := ddt.Node{
		ID:       0,
		ParentID: -1,
		Children: []*ddt.Node{&node1, &leaf1},
	}
	simpleTree, err := ddt.NewTree("simpleTree", &root)
	if err != nil {
		panic(err)
	}
	result, err := ddt.ResolveTree(simpleTree,int64(15))
	if err != nil {
		panic(err)
	}
	fmt.Println(result.(string))
}
```
### Create simple tree from json
Create or modify the simple tree from json
```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/sgrodriguez/ddt"
)

func main() {
	// define empty tree
	tree, err := ddt.NewTree("newTree", &ddt.Node{ID: 0, ParentID: -1})
	if err != nil {
		panic(err)
	}
	treeFromJson := []byte(`
{
   "nodes":[
      {
         "preProcessFnName":"",
         "id":0,
         "parentId":-1
      },
      {
         "preProcessFnName":"",
         "id":2,
         "parentId":0,
         "comparer":{
            "type":"lt",
            "equal":true
         },
         "valueToCompare":{
            "Value":60,
            "Type":"int64"
         }
      },
      {
         "preProcessFnName":"",
         "id":1,
         "parentId":0,
         "comparer":{
            "type":"gt",
            "equal":false
         },
         "valueToCompare":{
            "Value":60,
            "Type":"int64"
         },
         "result":{
            "Value":"prize1",
            "Type":"string"
         }
      },
      {
         "preProcessFnName":"",
         "id":3,
         "parentId":2,
         "comparer":{
            "type":"eq"
         },
         "valueToCompare":{
            "Value":30,
            "Type":"int64"
         },
         "result":{
            "Value":"prize2",
            "Type":"string"
         }
      },
      {
         "preProcessFnName":"",
         "id":4,
         "parentId":2,
         "comparer":{
            "type":"gt",
            "equal":false
         },
         "valueToCompare":{
            "Value":30,
            "Type":"int64"
         },
         "result":{
            "Value":"prize3",
            "Type":"string"
         }
      },
      {
         "preProcessFnName":"",
         "id":5,
         "parentId":2,
         "comparer":{
            "type":"lt",
            "equal":false
         },
         "valueToCompare":{
            "Value":30,
            "Type":"int64"
         },
         "result":{
            "Value":"prize4",
            "Type":"string"
         }
      }
   ],
   "name":"simpleTree"
}`)
	err = json.Unmarshal(treeFromJson, tree)
	if err != nil {
		panic(err)
	}
	result, err := ddt.ResolveTree(tree, int64(15))
	if err != nil {
		panic(err)
	}
	// result prize4
	fmt.Println(result.(string))
	// change some property of the tree for example the result of prize4
	modifiedTree := []byte(`{"nodes":[{"preProcessFnName":"","id":0,"parentId":-1},{"preProcessFnName":"","id":2,"parentId":0,"comparer":{"type":"lt","equal":true},"valueToCompare":{"Value":60,"Type":"int64"}},{"preProcessFnName":"","id":1,"parentId":0,"comparer":{"type":"gt","equal":false},"valueToCompare":{"Value":60,"Type":"int64"},"result":{"Value":"prize1","Type":"string"}},{"preProcessFnName":"","id":3,"parentId":2,"comparer":{"type":"eq"},"valueToCompare":{"Value":30,"Type":"int64"},"result":{"Value":"prize2","Type":"string"}},{"preProcessFnName":"","id":4,"parentId":2,"comparer":{"type":"gt","equal":false},"valueToCompare":{"Value":30,"Type":"int64"},"result":{"Value":"prize3","Type":"string"}},{"preProcessFnName":"","id":5,"parentId":2,"comparer":{"type":"lt","equal":false},"valueToCompare":{"Value":30,"Type":"int64"},"result":{"Value": 420,"Type":"int64"}}],"name":"simpleTree"}`)
	err = json.Unmarshal(modifiedTree, tree)
	if err != nil {
		panic(err)
	}
	result, err = ddt.ResolveTree(tree, int64(15))
	if err != nil {
		panic(err)
	}
	// result 420
	fmt.Println(result.(int64))
}
```

## Overview
#### Node
* ID: id of the node, root node must have 0.
* ParentID: parent id, root node must have -1.
* Result: if the node is leaf and is the next node of the tree, this is the result.
* Comparer.
* ValueToCompare: value 
* PreProcessFn: function to pre-process the input before comparing.
* PreProcessArgs.
    
#### Value
Basic types available for comparing, result and as PreProcessArgs.
   * Int
   * Int64
   * Bool
   * String
   * Float64
#### Comparators
   * Greater (or Equal)
   * Lesser  (or Equal)
   * Equal
#### Pre-Process Functions
Functions to pre-process the input before comparing with the next level of the tree.
   * CallStructMethod 
   * GetStructAttribute 


