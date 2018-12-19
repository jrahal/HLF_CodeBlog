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
// KL 14 Mar 2016 backport from 4.0
// KL 27 Mar 2016 support for loose JSON-RPC naming for 3.0.5
// KL 10 Aug 2016 added more flexibility in "getObjectAsInteger"
// ************************************

package main // sitting beside the main file for now

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// CASESENSITIVEMODE defines whether property names in the EVENT have to strictly
// follow JSON RPC conventions of case matching. Default is loose matching, but
// secure environments should turn this to true.
var CASESENSITIVEMODE = false

func asMap(obj interface{}) (toMap map[string]interface{}, ok bool) {
	var err error
	am, found := obj.(ArgsMap)
	if found {
		toMap = map[string]interface{}(am)
		return toMap, true
	}
	toMap, found = obj.(map[string]interface{})
	if found {
		return toMap, true
	}
	as, found := obj.(string)
	if found {
		var data interface{}
		err := json.Unmarshal([]byte(as), &data)
		if err == nil {
			return asMap(interface{}(data))
		}
	}
	err = fmt.Errorf("asMap: incoming type is %T and is not understood", obj)
	log.Error(err)
	return nil, false
}

func asStringArray(obj interface{}) (toSarr []string, ok bool) {
	var err error
	// 1. array of interface{}, which should of course contain strings
	sa, ok := obj.([]interface{})
	if ok {
		for i, el := range sa {
			sel, ok := el.(string)
			if !ok {
				err = fmt.Errorf("asStringArray: incoming element %d type is %T from array %#v and is not understood", i, el, obj)
				log.Error(err)
				return nil, false
			}
			toSarr = append(toSarr, sel)
		}
		return toSarr, true
	}
	// 2. array of strings, nothing to do
	toSarr, ok = obj.([]string)
	if ok {
		return toSarr, true
	}
	// what about a string argument?
	as, ok := obj.(string)
	if ok {
		if len(as) > 0 && as[0] == '[' {
			// 3. encoded JSON array of strings, unmarshall and call recursively if successful
			var data interface{}
			err := json.Unmarshal([]byte(as), &data)
			if err == nil {
				return asStringArray(interface{}(data))
			}
			log.Error(err)
			return nil, false
		}
		// 4. a non-JSON string, just return that as an array
		return []string{as}, true
	}
	err = fmt.Errorf("asStringArray: incoming type is %T and is not understood", obj)
	log.Error(err)
	return nil, false
}

// finds an object by its qualified name, which looks like "location.latitude"
// as one example. Returns as map[string]interface{}
func getObject(objIn interface{}, qname string) (interface{}, bool) {
	// return a copy of the selected object
	// handles full qualified name, starting at object's root
	obj, found := objIn.(map[string]interface{})
	if !found {
		objam, found := objIn.(ArgsMap)
		if !found {
			log.Errorf("getObject for %s passed a non-map / non-ArgsMap: type %T\n%#v", qname, objIn, objIn)
			return nil, false
		}
		obj = map[string]interface{}(objam)
	}
	searchObj := map[string]interface{}(obj)
	s := strings.Split(qname, ".")
	// crawl the levels
	for i, v := range s {
		//fmt.Printf("**** FIND level [%d] %s\n", i, v)
		//fmt.Printf("**** FIND level [%d] %s in %+v\n", i, v, searchObj)
		if i+1 < len(s) {
			tmp, found := searchObj[v]
			//fmt.Printf("** tmp is %+v\n", tmp)
			if found {
				searchObj, found = tmp.(map[string]interface{})
				//fmt.Printf("** tmp->searchObj AS MAP is %+v\n", searchObj)
				if !found {
					searchObj, found = tmp.(ArgsMap)
					if !found {
						log.Debugf("getObject not map or ArgsMap shape: %s", v)
						//fmt.Printf("** tmp->searchObj AS ARGSMAP is %+v\n", searchObj)
					}
				}
			} else {
				log.Debugf("getObject cannot find level: %s in %s", v, qname)
				return nil, false
			}
		} else {
			returnObj, found := searchObj[v]
			if !found {
				// this debug statement is not useful normally as we must be able to
				// handle assetID as part of iot common and as parameter on its own
				// so we get false warnings on read functions, but do enable it if
				// having problems with deep nested structures
				log.Debugf("getObject cannot find final level: %s in %s", v, qname)
				return nil, false
			}
			//fmt.Printf("**** Found level [%d] %s\n", i, v)
			return returnObj, true
		}
	}
	return nil, false
}

// inserts an object by its qualified name, which looks like "location.latitude"
// as one example. Creates missing levels.
func putObject(objIn interface{}, qname string, value interface{}) (interface{}, bool) {
	// overwrite the value of the selected object, create if necessary
	// handles full qualified name, starting at object's root
	searchObj, ok := asMap(objIn)
	if !ok {
		log.Errorf("getObject passed a non-map / non-ArgsMap: %#v", objIn)
		return objIn, false
	}
	s := strings.Split(qname, ".")
	// crawl the levels
	for i, v := range s {
		//fmt.Printf("**** FIND level [%d] %s\n", i, v)
		//fmt.Printf("**** FIND level [%d] %s in %+v\n", i, v, searchObj)
		if i+1 < len(s) {
			tmp, found := searchObj[v]
			//fmt.Printf("** tmp is %+v\n", tmp)
			if found {
				searchObj, found = tmp.(map[string]interface{})
				//fmt.Printf("** tmp->searchObj AS MAP is %+v\n", searchObj)
				if !found {
					searchObj, found = tmp.(ArgsMap)
					//fmt.Printf("** tmp->searchObj AS ARGSMAP is %+v\n", searchObj)
					if !found {
						// unknown object shape for a non-leaf level
						log.Errorf("putObject: unknown object shape for a non-leaf level: %+v", tmp)
						return objIn, false
					}
				}
			} else {
				//fmt.Printf("** putObject level not found in obj %+v, creating %s\n", searchObj, v)
				// level not found, create it and reset searchObj
				searchObj[v] = make(map[string]interface{})
				searchObj = searchObj[v].(map[string]interface{})
			}
		} else {
			//fmt.Printf("** putObject leaf node to be written into obj %+v, creating %s with value %+v\n", searchObj, v, value)
			// leaf node, assign the value and return
			searchObj[v] = value
			//fmt.Printf("**** Found level [%d] %s\n", i, v)
			return objIn, true
		}
	}
	log.Error("putObject: unknown error -- fell out of loop without returning")
	return objIn, false
}

// removes an object by its qualified name, which looks like "location.latitude"
// as one example. Removes missing levels.
func removeObject(objIn interface{}, qname string) (interface{}, bool) {
	// remove the value of the selected object, remove empty ancestors
	// handles full qualified name, starting at object's root
	searchObj, ok := asMap(objIn)
	if !ok {
		log.Errorf("getObject passed a non-map / non-ArgsMap: %#v", objIn)
		return objIn, false
	}
	s := strings.Split(qname, ".")
	// crawl the levels
	for i, v := range s {
		//fmt.Printf("**** FIND level [%d] %s\n", i, v)
		//fmt.Printf("**** FIND level [%d] %s in %+v\n", i, v, searchObj)
		if i+1 < len(s) {
			tmp, found := searchObj[v]
			//fmt.Printf("** tmp is %+v\n", tmp)
			if found {
				searchObj, found = tmp.(map[string]interface{})
				//fmt.Printf("** tmp->searchObj AS MAP is %+v\n", searchObj)
				if !found {
					searchObj, found = tmp.(ArgsMap)
					//fmt.Printf("** tmp->searchObj AS ARGSMAP is %+v\n", searchObj)
					if !found {
						// unknown object shape for a non-leaf level
						log.Errorf("putObject: unknown object shape for a non-leaf level: %+v", tmp)
						return objIn, false
					}
				}

			} else {
				//fmt.Printf("** putObject level not found in obj %+v, creating %s\n", searchObj, v)
				// level not found, create it and reset searchObj
				searchObj[v] = make(map[string]interface{})
				searchObj = searchObj[v].(map[string]interface{})
			}
		} else {
			//fmt.Printf("** putObject leaf node to be written into obj %+v, creating %s with value %+v\n", searchObj, v, value)
			// leaf node, assign the value and return
			delete(searchObj, v)
			//fmt.Printf("**** Found level [%d] %s\n", i, v)
			return objIn, true
		}
	}
	log.Error("putObject: unknown error -- fell out of loop without returning")
	return objIn, false
}

// Merges a specified object (usually in asset state) by qualified name with an incoming
// string or string array. Keeps only unique members (as in a set.)
func addToStringArray(objIn interface{}, qname string, value interface{}) (interface{}, bool) {
	_, ok := asMap(objIn)
	if !ok {
		log.Errorf("addToStringArray: incoming object not a map type %T[%#v]", objIn, objIn)
		return objIn, false
	}
	arr, found := getObject(objIn, qname)
	if !found {
		// array object does not exist yet, for convenience we create it
		log.Debugf("addToStringArray: redirecting to putObject of %s : %#v", qname, value)
		return putObject(objIn, qname, value)
	}
	log.Debugf("addToStringArray: adding %#v to %s : %#v\n", value, qname, arr)
	// we have an array to modify
	sarr, ok := asStringArray(arr)
	if !ok {
		log.Errorf("addToStringArray: incoming object is not a string array %s : %T[%#v]", qname, arr, arr)
		return objIn, false
	}
	sval, ok := asStringArray(value)
	if !ok {
		log.Errorf("addToStringArray: incoming value is not a string or string array %s : %T[%#v]", qname, value, value)
		return objIn, false
	}
	for _, v := range sval {
		if contains(sarr, v) {
			continue
		} else {
			sarr = append(sarr, v)
		}
	}
	sort.Strings(sarr)
	log.Debugf("addToStringArray: calling putObject with result %#v\n", sarr)
	return putObject(objIn, qname, sarr)
}

// Removes from a named object in asset state or other map, an incoming
// string or string array. Assumes unique members (as in a set.)
func removeFromStringArray(objIn interface{}, qname string, value interface{}) (interface{}, bool) {
	_, ok := asMap(objIn)
	if !ok {
		log.Errorf("removeFromStringArray: incoming object not a map type %T[%#v]", objIn, objIn)
		return objIn, false
	}
	arr, found := getObject(objIn, qname)
	if !found {
		log.Debugf("addToStringArray: array %s not found", qname)
		return objIn, false
	}
	log.Debugf("removeFromStringArray: removing %#v from %s : %#v\n", value, qname, arr)
	sarr, ok := asStringArray(arr)
	if !ok {
		log.Errorf("removeFromStringArray: incoming object is not a string array %s : %T[%#v]", qname, arr, arr)
		return objIn, false
	}
	sval, ok := asStringArray(value)
	if !ok {
		log.Errorf("removeFromStringArray: incoming value is not a string or string array %s : %T[%#v]", qname, value, value)
		return objIn, false
	}
	r := sarr[:0]
	for _, v := range sarr {
		fmt.Printf("r: %#v v: %#v sval: %#v\n", r, v, sval)
		if !contains(sval, v) {
			r = append(r, v)
		}
	}
	sort.Strings(r)
	log.Debugf("addToStringArray: calling putObject with result %#v\n", r)
	return putObject(objIn, qname, r)
}

func getObjectAsMap(objIn interface{}, qname string) (map[string]interface{}, bool) {
	amap, found := getObject(objIn, qname)
	if found {
		t, found := asMap(amap)
		if found {
			return t, true
		}
		log.Warningf("getObjectAsMap object is not a map: %s but rather %T", qname, objIn)
	}
	return nil, false
}

func getObjectAsString(objIn interface{}, qname string) (string, bool) {
	tbytes, found := getObject(objIn, qname)
	if found {
		t, found := tbytes.(string)
		if found {
			return t, true
		}
		log.Warningf("getObjectAsString object is not a string: %s", qname)
	}
	return "", false
}

func getObjectAsStringArray(objIn interface{}, qname string) ([]string, bool) {
	tbytes, found := getObject(objIn, qname)
	if found {
		return asStringArray(tbytes)
	}
	return []string{}, false
}

func getObjectAsBoolean(objIn interface{}, qname string) (bool, bool) {
	tbytes, found := getObject(objIn, qname)
	if found {
		t, found := tbytes.(bool)
		if found {
			return t, true
		}
		log.Warningf("getObjectAsBoolean object is not a boolean: %s", qname)
	}
	return false, false
}

func getObjectAsNumber(objIn interface{}, qname string) (float64, bool) {
	tbytes, found := getObject(objIn, qname)
	if found {
		t, found := tbytes.(float64)
		if found {
			return t, true
		}
		log.Warningf("getObjectAsNumber object is not a number (float64): %s", qname)
	}
	return 0, false
}

func getObjectAsInteger(objIn interface{}, qname string) (int, bool) {
	tbytes, found := getObject(objIn, qname)
	if found {
		// try as int first
		i, found := tbytes.(int)
		if found {
			return i, true
		}
		// try as JSON number and then cast
		f, found := tbytes.(float64)
		if found {
			return int(f), true
		}
		log.Warningf("getObjectAsInteger object is not an integer: %s", qname)
	}
	return 0, false
}

// this small function isolates the getting of the object in case
// sensitive or case insensitive modes because they are quite different
// we must not modify the destination key

func findObjectByKey(objIn interface{}, key string) (interface{}, bool) {
	objMap, found := objIn.(map[string]interface{})
	if found {
		dstKey, found := findMatchingKey(objMap, key)
		if found {
			objOut, found := objMap[dstKey]
			if found {
				return objOut, found
			}
		}
	}
	return nil, false
}

// finds a key that matches the incoming key, very useful to remove the
// complexity of switching case insensitivity because we always need
// the destination key to stay intact to avoid making copies of that
// substructure as we copy fields from the incoming structure
func findMatchingKey(objIn interface{}, key string) (string, bool) {
	objMap, found := objIn.(map[string]interface{})
	if !found {
		// not a map, cannot proceed
		log.Warningf("findMatchingKey objIn is not a map shape %+v", objIn)
		return "", false
	}
	if CASESENSITIVEMODE {
		// we can just use the key directly
		_, found := objMap[key]
		return key, found
	}
	// we must visit all keys and compare using tolower on each side
	for k := range objMap {
		if strings.ToLower(k) == strings.ToLower(key) {
			//log.Debugf("findMatchingKey found match! %+v %+v", k, key)
			return k, true
		}
	}
	//log.Debugf("findMatchingKey did not find key %+v", key)
	return "", false
}

// in a contract, src is usually the incoming update event,
// and dst is the existing state from the ledger

func contains(arr interface{}, val interface{}) bool {
	switch t := arr.(type) {
	case []string:
		arr2 := arr.([]string)
		for _, v := range arr2 {
			if v == val {
				return true
			}
		}
	case []int:
		arr2 := arr.([]int)
		for _, v := range arr2 {
			if v == val {
				return true
			}
		}
	case []float64:
		arr2 := arr.([]float64)
		for _, v := range arr2 {
			if v == val {
				return true
			}
		}
	case []interface{}:
		//todo: try cast instead of assertion
		//todo: use schema to determine if we even call this function or just add the value
		arr2 := arr.([]interface{})
		for _, v := range arr2 {
			switch tt := val.(type) {
			case string:
				if v.(string) == val.(string) {
					return true
				}
			case int:
				if v.(int) == val.(int) {
					return true
				}
			case float64:
				if v.(float64) == val.(float64) {
					return true
				}
			case interface{}:
				if v.(interface{}) == val.(interface{}) {
					return true
				}
			default:
				log.Errorf("contains passed array containing unknown type: %+v\n", tt)
				return false
			}
		}
	default:
		log.Errorf("contains passed array of unknown type: %+v\n", t)
		return false
	}
	return false
}

// deep merge src into dst and return dst
func deepMerge(srcIn interface{}, dstIn interface{}) interface{} {
	src, found := srcIn.(map[string]interface{})
	if !found {
		src, found = srcIn.(ArgsMap)
		if !found {
			log.Criticalf("Deep Merge passed source map of type: %s", reflect.TypeOf(srcIn))
			return nil
		}
	}
	dst, found := dstIn.(map[string]interface{})
	if !found {
		dst, found = dstIn.(ArgsMap)
		if !found {
			log.Criticalf("Deep Merge passed dest map of type: %s", reflect.TypeOf(dstIn))
			return nil
		}
	}
	for k, v := range src {
		switch v.(type) {
		case map[string]interface{}:
			// don't try hoisting dstKey calculation
			dstKey, found := findMatchingKey(dst, k)
			if found {
				dstChild, found := dst[dstKey].(map[string]interface{})
				if found {
					// recursive deepMerge into existing key
					dst[dstKey] = deepMerge(v.(map[string]interface{}), dstChild)
				}
			} else {
				// copy entire map to incoming key
				dst[k] = v
			}
		case []interface{}:
			dstKey, found := findMatchingKey(dst, k)
			if found {
				dstChild, found := dst[dstKey].([]interface{})
				if found {
					// union
					for elem := range v.([]interface{}) {
						if !contains(dstChild, elem) {
							dstChild = append(dstChild, elem)
						}
					}
				}
			} else {
				// copy
				dst[k] = v
			}
		default:
			// copy discrete types
			dstKey, found := findMatchingKey(dst, k)
			if found {
				dst[dstKey] = v
			} else {
				dst[k] = v
			}
		}
	}
	return dst
}

// returns a string that is nicely indented
// if json fails for some reason, returns the %#v representation
func prettyPrint(m interface{}) string {
	bytes, err := json.MarshalIndent(m, "", "  ")
	if err == nil {
		return string(bytes)
	}
	return fmt.Sprintf("%#v", m)
}
