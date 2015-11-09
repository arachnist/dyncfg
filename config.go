package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Config struct {
	Nick      string
	Host      string
	RealName  string
	User      string
	Networks  []string
	Servers   map[string][]string
	Channels  map[string][]string
	Passwords map[string]string
	Plugins   []string
	Ignore    []string
	Logpath   string
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
		return ""
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
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
					return value
				}
			}
			return lookupvar(key, path.Join(c.ConfigDir, context.Network, context.Source+".json"))
		}
		return lookupvar(key, path.Join(c.ConfigDir, context.Network+".json"))
	}
	return lookupvar(key, path.Join(c.ConfigDir, "common.json"))
}
