package deilephila

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/ofirbarak/deilephila/drivers/mongodb"
	"gopkg.in/yaml.v3"
)

type Driver interface {
	Init(config []byte)
	Read(context.Context, chan mongodb.ChangeEvent) error
}

type Mapper func(event mongodb.ChangeEvent, srcDriver Driver, dstDriver Driver) error

type SyncOptions struct {
	Workers     int `yaml:"workers"`
	MapFunction Mapper
}

type Config struct {
	Src  interface{} `yaml:"src"`
	Dst  interface{} `yaml:"dst"`
	Sync SyncOptions `yaml:"sync"`
}

// Parse yaml file to Config struct
func ReadConfig(configPath string) *Config {
	filepath, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer filepath.Close()

	data, err := ioutil.ReadAll(filepath)
	if err != nil {
		log.Fatal(err)
	}

	config := new(Config)
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	return config
}

func marshal(config interface{}) []byte {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
