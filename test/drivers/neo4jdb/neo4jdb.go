package neo4jdb

import (
	"context"
	"errors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/ofirbarak/deilephila/drivers/neo4jdb"
)

type TestNeo4jDriver struct {
	neo4jdb.Neo4jDriver
}

func Neo4jConfig() *neo4jdb.Neo4jOptions {
	return &neo4jdb.Neo4jOptions{
		URI:      "neo4j://localhost:7687",
		Username: "neo4j",
		Password: "neo4j",
		Database: "neo4j",
	}
}

const COUNT_ALL_NODES = "MATCH (n) RETURN count(n)"

func (d *TestNeo4jDriver) Init(config []byte) {
	d.Neo4jDriver.Init(config)
}

func (d *TestNeo4jDriver) Count(filter string) (int64, error) {
	neo4j_result, err := neo4j.ExecuteQuery(context.TODO(), d.Driver, filter, map[string]any{}, neo4j.EagerResultTransformer)
	if err != nil {
		return -1, err
	}

	if !(len(neo4j_result.Records) == 1 && len(neo4j_result.Records[0].Values) == 1) {
		return -1, errors.New("can't parse results")
	}

	return neo4j_result.Records[0].Values[0].(int64), nil
}

func (d *TestNeo4jDriver) CleanDatabase() error {
	_, err := neo4j.ExecuteQuery(
		context.TODO(),
		d.Driver,
		"MATCH (n) DETACH DELETE n",
		map[string]any{},
		neo4j.EagerResultTransformer)

	return err
}
