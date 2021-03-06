// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dyncfg

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

// A Dyncfg is a dynamic configuration key lookup mechanism.
type Dyncfg struct {
	cache         map[string]cacheEntry
	buildFileList func(map[string]string) []string
	l             sync.Mutex
}

// New takes a file list builder and returns a new instance of Dyncfg.
//
// The instance is ready to use.
func New(f func(map[string]string) []string) *Dyncfg {
	var c Dyncfg
	c.cache = make(map[string]cacheEntry)
	c.buildFileList = f
	return &c
}

func (c *Dyncfg) lookupvar(key, path string) interface{} {
	var f interface{}
	i, err := os.Stat(path)
	_, ok := c.cache[path]
	if os.IsNotExist(err) {
		log.Println("Dyncfg does not exist", path)
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
			log.Println("Purging", path, "from cache:", err)
			delete(c.cache, path)
			return nil
		}

		err = json.Unmarshal(data, &f)
		if err != nil {
			log.Println("Cannot parse", path)
			log.Println("Purging", path, "from cache:", err)
			delete(c.cache, path)
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

// Lookup searches for requested configuration key in file list built using
// context.
func (c *Dyncfg) Lookup(context map[string]string, key string) interface{} {
	var value interface{}

	c.l.Lock()
	defer c.l.Unlock()

	for _, fpath := range c.buildFileList(context) {
		log.Println("Context:", context, "Looking up", key, "in", fpath)
		value = c.lookupvar(key, fpath)
		if value != nil {
			log.Println("Context:", context, "Found", key, value)
			return value
		}
	}

	log.Println("Context:", context, "Key", key, "not found")
	return nil
}

// LookupString is analogous to Lookup(), but does the cast to string.
func (c *Dyncfg) LookupString(context map[string]string, key string) string {
	var value interface{}

	c.l.Lock()
	defer c.l.Unlock()

	for _, fpath := range c.buildFileList(context) {
		log.Println("Context:", context, "Looking up", key, "in", fpath)
		value = c.lookupvar(key, fpath)
		if value != nil {
			log.Println("Context:", context, "Found", key, value)
			return value.(string)
		}
	}

	log.Println("Context:", context, "Key", key, "not found")
	return ""
}

// LookupInt is analogous to Lookup(), but does the cast to int.
func (c *Dyncfg) LookupInt(context map[string]string, key string) int {
	var value interface{}

	c.l.Lock()
	defer c.l.Unlock()

	for _, fpath := range c.buildFileList(context) {
		log.Println("Context:", context, "Looking up", key, "in", fpath)
		value = c.lookupvar(key, fpath)
		if value != nil {
			log.Println("Context:", context, "Found", key, value)
			return int(value.(float64))
		}
	}

	log.Println("Context:", context, "Key", key, "not found")
	return -1
}

// LookupStringSlice is analogous to Lookup(), but does the cast to []string
func (c *Dyncfg) LookupStringSlice(context map[string]string, key string) (retval []string) {
	var value interface{}

	c.l.Lock()
	defer c.l.Unlock()

	for _, fpath := range c.buildFileList(context) {
		log.Println("Context:", context, "Looking up", key, "in", fpath)
		value = c.lookupvar(key, fpath)
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
func (c *Dyncfg) LookupStringMap(context map[string]string, key string) (retval map[string]bool) {
	var value interface{}
	retval = make(map[string]bool)

	c.l.Lock()
	defer c.l.Unlock()

	for _, fpath := range c.buildFileList(context) {
		log.Println("Context:", context, "Looking up", key, "in", fpath)
		value = c.lookupvar(key, fpath)
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
