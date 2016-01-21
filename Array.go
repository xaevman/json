package json

import (
    "encoding/json"
    "fmt"
)

type Array struct {
    children []Element
    dType    int
    name     string
    parent   Element
    value    []interface{}
}

func NewArray(name string) Element {
    return &Array {
        children : make([]Element, 0),
        name     : name,
        value    : make([]interface{}, 0),
    }
}

func (this *Array) AppendChild(child Element) error {
    if child == nil {
        return fmt.Errorf("Can't append nil child")
    }

    child.SetParent(this)
    this.children = append(this.children, child)

    this.value = append(this.value, child.Value())
    
    return nil
}

func (this *Array) Children() []Element {
    return this.children
}

func (this *Array) ChildrenLen() int {
    return len(this.children)
}

func (this *Array) MarshalJSON() ([]byte, error) {
    return json.Marshal(this.value)
}

func (this *Array) Name() string {
    return this.name
}

func (this *Array) Parent() Element {
    return this.parent
}

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

func (this *Array) SetParent(parent Element) {
    this.parent = parent
}

func (this *Array) SetName(name string) {
    this.name = name
}

func (this *Array) String() string {
    return fmt.Sprintf(
        "type: %s, children: %d, value: %v",
        typeLookup[this.dType],
        len(this.children),
        this.value,
    )
}

func (this *Array) Type() int {
    return this.dType
}

func (this *Array) Value() interface{} {
    return this.value
}
