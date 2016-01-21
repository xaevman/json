package json

import (
    gojson "encoding/json"
    "log"
    "testing"
)

func TestParseJSON(t *testing.T) {
    // test parse
    e, err := ParseJSON([]byte(jsonTxt))
    if err != nil {
        t.Fatal(err)
    }

    // test object key-val
    val, err := Search(e, "object.key")
    if err != nil {
        t.Fatal(err)
    }
    if val.Value().(string) != "val" {
        t.Fatal()
    }
    log.Printf("key : %s", val.Value())

    // test array indexing with sub-object
    ex1 := map[string]float64 { 
        "arr1o2v1" : 4, 
        "arr1o2v2" : 5, 
        "arr1o2v3" : 6,
    }
    val, err = Search(e, "object.arrayObj[1]")
    if err != nil {
        t.Fatal(err)
    }
    if val.ChildrenLen() != 3 {
        t.Fatal()
    }
    for _, v := range val.Children() {
        if ex1[v.Name()] != v.Value().(float64) {
            t.Fatalf(
                "%s expected %f, got %f\n",
                v.Name(),
                ex1[v.Name()],
                v.Value(),
            )
        }
        
        log.Printf("%s : %v", v.Name(), v.Value())
    }

    // test array ordering
    ex2 := []string { "sub1", "sub2", "sub3" } 
    val, err = Search(e, "object.array2[3]")
    if err != nil {
        t.Fatal(err)
    }
    if val.ChildrenLen() != 3{
        t.Fatal()
    }
    for i, v := range val.Children() {
        if ex2[i] != v.Value().(string) {
            t.Fatalf("Expected %s, got %s\n", ex2[i], v.Value())
        }
    }

    // test array indexing traversal
    val, err = Search(e, "object.arrayObj[0].arr1o1v2")
    if err != nil {
        t.Fatal(err)
    }
    if val.Value().(float64) != 2 {
        t.Fatal()
    }
    log.Printf("val: %.0f", val.Value())

    // test multiple array traversal
    val, err = Search(e, "object.array2[3][2]")
    if err != nil {
        t.Fatal(err)
    }
    if val.Value().(string) != "sub3" {
        t.Fatal()
    }
    log.Print(val.Value())

    // test a deeply embedded value
    val, err = Search(e, "object.obj2.obj3.obj3Name")
    if err != nil {
        t.Fatal(err)
    }
    if val.Value().(string) != "indepth" {
        t.Fatal()
    }
    log.Print(val.Value())

    // test all
    val, err = Search(e, "object.arrayObj[0].arr1o1v3.arrb[0]")
    if err != nil {
        t.Fatal(err)
    }
    if val.Value().(float64) != 1773 {
        t.Fatal()
    }
    log.Print(val.Value())
}

func TestBuildJSON(t *testing.T) {
    body := NewObject("")

    myObj := NewObject("myobject")
    myArray := NewArray("myarray")
    myArray.AppendChild(NewValue("", 1))
    myArray.AppendChild(NewValue("", 3))
    myArray.AppendChild(NewValue("", 5))
    myArray.AppendChild(NewValue("", 7))
    myObj.AppendChild(myArray)
    myObj.AppendChild(NewValue("mykey", "myval"))

    myOtherObj := NewObject("myOtherObject")
    myOtherObj.AppendChild(NewValue("myKey2", "myVal2"))

    body.AppendChild(myObj)
    body.AppendChild(myOtherObj)

    val, err := Search(body, "myobject.myarray[2]")
    if err != nil {
        t.Fatal(err)
    }
    if val.Value().(int) != 5 {
        t.Fatal()
    }

    val, err = Search(body, "myobject.mykey")
    if err != nil {
        t.Fatal(err)
    }
    if val.Value().(string) != "myval" {
        t.Fatal()
    }

    val, err = Search(body, "myOtherObject")
    if err != nil {
        t.Fatal(err)
    }

    val2, err := Search(val, "myKey2")
    if err != nil {
        t.Fatal(err)
    }
    if val2.Value().(string) != "myVal2" {
        t.Fatal()
    }

    jsonBytes, err := gojson.MarshalIndent(&body, "", "    ")
    if err != nil {
        t.Fatal()
    }

    log.Println(string(jsonBytes))
}

const targetJson = `
{
    "myobject" : {
        "myarray" : [
            1,
            3,
            5,
            7
        ],
        "mykey" : "myval"
    },
    "myOtherObject" : {
        "myKey2" : "myVal2"
    }
}
`

const jsonTxt = `
{
    "object" : {
        "key" : "val",
        "arrayObj" : [
            {
                "arr1o1v1" : 1,
                "arr1o1v2" : 2,
                "arr1o1v3" : {
                    "arrb" : [
                        1773
                    ]
                }
            },
            {
                "arr1o2v1" : 4,
                "arr1o2v2" : 5,
                "arr1o2v3" : 6
            }
        ],
        "key2" : "val2",
        "obj2" : {
            "obj2Name" : "test123",
            "obj2Val"  : 123,
            "obj3"     : {
                "obj3Name" : "indepth"
            }
        },
        "array2" : [
            1,
            2,
            3,
            [
                "sub1",
                "sub2",
                "sub3"
            ]
        ]
    },
    "test" : "val"
}

`
