package mongodb

import (
	"context"

	"github.com/ofirbarak/deilephila/drivers/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

func MongoTestConfig() *mongodb.MongoOptions {
	return &mongodb.MongoOptions{
		MongoURI: "mongodb://localhost:27017/?replicaSet=rs0&directConnection=true",
		Database: "test",
	}
}

type TestMongoDriver struct {
	Base       mongodb.MongoDriver
	Collection *mongo.Collection
}

func (d *TestMongoDriver) Init(config []byte) {
	d.Base.Init(config)
	d.Collection = d.Base.Driver.Database(d.Base.Config.Database).Collection("test")
}

func (d *TestMongoDriver) CleanDatabase() error {
	return d.Collection.Drop(context.TODO())
}
