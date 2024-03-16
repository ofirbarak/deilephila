package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"plugin"

	"github.com/ofirbarak/deilephila"
	"github.com/ofirbarak/deilephila/drivers/mongodb"
)

const PLUGIN_FUNC_NAME = "Map"

func getFuncFromPlugin(soPath string) plugin.Symbol {
	plugin, err := plugin.Open(soPath)
	if err != nil {
		log.Fatal(err)
	}

	f, err := plugin.Lookup(PLUGIN_FUNC_NAME)
	if err != nil {
		log.Fatal(err)
	}

	return f
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Wrong usage", os.Args)
		fmt.Println("Usage: deilephila <path-to-config.yml> <path-to-plugin.so>")
		os.Exit(1)
	}

	filepath := os.Args[1]
	soPath := os.Args[2]

	f := getFuncFromPlugin(soPath)
	config := deilephila.ReadConfig(filepath)
	config.Sync.MapFunction = f.(func(event mongodb.ChangeEvent, srcDriver deilephila.Driver, dstDriver deilephila.Driver) error)

	deilephila.New(config).Run(context.Background())
}
