package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

type cacheEntry struct {
	ModTime  time.Time
	Contents interface{}
}

type Config struct {
	ConfigDir string
	Cache     map[string]cacheEntry
}

type Context struct {
	Network string
	Target  string
	Source  string
}

var C Config

func init() {
	log.Println("Initializing configuration cache")
	C.Cache = make(map[string]cacheEntry)
}

func lookupvar(key, path string) interface{} {
	var f interface{}
	i, err := os.Stat(path)
	_, ok := C.Cache[path]
	if os.IsNotExist(err) {
		log.Println("Config does not exist", path)
		if ok {
			log.Println("Purging", path, "from cache")
			delete(C.Cache, path)
		}
	}
	if err != nil {
		return nil
	}

	if C.Cache[path].ModTime.Before(i.ModTime()) || !ok {
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
		C.Cache[path] = cacheEntry{
			ModTime:  i.ModTime(),
			Contents: f,
		}

		return f.(map[string]interface{})[key]
	} else {
		log.Println("Found cache for", path)
		return C.Cache[path].Contents.(map[string]interface{})[key]
	}
}

func (c *Config) Lookup(context Context, key string) interface{} {
	var fpath string
	var value interface{}
	if context.Network != "" {
		if context.Source != "" {
			if context.Target != "" {
				fpath = path.Join(c.ConfigDir, context.Network, context.Source, context.Target+".json")

				log.Println("Context:", context, "Looking up", key, "in", fpath)
				value = lookupvar(key, fpath)
				if value != nil {
					log.Println("Context:", context, "Found", key, value)
					return value
				}
			}

			fpath = path.Join(c.ConfigDir, context.Network, context.Source+".json")

			log.Println("Context:", context, "Looking up", key, "in", fpath)
			value = lookupvar(key, fpath)
			if value != nil {
				log.Println("Context:", context, "Found", key, value)
				return value
			}
		}

		fpath = path.Join(c.ConfigDir, context.Network+".json")

		log.Println("Context:", context, "Looking up", key, "in", fpath)
		value = lookupvar(key, fpath)
		if value != nil {
			log.Println("Context:", context, "Found", key, value)
			return value
		}
	}

	fpath = path.Join(c.ConfigDir, "common.json")

	log.Println("Context:", context, "Looking up", key, "in", fpath)
	value = lookupvar(key, fpath)
	if value != nil {
		log.Println("Context:", context, "Found", key, value)
		return value
	} else {
		log.Println("Key", key, "not found")
		return nil
	}
}
