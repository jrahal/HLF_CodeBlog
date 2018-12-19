/*
Copyright (c) 2016 IBM Corporation and other Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and limitations under the License.

Contributors:
Kim Letkeman - Initial Contribution
*/

// ************************************
// KL 27 Mar 2016 add testing for mapUtils as strict RPC is coming in
// ************************************

package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

var samplesStartLine = 36

var testsamples = `
{
    "event1": {
        "assetID": "ASSET001",
        "carrier": "UPS",
        "extension": {
            "arr": ["s1", "s2", "s3"]
        },
        "location": {
        },
        "temperature": 123.456,
        "timestamp": "2016-03-17T01:51:23.51620144Z"
    },
    "event2": {
        "AssetID": "ASSET001",
        "carrier": "UPS",
        "extension": {
            "arrint": [1, 2]
        },
        "location": {
            "latitude": 123.456
        },
        "Temperature": 123.456,
        "timestamp": "2016-03-17T01:51:23.51620144Z"
    },
    "event3": {
        "assetid": "ASSET001",
        "carrier": "UPS",
        "extension": {
            "arr": []
        },
        "location": {
            "longitude": 123.456
        },
        "tEmperature": 123.456,
        "timestamp": "2016-03-17T01:51:23.51620144Z"
    }
}`

var testparm1 = `
{
    "assetID": "ASSET001",
    "carrier": "UPS",
    "temperature": 2.2,
    "integer": 2,
    "bool": true,
	"sarr": ["a","b"],
	"aa" : {
		"bb" : {
			"cc" : "d"
		}
	}
}`

func printUnmarshalError(js string, err interface{}) {
	syntax, ok := err.(*json.SyntaxError)
	if !ok {
		fmt.Println("*********** ERR trying to get syntax error location **************\n", err)
		return
	}

	start, end := strings.LastIndex(js[:syntax.Offset], "\n")+1, len(js)
	if idx := strings.Index(js[start:], "\n"); idx >= 0 {
		end = start + idx
	}

	line, pos := strings.Count(js[:start], "\n"), int(syntax.Offset)-start-1
	// note, the offset here is the line number in this file
	// of the test samples json string definition (it happens to work out)
	fmt.Printf("Error in line %d: %s \n", line+samplesStartLine, err)
	fmt.Printf("%s\n%s^\n\n", js[start:end], strings.Repeat(" ", pos))
}

func getTestObjects(t *testing.T) map[string]interface{} {
	var o interface{}
	err := json.Unmarshal([]byte(testsamples), &o)
	if err != nil {
		printUnmarshalError(testsamples, err)
		t.Fatalf("unmarshal test samples failed: %s", err)
	} else {
		omap, found := o.(map[string]interface{})
		if found {
			return omap
		}
		t.Fatalf("test samples not map shape, is: %s", reflect.TypeOf(o))
	}
	return make(map[string]interface{})
}

func getTestParms(t *testing.T) interface{} {
	var o interface{}
	err := json.Unmarshal([]byte(testparm1), &o)
	if err != nil {
		printUnmarshalError(testsamples, err)
		t.Fatalf("unmarshal test samples failed: %s", err)
	}
	return o
}

func TestContains(t *testing.T) {
	t.Log("Enter TestContains")
	o := getTestObjects(t)
	ev1, found := getObject(o, "event1.extension.arr")
	if !found {
		t.Fatal("event1.extension.arr not found")
	}
	if !contains(ev1, "s2") {
		t.Fatal("event1.extension.arr should contain s2")
	}
	if contains(ev1, "s6") {
		t.Fatal("event1.extension.arr should not contain s6")
	}
	ev2, found := getObject(o, "event2.extension.arrint")
	if !found {
		t.Fatal("event2.extension.arrint not found")
	}
	// for the next 2, remember that JSON unmarshals numbers as float64
	if !contains(ev2, float64(2)) {
		t.Fatal("event2.extension.arr should contain 2")
	}
	if contains(ev2, float64(3)) {
		t.Fatal("event2.extension.arr should not contain 3")
	}
	ev3, found := getObject(o, "event3.extension.arr")
	if !found {
		t.Fatal("event3.extension.arr not found")
	}
	if contains(ev3, "s2") {
		t.Fatal("event2.extension.arr should not contain s2")
	}
}

func TestDeepMerge(t *testing.T) {
	t.Log("Enter TestDeepMerge")
	o := getTestObjects(t)
	ev1, found := getObject(o, "event1")
	if !found {
		t.Fatal("event1 not found")
	}
	ev2, found := getObject(o, "event2")
	if !found {
		t.Fatal("event2 not found")
	}
	ev3, found := getObject(o, "event3")
	if !found {
		t.Fatal("event3 not found")
	}
	state1 := ev1
	//fmt.Printf("*** State1: %s\n", prettyPrint(state1))
	_, found = getObject(state1, "location.latitude")
	if found {
		t.Fatal("state1.location should not contain latitude")
	}
	_, found = getObject(state1, "location.longitude")
	if found {
		t.Fatal("state1.location should not contain longitude")
	}
	state2 := deepMerge(ev2, state1)
	//fmt.Printf("*** State2: %s\n", prettyPrint(state2))
	_, found = getObject(state1, "location.latitude")
	if !found {
		t.Fatal("state2.location should contain latitude")
	}
	_, found = getObject(state1, "location.longitude")
	if found {
		t.Fatal("state2.location should not contain longitude")
	}
	state3 := deepMerge(ev3, state2)
	//fmt.Printf("*** State3: %s\n", prettyPrint(state3))
	_, found = getObject(state3, "location.latitude")
	if !found {
		t.Fatal("state2.location should contain latitude")
	}
	_, found = getObject(state3, "location.longitude")
	if !found {
		t.Fatal("state3.location should contain longitude")
	}
}

func TestParms(t *testing.T) {
	fmt.Println("Enter TestContains")
	o := getTestParms(t)
	_, found := getObject(o, "assetID")
	if !found {
		t.Fatal("assetID not found")
	}
}

func TestArgsMap(t *testing.T) {
	fmt.Println("Enter TestArgsMap")
	o := getTestParms(t)
	var a ArgsMap = o.(map[string]interface{})
	_, found := getObject(a, "assetID")
	if !found {
		t.Fatal("assetID not found")
	}
}

func TestGetByType(t *testing.T) {
	fmt.Println("Enter TestByType")
	o := getTestParms(t)
	_, found := getObjectAsString(o, "assetID")
	if !found {
		t.Fatal("typeof assetID should be string")
	}

	_, found = getObjectAsStringArray(o, "sarr")
	if !found {
		t.Fatalf("typeof sarr should be []string")
	}

	_, found = getObjectAsNumber(o, "temperature")
	if !found {
		t.Fatal("typeof temperature should be number")
	}

	_, found = getObjectAsInteger(o, "temperature")
	if !found {
		t.Fatal("type of temperature should be integer")
	}

	_, found = getObjectAsInteger(o, "integer")
	if !found {
		t.Fatal("typeof integer should be integer")
	}

	_, found = getObjectAsMap(o, "aa")
	if !found {
		t.Fatal("typeof aa should be map")
	}
}

func TestPutObject(t *testing.T) {
	fmt.Println("Enter TestPutObject")
	o := getTestParms(t)

	fmt.Printf("Object before: %+v\n\n", o)

	o, ok := putObject(o, "time", time.Now())
	if !ok {
		t.Fatal("could not put time")
	}

	o, ok = putObject(o, "anInt", 1)
	if !ok {
		t.Fatal("could not put anInt")
	}

	o, ok = putObject(o, "aFloat", 1.567)
	if !ok {
		t.Fatal("could not put aFloat")
	}

	i, found := getObjectAsInteger(o, "anInt")
	if !found {
		t.Fatal("anInt not an integer")
	}
	fmt.Println("anInt: ", i, " TypeOF i: ", reflect.TypeOf(i))

	n, found := getObjectAsNumber(o, "aFloat")
	if !found {
		t.Fatal("aFloat not a float")
	}
	fmt.Println("aFloat: ", n, " TypeOF n: ", reflect.TypeOf(n))

	o, ok = putObject(o, "maintenance.status", "inventory")
	if !ok {
		t.Fatal("could not put maintenance.status")
	}

	o, ok = putObject(o, "a.b.c.d.lastmaplevel.status", "installed")
	if !ok {
		t.Fatal("could not put a.b.c.d.lastmaplevel.status")
	}

	fmt.Printf("Object after: %+v\n\n", o)

}

func TestRemoveObject(t *testing.T) {
	fmt.Println("Enter TestRemoveObject")
	o := getTestParms(t)

	fmt.Printf("Object before: %+v\n\n", o)

	o, ok := removeObject(o, "assetID")
	if !ok {
		t.Fatal("could not remove assetID")
	}

	o, ok = removeObject(o, "carrier")
	if !ok {
		t.Fatal("could not remove carrier")
	}

	o, ok = removeObject(o, "aa.bb.cc")
	if !ok {
		t.Fatal("could not remove aa.bb.cc")
	}

	fmt.Printf("Object after removal of aa.bb.cc: %+v\n\n", o)

	o, ok = removeObject(o, "aa")
	if !ok {
		t.Fatal("could not remove aa")
	}

	fmt.Printf("Object after: %+v\n\n", o)
}

func TestAsStringArray(t *testing.T) {
	fmt.Println("Enter TestAsStringArray")

	s, ok := asStringArray([]string{"a"})
	if !ok {
		t.Fatal("could convert []string{'a'} to string array")
	}
	fmt.Printf("TestAsStringArray: conversion of []string{'a'} created %#v\n", s)

	_, ok = asStringArray([]int{2, 3, 4})
	if ok {
		t.Fatal("converted []int{2, 3, 4} to string array, how?")
	}

	s, ok = asStringArray("astring")
	if !ok {
		t.Fatal("failed to convert 'astring' to string array")
	}
	fmt.Printf("TestAsStringArray: conversion of 'astring' created %#v\n", s)

	s, ok = asStringArray(`["a", "b", "c"]`)
	if !ok {
		t.Fatal("failed to convert JSON a,b,c to string array")
	}
	fmt.Printf("TestAsStringArray: conversion of JSON array ['a', 'b', 'c'] created %#v\n", s)
}

func TestAddToStringArray(t *testing.T) {
	fmt.Println("Enter TestAddToStringArray")
	o := getTestParms(t)

	fmt.Printf("Object before: %+v\n\n", o)

	o, ok := addToStringArray(o, "sarr", []string{"d", "b", "c"})
	if !ok {
		t.Fatal("could not merge [d,b,c] into sarr")
	}
	fmt.Printf("Object d,b,c added: %+v\n\n", o)

	o, ok = addToStringArray(o, "unknown", []string{"unk"})
	if !ok {
		t.Fatal("could not add new array unknown : [unk]")
	}
	r, ok := getObjectAsStringArray(o, "unknown")
	if ok {
		if r[0] != "unk" {
			t.Fatal("o[unknown] != unk")
		}
	} else {
		t.Fatal("o[unknown] is missing")
	}
	fmt.Printf("Object unkown added: %+v\n\n", o)

	o, ok = addToStringArray(o, "sarr", "astring")
	if !ok {
		t.Fatal("could not merge 'astring' into sarr")
	}
	fmt.Printf("Object astring added: %+v\n\n", o)

	fmt.Printf("Object after: %+v\n\n", o)

	// next one destroys o
	o, ok = addToStringArray("", "sarr", "astring")
	if ok {
		t.Fatal("passed in string instead of map, but it did not fail")
	}

	fmt.Printf("Object after destruction by not checking for nil: %+v\n\n", o)

}

func TestRemoveFromStringArray(t *testing.T) {
	fmt.Println("Enter TestRemoveFromStringArray")
	o := getTestParms(t)

	fmt.Printf("Object before: %+v\n\n", o)

	o, ok := removeFromStringArray(o, "sarr", []string{"d", "b", "c"})
	if !ok {
		t.Fatal("could not remove [d,b,c] from sarr")
	}
	if test, ok := getObjectAsStringArray(o, "sarr"); !ok {
		t.Fatal("sarr is missing from test data")
	} else {
		if len(test) != 1 {
			t.Fatal("sarr should just contain one entry")
		}
	}

	o, ok = removeFromStringArray(o, "unknown", []string{"unk"})
	if ok {
		t.Fatal("successfully removed non-existent entry from unknown : [unk]")
	}

	o, ok = removeFromStringArray(o, "sarr", "a")
	if !ok {
		t.Fatal("could not remove 'a' from sarr")
	}
	if test, ok := getObjectAsStringArray(o, "sarr"); !ok {
		t.Fatal("sarr is missing from test data")
	} else {
		if len(test) > 0 {
			t.Fatal("sarr should be empty")
		}
	}

	fmt.Printf("Object after: %+v\n\n", o)

}
