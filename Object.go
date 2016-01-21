package json

import (
    "encoding/json"
    "fmt"
)

type Object struct {
    children map[string]Element
    dType    int
    name     string
    parent   Element
    value    map[string]interface{}
}

func NewObject(name string) Element {
    return &Object {
        children : make(map[string]Element),
        name     : name,
        value    : make(map[string]interface{}),
    }
}

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

func (this *Object) Children() []Element {
    results := make([]Element, 0, len(this.children))

    for k, _ := range this.children {
        results = append(results, this.children[k])
    }

    return results
}

func (this *Object) ChildrenLen() int {
    return len(this.children)
}

func (this *Object) MarshalJSON() ([]byte, error) {
    return json.Marshal(this.value)
}

func (this *Object) Name() string {
    return this.name
}

func (this *Object) Parent() Element {
    return this.parent
}

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

func (this *Object) SetName(name string) {
    this.name = name
}

func (this *Object) SetParent(parent Element) {
    this.parent = parent
}

func (this *Object) String() string {
    return fmt.Sprintf(
        "%s, type: %s, children: %d, value: %v",
        this.name,
        typeLookup[this.dType],
        len(this.children),
        this.value,
    )
}

func (this *Object) Type() int {
    return this.dType
}

func (this *Object) Value() interface{} {
    return this.value
}

