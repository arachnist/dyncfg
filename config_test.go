// +build go1.4

// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
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
		v := c.Lookup(emptyContext, e.key)
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
		v := c.LookupString(emptyContext, e.key)
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
		v := c.LookupInt(emptyContext, e.key)
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
		v := c.LookupStringSlice(emptyContext, e.key)
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
		v := c.LookupStringMap(emptyContext, e.key)
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

func complicatedLookupHelper(map[string]string) []string {
	return []string{
		".to-be-removed.json",
		".to-be-replaced-with-gibberish.json",
		".to-be-chmoded-000.json",
		".to-be-replaced.json",
	}
}

var testEvents = []struct {
	desc   string
	action func()
	ev     int
}{
	{
		desc:   "remove file",
		action: func() { _ = os.Remove(".to-be-removed.json") },
		ev:     1,
	},
	{
		desc: "replace with gibberish",
		action: func() {
			_ = ioutil.WriteFile(".to-be-replaced-with-gibberish.json", []byte("!@#%^&!%@$&*#@@%*@^#$^&*@&(!"), 0644)
		},
		ev: 2,
	},
	{
		desc: "chmod 000",
		action: func() {
			_ = os.Chmod(".to-be-chmoded-000.json", 0000)
			_ = os.Chtimes(".to-be-chmoded-000.json", time.Now().Local(), time.Now().Local())
		},
		ev: 3,
	},
	{
		desc:   "replace",
		action: func() { _ = ioutil.WriteFile(".to-be-replaced.json", []byte("{\"testkey\":99}"), 0644) },
		ev:     99,
	},
}

func TestLookupvarConditions(t *testing.T) {
	// create files and put their parsed contents in cache
	for i := len(complicatedLookupHelper(emptyContext)) - 1; i >= 0; i-- {
		path := complicatedLookupHelper(emptyContext)[i]
		err := ioutil.WriteFile(path, []byte(fmt.Sprintf("{\"testkey\":%d}", i)), 0644)
		if err != nil {
			t.Log("file write failed")
		}

		if i != c2.LookupInt(emptyContext, "testkey") {
			t.Log("LookupInt did not return correct value", i)
			t.Fail()
		}
	}

	time.Sleep(time.Second) // make sure we get to update cache

	for _, e := range testEvents {
		e.action()
		v := c2.LookupInt(emptyContext, "testkey")
		if e.ev != v {
			t.Log("Test failed:", e.desc, "expected:", e.ev, "got:", v)
			t.Fail()
		}
	}

	for _, p := range complicatedLookupHelper(emptyContext) {
		_ = os.Remove(p)
	}
}

var c *Config = New()
var c2 *Config = New()

func TestMain(m *testing.M) {
	c.SetFileListBuilder(configLookupHelper)
	c2.SetFileListBuilder(complicatedLookupHelper)
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}
