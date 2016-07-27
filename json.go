//  ---------------------------------------------------------------------------
//
//  json.go
//
//  Copyright (c) 2016, Jared Chavez.
//  All rights reserved.
//
//  Use of this source code is governed by a BSD-style
//  license that can be found in the LICENSE file.
//
//  -----------

// Package json provides types and functions for parsing, building, and
// working with JSON documents.
package json

import (
    "encoding/json"
    "fmt"
    "strconv"
    "strings"
)

// Element type enumeration.
const (
    TypeObject = iota
    TypeArray
    TypeValue
)

var typeLookup = map[int]string{
    0: "Object",
    1: "Array",
    2: "Value",
}

// EmptyElement represents an empty JSON Element of any type
// (though it is presented as an Object). EmptyElement is returned
// during Search functions when no element is found.
var EmptyElement = NewObject("empty")

// Element presents the interface that the different types of JSON
// elements must implement. The existing types of elements are Array,
// Object, and Value.
type Element interface {
    AppendChild(Element) error
    Children() []Element
    ChildrenLen() int
    Delete()
    Name() string
    Parent() Element
    Set(interface{}) error
    SetParent(Element)
    Value() interface{}
}

// Delete attempts to delete the element at the given path from the
// supplied Element object.
func Delete(e Element, path string) error {
    el := Search(e, path)
    if el == EmptyElement {
        return fmt.Errorf("Path not found: %s", path)
    }

    el.Delete()

    return nil
}

// Parse reads a block of JSON data and attempts to parse it
// into a JSON Element object and sub-objects.
func Parse(data []byte) (Element, error) {
    var jsonData interface{}

    err := json.Unmarshal(data, &jsonData)
    if err != nil {
        return EmptyElement, err
    }

    obj, err := newItem(nil, "", jsonData)
    if err != nil {
        return EmptyElement, err
    }

    return obj, nil
}

// Search allows you to search through the given JSON element
// using dotted notation. Square brackets [] are used to denote
// array element indexes where needed.
// Ex: myobject.myarray[1].myvalue
func Search(e Element, path string) Element {
    curVal := e

    pathParts := strings.Split(path, ".")
    for i := range pathParts {
        key, idxList := parsePathPart(pathParts[i])

        newVal, err := get(curVal, key)
        if err != nil {
            return EmptyElement
        }

        if len(idxList) > 0 {
            for x := range idxList {
                // array
                val, err := getIdx(newVal, idxList[x])
                if err != nil {
                    return EmptyElement
                }

                newVal = val
            }
        }

        curVal = newVal
    }

    return curVal
}

func deleteElement(parent Element, child Element) {
    if parent == nil {
        return
    }

    if child == nil {
        return
    }

    switch p := parent.(type) {
    case *Object:
        delete(p.children, child.Name())
        delete(p.value, child.Name())
    case *Array:
        i := 0
        for ; i < len(p.children); i++ {
            if p.children[i] == child {
                break
            }
        }

        s := p.children
        s = append(s[:i], s[i+1:]...)
        p.children = s

        ns := p.value
        ns = append(ns[:i], ns[i+1:]...)
        p.value = ns
    }
}

// get attempts to retrieve a child of the given element by key name.
func get(e Element, key string) (Element, error) {
    switch obj := e.(type) {
    case *Object:
        val, ok := obj.children[key]
        if !ok {
            return nil, fmt.Errorf("Key missing (%s)", key)
        }

        return val, nil
    case *Array:
        return nil, fmt.Errorf("Array type does not support get by key")
    case *Value:
        return nil, fmt.Errorf("Value does not support get by key")
    default:
        err := fmt.Errorf("get : Invalid type (%v)", obj)
        return nil, err
    }

    err := fmt.Errorf("get unknown error")
    return nil, err
}

// getIdx attempts to retrieve a child of the given elemnent by array index.
func getIdx(e Element, idx int) (Element, error) {
    switch obj := e.(type) {
    case *Object:
        return nil, fmt.Errorf("GetIdx not implemented on object type")
    case *Array:
        if idx < 0 {
            return nil, fmt.Errorf("Invalid index (%d)", idx)
        }

        if idx >= len(obj.children) {
            err := fmt.Errorf(
                "Idx outside array bounds (%d >= %d)",
                idx,
                len(obj.children),
            )

            return nil, err
        }
        return obj.children[idx], nil
    case *Value:
        return nil, fmt.Errorf("Value does not support get by idx")
    default:
        err := fmt.Errorf("getIdx : Invalid type (%v)", obj)
        return nil, err
    }

    err := fmt.Errorf("getIdx unknown error")
    return nil, err
}

// newItem recursively builds up a new JSON Element with the given name and parent
// from the provided data. Simple values are converted into Value elements.
// map[string]interface{} data value are converted to JSON Objects. interface{}
// arrays are converted into Array elements.
func newItem(parent Element, name string, data interface{}) (Element, error) {
    var result Element

    switch val := data.(type) {
    case map[string]interface{}:
        result = NewObject(name)
        err := result.Set(val)
        if err != nil {
            return nil, err
        }
    case []interface{}:
        result = NewArray(name)
        err := result.Set(val)
        if err != nil {
            return nil, err
        }
    case interface{}:
        result = NewValue(name, val)
    default:
        if name != "" {
            result = NewValue(name, nil)
        } else {
            result = EmptyElement
        }
    }

    return result, nil
}

// parsePathPart takes a given path element and attempts to parse it
// for key names and array indexes.
func parsePathPart(item string) (string, []int) {
    key := ""
    startIdx := 0
    idxList := make([]int, 0)

    startIdx = strings.Index(item, "[")

    if startIdx == -1 {
        return item, idxList
    }

    key = item[:startIdx]

    for {
        endIdx := startIdx + strings.Index(item[startIdx:], "]")

        if endIdx == -1 {
            return key, idxList
        }

        idx, err := strconv.ParseInt(item[startIdx+1:endIdx], 10, 32)
        if err != nil {
            return key, idxList
        }

        idxList = append(idxList, int(idx))
        tmp := strings.Index(item[endIdx:], "[")

        if tmp == -1 {
            return key, idxList
        }

        startIdx = endIdx + tmp
    }

    return key, idxList
}
