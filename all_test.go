//  ---------------------------------------------------------------------------
//
//  all_test.go
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
    "log"
    "testing"
)

func TestParseJSON(t *testing.T) {
    // test parse
    e := Parse([]byte(jsonTxt))
    if e == EmptyElement {
        t.Fatal()
    }

    // test object key-val
    if Search(e, "object.key").Value().(string) != "val" {
        t.Fatal()
    }

    // test array indexing with sub-object
    ex1 := map[string]float64 { 
        "arr1o2v1" : 4, 
        "arr1o2v2" : 5, 
        "arr1o2v3" : 6,
    }
    val := Search(e, "object.arrayObj[1]")
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
    }

    // test array ordering
    ex2 := []string { "sub1", "sub2", "sub3" } 
    val = Search(e, "object.array2[3]")
    if val.ChildrenLen() != 3{
        t.Fatal()
    }
    for i, v := range val.Children() {
        if ex2[i] != v.Value().(string) {
            t.Fatalf("Expected %s, got %s\n", ex2[i], v.Value())
        }
    }

    // test array indexing traversal
    if Search(e, "object.arrayObj[0].arr1o1v2").Value().(float64) != 2 {
        t.Fatal()
    }

    // test multiple array traversal
    if Search(e, "object.array2[3][2]").Value().(string) != "sub3" {
        t.Fatal()
    }

    // test a deeply embedded value
    if Search(e, "object.obj2.obj3.obj3Name").Value().(string) != "indepth" {
        t.Fatal()
    }

    // test all
    if Search(e, "object.arrayObj[0].arr1o1v3.arrb[0]").Value().(float64) != 1773 {
        t.Fatal()
    }

    // test delete
    Search(e, "object.arrayObj[0].arr1o1v1").Delete()
    if Search(e, "object.arrayObj[0].arr1o1v1") != EmptyElement {
        t.Fatal()
    }

    log.Println(e)
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

    if Search(body, "myobject.myarray[2]").Value().(int) != 5 {
        t.Fatal()
    }

    if Search(body, "myobject.mykey").Value().(string) != "myval" {
        t.Fatal()
    }

    val := Search(body, "myOtherObject")
    if val == EmptyElement {
        t.Fatal()
    }

    if Search(val, "myKey2").Value().(string) != "myVal2" {
        t.Fatal()
    }

    log.Println(body)
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
