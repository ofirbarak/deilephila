package test

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	"gopkg.in/yaml.v3"
)

func Marshal(config interface{}) []byte {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func WaitForCondition(condition func() bool, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		if condition() {
			return nil
		}

		if ctx.Err() != nil {
			return errors.New("timeout! condition is not met")
		}

		time.Sleep(50 * time.Millisecond)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateDoc() TestDocument {
	return TestDocument{Name: RandString(10)}
}
