
package json

import (
    "encoding/json"
    "fmt"
    "strings"
    "strconv"
)

const (
    TypeObject = iota
    TypeArray
    TypeValue
)

var typeLookup = map[int]string {
    0 : "Object",
    1 : "Array",
    2 : "Value",
}

type Element interface {
    AppendChild(Element) error
    Children()           []Element
    ChildrenLen()        int
    Name()               string
    MarshalJSON()        ([]byte, error)
    Parent()             Element
    Set(interface{})     error
    SetName(string)
    SetParent(Element)
    String()             string
    Type()               int
    Value()              interface{}
}

func ParseJSON(data []byte) (Element, error) {
    var jsonData interface{}

    err := json.Unmarshal(data, &jsonData)
    if err != nil {
        return nil, err
    }

    obj, err := newItem(nil, "", jsonData)
    if err != nil {
        return nil, err
    }

    return obj, err
}

func Search(e Element, path string) (Element, error) {
    curVal := e
    
    pathParts := strings.Split(path, ".")
    for i := range pathParts {
        key, idxList := parsePathPart(pathParts[i])
        
        newVal, err := get(curVal, key)
        if err != nil {
            return nil, err
        }

        if len(idxList) > 0 {
            for x := range idxList {
                // array
                val, err := getIdx(newVal, idxList[x])
                if err != nil {
                    return nil, err
                }

                newVal = val
            }
        }

        curVal = newVal
    }

    return curVal, nil
}

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
            err := fmt.Errorf("Invalid type (%d)", e.Type())
            return nil, err
    }

    err := fmt.Errorf("Invalid type (%d)", e.Type())
    return nil, err
}

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
            err := fmt.Errorf("Invalid type (%d)", e.Type())
            return nil, err
    }

    err := fmt.Errorf("Invalid type (%d)", e.Type())
    return nil, err
}

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
        return nil, fmt.Errorf("Invalid type: %v", val)
    }

    return result, nil
}

func parsePathPart(item string) (string, []int) {
    key      := ""
    startIdx := 0
    idxList  := make([]int, 0)

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

        idx, err := strconv.ParseInt(item[startIdx + 1:endIdx], 10, 32)
        if err != nil {
            return key, idxList
        }

        idxList = append(idxList, int(idx))
        tmp    := strings.Index(item[endIdx:], "[")

        if tmp == -1 {
            return key, idxList
        }

        startIdx = endIdx + tmp
    }

    return key, idxList
}
