package config

import (
    "encoding/json"
    "io/ioutil"
)

type Config struct {
    Nick        string
    Host        string
    Networks    []string
    Servers     map[string] []string
    Channels    map[string] []string
    Passwords   map[string] string
    Plugins     []string
    Ignore      []string
    Logpath     string
}

func ReadConfig(path string) (Config, error) {
    var config Config
    
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return config, err
    }

    err = json.Unmarshal(data, &config)
    if err != nil {
        return config, err
    }

    return config, nil
}
