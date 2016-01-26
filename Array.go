//  ---------------------------------------------------------------------------
//
//  Array.go
//
//  Copyright (c) 2016, Jared Chavez. 
//  All rights reserved.
//
//  Use of this source code is governed by a BSD-style
//  license that can be found in the LICENSE file.
//
//  -----------

package json

import (
    "encoding/json"
    "fmt"
)

// Array implements the Element interface for JSON Array types.
type Array struct {
    children []Element
    dType    int
    name     string
    parent   Element
    value    []interface{}
}

// NewArray returns a new Array object as an Element interface type.
func NewArray(name string) Element {
    return &Array {
        children : make([]Element, 0),
        dType    : TypeArray,
        name     : name,
        value    : make([]interface{}, 0),
    }
}

// AppendChild appends a new child Element to this Array object.
func (this *Array) AppendChild(child Element) error {
    if child == nil {
        return fmt.Errorf("Can't append nil child")
    }

    child.SetParent(this)
    this.children = append(this.children, child)

    this.value = append(this.value, child.Value())
    
    return nil
}

// Children returns a slice of the child Elements in this Array.
func (this *Array) Children() []Element {
    return this.children
}

// ChildrenLen returns the number of children present in this
// Array.
func (this *Array) ChildrenLen() int {
    return len(this.children)
}

// Delete attempts to remove this Element from its parent Object
// or Array.
func (this *Array) Delete() {
    deleteElement(this.parent, this)
}

// MarshalJSON implements the standard golang json marshaller 
// interface for Array.
func (this *Array) MarshalJSON() ([]byte, error) {
    return json.Marshal(this.value)
}

// Name returns the parent object's key name for this Array.
func (this *Array) Name() string {
    return this.name
}

// Parent returns a reference to this Array's parent element,
// if one exists.
func (this *Array) Parent() Element {
    return this.parent
}

// Set sets the internal value of this Array element to the
// provided value. val should be of the type []interface{}
// or an error will occur. Set will also cascade the creation
// of all child elements present within val.
func (this *Array) Set(val interface{}) error {
    valArr, ok := val.([]interface{})
    if !ok {
        return fmt.Errorf("Cannot set value. Invalid type")
    }

    this.children = make([]Element, 0)
    this.value    = make([]interface{}, 0)
    
    for i := range valArr {
        item, err := newItem(this, "", valArr[i])
        if  err != nil {
            return err
        }

        err = this.AppendChild(item)
        if err != nil {
            return err
        }
    }

    return nil
}

// SetName sets the parent object's key value for this Array.
func (this *Array) SetName(name string) {
    this.name = name
}

// SetParent sets the parent Element for this Array.
func (this *Array) SetParent(parent Element) {
    this.parent = parent
}

// String pretty-prints this Array object.
func (this *Array) String() string {
    jsonBytes, _ := json.MarshalIndent(&this.value, "", "    ")
    return string(jsonBytes)
}

// Type returns the type of JSON Element (TypeArray)
func (this *Array) Type() int {
    return this.dType
}

// Value returns the underlying value object (of type []interface{})
// for this Array.
func (this *Array) Value() interface{} {
    return this.value
}
