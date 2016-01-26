//  ---------------------------------------------------------------------------
//
//  Object.go
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

// Object implements the Element interface for JSON Object types.
type Object struct {
    children map[string]Element
    dType    int
    name     string
    parent   Element
    value    map[string]interface{}
}

// NewObject returns a new Object object as an Element interface type.
func NewObject(name string) Element {
    return &Object {
        children : make(map[string]Element),
        name     : name,
        value    : make(map[string]interface{}),
    }
}

// AppendChild appends a new child Element to this Object.
func (this *Object) AppendChild(child Element) error {
    cName := child.Name()
    if cName == "" {
        return fmt.Errorf("Key is missing or nil")
    }

    _, ok := this.children[cName]
    if ok {
        return fmt.Errorf("Key already present (%s)", cName)
    }

    this.value[cName] = child.Value()

    child.SetParent(this)
    this.children[cName] = child

    return nil
}

// Children returns a slice of the child elements within this
// JSON Object. Since the underlying storage for object children
// is a map, order cannot be guaranteed.
func (this *Object) Children() []Element {
    results := make([]Element, 0, len(this.children))

    for k, _ := range this.children {
        results = append(results, this.children[k])
    }

    return results
}

// ChildrenLen returns the number of child objects contained 
// within this Object.
func (this *Object) ChildrenLen() int {
    return len(this.children)
}

// Delete attempts to remove this object from its parent object.
func (this *Object) Delete() {
    deleteElement(this.parent, this)
}

// MarshalJSON implements the standard golang json marshaller
// for this Object.
func (this *Object) MarshalJSON() ([]byte, error) {
    return json.Marshal(this.value)
}

// Name returns the parent object's key name for this Object.
func (this *Object) Name() string {
    return this.name
}

// Parent returns a reference to this Object's parent element,
// if one exists.
func (this *Object) Parent() Element {
    return this.parent
}

// Set sets the internal value of this Object element to the
// provided value. val should be of the type map[string]interface{}
// or an error will occur. Set will also cascade the creation
// of all child elements present within val.
func (this *Object) Set(val interface{}) error {
    valMap, ok := val.(map[string]interface{})
    if !ok {
        return fmt.Errorf("Cannot set value. Invalid type")
    }

    this.value    = make(map[string]interface{})
    this.children = make(map[string]Element)

    for k, _ := range valMap {
        child, err := newItem(this, k, valMap[k])
        if err != nil {
            return err
        }

        err = this.AppendChild(child)
        if err != nil {
            return err
        }
    }

    return nil
}

// SetName sets the parent object's key value for this Object.
func (this *Object) SetName(name string) {
    this.name = name
}

// SetParent sets the parent Element for this Object.
func (this *Object) SetParent(parent Element) {
    this.parent = parent
}

// String pretty-prints this Object.
func (this *Object) String() string {
    jsonBytes, _ := json.MarshalIndent(&this.value, "", "    ")
    return string(jsonBytes)
}

// Type returns the type of this JSON Element (TypeObject).
func (this *Object) Type() int {
    return this.dType
}

// Value the underlying value object (of type map[string]interface{})
// for this Object.
func (this *Object) Value() interface{} {
    return this.value
}

