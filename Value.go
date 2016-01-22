//  ---------------------------------------------------------------------------
//
//  Value.go
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

// Value implements the Element interface for Value types.
type Value struct {
    name   string
    dType  int
    parent Element
    value  interface{}
}

// NewValue returns a new Value object as an Element interface type.
func NewValue(name string, value interface{}) Element {
    return &Value{
        name  : name,
        value : value,
    }
}

// AppendChild is not implemented for Value objects.
func (this *Value) AppendChild(child Element) error {
    return nil
}

// Children is not implemented for Value objects.
func (this *Value) Children() []Element {
    return nil
}

// ChildrenLen is not implemented for Value objects.
func (this *Value) ChildrenLen() int {
    return 0
}

// MarshalJSON implements the golang json marshaller for Value 
// types.
func (this *Value) MarshalJSON() ([]byte, error) {
    return json.Marshal(this.value)
}

// Name returns the parent object's key name for this Value.
func (this *Value) Name() string {
    return this.name
}

// Parent returns a reference to this Value's parent element,
// if one exists.
func (this *Value) Parent() Element {
    return this.parent
}

// Set sets the internal value of this Value element to the
// provided value.
func (this *Value) Set(val interface{}) error {
    this.value = val
    return nil
}

// SetName sets the parent object's key value for this Value.
func (this *Value) SetName(name string) {
    this.name = name
}

// SetParent sets the parent Element for this Object.
func (this *Value) SetParent(parent Element) {
    this.parent = parent
}

// String pretty-prints this Value.
func (this *Value) String() string {
    return fmt.Sprintf(
        "name: %s, type: %s, value: %v",
        this.name,
        typeLookup[this.dType],
        this.value,
    )
}

// Type returns the type of this JSON Element (TypeValue).
func (this *Value) Type() int {
    return this.dType
}

// Value the underlying raw value for this Value object.
func (this *Value) Value() interface{} {
    return this.value
}

