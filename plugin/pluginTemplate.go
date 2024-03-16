package main

import (
	"github.com/ofirbarak/deilephila"
	"github.com/ofirbarak/deilephila/drivers/mongodb"
)

func Map(event mongodb.ChangeEvent, src deilephila.Driver, dst deilephila.Driver) error {
	return nil
}
