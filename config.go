package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Config struct {
	ConfigDir string
}

type Context struct {
	Network string
	Target  string
	Source  string
}

var C Config

func lookupvar(key, path string) interface{} {
	var f interface{}
	_, err := os.Stat(path)
	if err != nil {
		return nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(data, &f)

	return f.(map[string]interface{})[key]
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
