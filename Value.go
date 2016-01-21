package json

import (
    "encoding/json"
    "fmt"
)

type Value struct {
    name   string
    dType  int
    parent Element
    value  interface{}
}

func NewValue(name string, value interface{}) Element {
    return &Value{
        name  : name,
        value : value,
    }
}

func (this *Value) AppendChild(child Element) error {
    return nil
}

func (this *Value) Children() []Element {
    return nil
}

func (this *Value) ChildrenLen() int {
    return 0
}

func (this *Value) MarshalJSON() ([]byte, error) {
    return json.Marshal(this.value)
}

func (this *Value) Name() string {
    return this.name
}

func (this *Value) Parent() Element {
    return this.parent
}

func (this *Value) Set(val interface{}) error {
    this.value = val
    return nil
}

func (this *Value) SetParent(parent Element) {
    this.parent = parent
}

func (this *Value) SetName(name string) {
    this.name = name
}

func (this *Value) String() string {
    return fmt.Sprintf(
        "name: %s, type: %s, value: %v",
        this.name,
        typeLookup[this.dType],
        this.value,
    )
}

func (this *Value) Type() int {
    return this.dType
}

func (this *Value) Value() interface{} {
    return this.value
}

