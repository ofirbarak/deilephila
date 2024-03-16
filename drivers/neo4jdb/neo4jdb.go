package neo4jdb

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/ofirbarak/deilephila/drivers/mongodb"
	"gopkg.in/yaml.v3"
)

type Neo4jOptions struct {
	URI      string `yaml:"URI"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Neo4jDriver struct {
	Config *Neo4jOptions
	Driver neo4j.DriverWithContext
}

func (d *Neo4jDriver) Init(config []byte) {
	neo4jConfig := new(Neo4jOptions)
	if err := yaml.Unmarshal(config, neo4jConfig); err != nil {
		log.Fatal(err)
	}

	driver, err := neo4j.NewDriverWithContext(
		neo4jConfig.URI,
		neo4j.BasicAuth(neo4jConfig.Username, neo4jConfig.Password, ""))
	if err != nil {
		log.Fatal(err)
	}

	if err := driver.VerifyConnectivity(context.Background()); err != nil {
		log.Fatal(err)
	}

	// TODO: read version and if pro create the DB
	// Database must be created manually since the community edition does not support it
	// res2, err2 := neo4j.ExecuteQuery(context.TODO(), driver, "CREATE DATABASE test", nil, neo4j.EagerResultdeilephila)
	// if err2 != nil {
	// log.Fatal(err2)
	// }

	log.Println("Connected successfully to Neo4j")

	d.Config = neo4jConfig
	d.Driver = driver
}

func (d *Neo4jDriver) Read(context.Context, chan mongodb.ChangeEvent) error {
	panic("Not implemented")
}
