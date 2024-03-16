package test

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/ofirbarak/deilephila"
	"github.com/ofirbarak/deilephila/drivers/mongodb"
	"github.com/ofirbarak/deilephila/drivers/neo4jdb"
	"go.mongodb.org/mongo-driver/bson"
)

type TestDocument struct {
	Name string `bson:"name"`
}

func Map(event mongodb.ChangeEvent, src deilephila.Driver, dst deilephila.Driver) error {
	if event.OperationType == "invalidate" {
		log.Fatal(errors.New("error operation type"))
	}

	if event.OperationType != "insert" {
		return nil
	}

	doc := new(TestDocument)
	{
		res, err := bson.Marshal(event.FullDocument)

		if err != nil {
			return err
		}
		if err := bson.Unmarshal(res, doc); err != nil {
			return err
		}

	}

	properties := map[string]interface{}{}
	{
		res, err := json.Marshal(doc)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(res, &properties); err != nil {
			return err
		}
	}

	dstDriver := dst.(*neo4jdb.Neo4jDriver)
	_, err := neo4j.ExecuteQuery(
		context.TODO(),
		dstDriver.Driver,
		"CREATE (n:Node $properties)",
		map[string]interface{}{"properties": properties},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dstDriver.Config.Database))

	return err
}
