package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type M map[string]interface{}

// if not call make, options.Set will complain "assignment to entry in nil map"
var options M = make(M)

func LoadOptions(configFile string) *M {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("error reading file: %v\n", err)
	} else if err = json.Unmarshal([]byte(data), &options); err != nil {
		log.Printf("error parsing conf: %s, %v\n", configFile, err)
	}

	confValue := os.Getenv("CONF")
	if confValue != "" {
		if err = json.Unmarshal([]byte(confValue), &options); err != nil {
			log.Printf("error parsing CONF: %v\n", err)
		}
	}
	return &options
}

func GetOptions() *M {
	return &options
}

func (self *M) Get(name string) interface{} {
	return (*self)[name]
}

func (self *M) Set(name string, value interface{}) {
	(*self)[name] = value
}

func (self *M) SetDefault(name string, value interface{}) {
	if _, ok := (*self)[name]; !ok {
		(*self)[name] = value
	}
}

func (self *M) GetBoolean(name string) bool {
	if value, ok := (*self)[name].(bool); ok {
		return value
	}
	return false
}

func (self *M) GetString(name string) string {
	if value, ok := (*self)[name].(string); ok {
		return value
	}
	return ""
}

func (self *M) GetStringSlice(name string) []string {
	// for options.Set
	if value, ok := (*self)[name].([]string); ok {
		return value
	}
	return []string{}
}

func (self *M) GetInt(name string) int {
	if value, ok := (*self)[name].(int); ok {
		return value
	}
	return 0
}

func (self *M) GetMap(name string) *M {
	if value, ok := (*self)[name].(M); ok {
		return &value
	}
	return nil
}
