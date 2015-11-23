// +build go1.2,!go1.4

// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

var emptyContext map[string]string

// func Lookup(context map[string]string, key string) interface{}
var testsLookup = []struct {
	key           string
	expectedValue string
}{
	{
		key:           "TestLookup",
		expectedValue: "this can be anything you can imagine, but now it's just a string",
	},
	{
		key:           "NotExisting",
		expectedValue: "<nil>",
	},
}

func TestLookup(t *testing.T) {
	for i, e := range testsLookup {
		v := Lookup(emptyContext, e.key)
		if fmt.Sprintf("%+v", v) != fmt.Sprintf("%+v", e.expectedValue) {
			t.Log("test#", i, "Lookup key", e.key)
			t.Logf("expected: %+v\n", e.expectedValue)
			t.Logf("result  : %+v\n", v)
			t.Fail()
		}
	}
}

// func LookupString(context map[string]string, key string) string
var testsLookupString = []struct {
	key           string
	expectedValue string
}{
	{
		key:           "TestLookupString",
		expectedValue: "this should be a string",
	},
	{
		key:           "NotExisting",
		expectedValue: "",
	},
}

func TestLookupString(t *testing.T) {
	for i, e := range testsLookupString {
		v := LookupString(emptyContext, e.key)
		if fmt.Sprintf("%+v", v) != fmt.Sprintf("%+v", e.expectedValue) {
			t.Log("test#", i+1, "Lookup key", e.key)
			t.Logf("expected: %+v\n", e.expectedValue)
			t.Logf("result  : %+v\n", v)
			t.Fail()
		}
	}
}

// func LookupInt(context map[string]string, key string) int
var testsLookupInt = []struct {
	key           string
	expectedValue int
}{
	{
		key:           "ThisWillBeAnInt",
		expectedValue: 42,
	},
	{
		key:           "NotExisting",
		expectedValue: -1,
	},
}

func TestLookupInt(t *testing.T) {
	for i, e := range testsLookupInt {
		v := LookupInt(emptyContext, e.key)
		if fmt.Sprintf("%+v", v) != fmt.Sprintf("%+v", e.expectedValue) {
			t.Log("test#", i+1, "Lookup key", e.key)
			t.Logf("expected: %+v\n", e.expectedValue)
			t.Logf("result  : %+v\n", v)
			t.Fail()
		}
	}
}

// func LookupStringSlice(context map[string]string, key string) []string
var testsLookupStringSlice = []struct {
	key           string
	expectedValue []string
}{
	{
		// TestLookupStringSlice return slice is ordered
		key:           "TestLookupStringSlice",
		expectedValue: []string{"a", "am", "be", "going", "i", "slice", "to"},
	},
	{
		key:           "NotExisting",
		expectedValue: []string{""},
	},
}

func TestLookupStringSlice(t *testing.T) {
	for i, e := range testsLookupStringSlice {
		v := LookupStringSlice(emptyContext, e.key)
		if fmt.Sprintf("%+v", v) != fmt.Sprintf("%+v", e.expectedValue) {
			t.Log("test#", i+1, "Lookup key", e.key)
			t.Logf("expected: %+v\n", e.expectedValue)
			t.Logf("result  : %+v\n", v)
			t.Fail()
		}
	}
}

// func LookupStringMap(context map[string]string, key string) map[string]bool
var testsLookupStringMap = []struct {
	key           string
	expectedValue map[string]bool
}{
	{
		key: "TestLookupStringMap",
		expectedValue: map[string]bool{
			"to":              true,
			"be":              true,
			"a":               true,
			"map[string]bool": true,
			"I":               true,
			"want":            true,
		},
	},
	{
		key:           "NotExisting",
		expectedValue: map[string]bool{},
	},
}

func TestLookupStringMap(t *testing.T) {
	for i, e := range testsLookupStringMap {
		v := LookupStringMap(emptyContext, e.key)
		if !reflect.DeepEqual(v, e.expectedValue) {
			t.Log("test#", i+1, "Lookup key", e.key)
			t.Logf("expected: %+v\n", e.expectedValue)
			t.Logf("result  : %+v\n", v)
			t.Fail()
		}
	}
}

func configLookupHelper(map[string]string) []string {
	return []string{"this/thing/does/not/exist.json", ".testconfig.json"}
}

func init() {
	SetFileListBuilder(configLookupHelper)
	log.SetOutput(ioutil.Discard)
}
