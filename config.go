// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

type cacheEntry struct {
	modTime  time.Time
	contents interface{}
}

type config struct {
	cache         map[string]cacheEntry
	buildFileList func(map[string]string) []string
	l             sync.Mutex
}

var c config

func init() {
	log.Println("Initializing configuration cache")
	c.cache = make(map[string]cacheEntry)
}

func lookupvar(key, path string) interface{} {
	var f interface{}
	i, err := os.Stat(path)
	_, ok := c.cache[path]
	if os.IsNotExist(err) {
		log.Println("Config does not exist", path)
		if ok {
			log.Println("Purging", path, "from cache")
			delete(c.cache, path)
		}
	}
	if err != nil {
		return nil
	}

	if c.cache[path].modTime.Before(i.ModTime()) || !ok {
		log.Println("Stale cache for", path)
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}

		err = json.Unmarshal(data, &f)
		if err != nil {
			log.Println("Cannot parse", path)
			return nil
		}

		log.Println("Updating cache for", path)
		c.cache[path] = cacheEntry{
			modTime:  i.ModTime(),
			contents: f,
		}

		return f.(map[string]interface{})[key]
	}

	log.Println("Found cache for", path)
	return c.cache[path].contents.(map[string]interface{})[key]
}

// SetFileListBuilder registers file list builder function.
//
// Registered function takes context (key-value map[string]string) as the only
// argument and return a slice of strings - file paths.
func SetFileListBuilder(f func(map[string]string) []string) {
	c.buildFileList = f
}

// Lookup searches for requested configuration key in file list built using
// context.
func Lookup(context map[string]string, key string) interface{} {
	var value interface{}

	c.l.Lock()
	defer c.l.Unlock()

	for _, fpath := range c.buildFileList(context) {
		log.Println("Context:", context, "Looking up", key, "in", fpath)
		value = lookupvar(key, fpath)
		if value != nil {
			log.Println("Context:", context, "Found", key, value)
			return value
		}
	}

	log.Println("Context:", context, "Key", key, "not found")
	return nil
}

// LookupString is analogous to Lookup(), but does the cast to string.
func LookupString(context map[string]string, key string) string {
	var value interface{}

	c.l.Lock()
	defer c.l.Unlock()

	for _, fpath := range c.buildFileList(context) {
		log.Println("Context:", context, "Looking up", key, "in", fpath)
		value = lookupvar(key, fpath)
		if value != nil {
			log.Println("Context:", context, "Found", key, value)
			return value.(string)
		}
	}

	log.Println("Context:", context, "Key", key, "not found")
	return ""
}

// LookupStringSlice is analogous to Lookup(), but does the cast to []string
func LookupStringSlice(context map[string]string, key string) (retval []string) {
	var value interface{}

	c.l.Lock()
	defer c.l.Unlock()

	for _, fpath := range c.buildFileList(context) {
		log.Println("Context:", context, "Looking up", key, "in", fpath)
		value = lookupvar(key, fpath)
		if value != nil {
			log.Println("Context:", context, "Found", key, value)
			for _, s := range value.([]interface{}) {
				retval = append(retval, s.(string))
			}
			sort.Strings(retval)
			return retval
		}
	}

	log.Println("Context:", context, "Key", key, "not found")
	return []string{""}
}

// LookupStringMap is analogous to Lookup(), but does the cast to
// map[string]bool for optimised lookups.
func LookupStringMap(context map[string]string, key string) (retval map[string]bool) {
	var value interface{}

	c.l.Lock()
	defer c.l.Unlock()

	for _, fpath := range c.buildFileList(context) {
		log.Println("Context:", context, "Looking up", key, "in", fpath)
		value = lookupvar(key, fpath)
		if value != nil {
			log.Println("Context:", context, "Found", key, value)
			for _, s := range value.([]interface{}) {
				retval[s.(string)] = true
			}
			return retval
		}
	}

	log.Println("Context:", context, "Key", key, "not found")
	return retval
}
