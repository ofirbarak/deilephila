package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
)

type documentKey struct {
	ID primitive.ObjectID `bson:"_id"`
}

type changeID struct {
	Data string `bson:"_data"`
}

type namespace struct {
	Db   string `bson:"db"`
	Coll string `bson:"coll"`
}

type ChangeEvent struct {
	ID            changeID            `bson:"_id"`
	OperationType string              `bson:"operationType"`
	ClusterTime   primitive.Timestamp `bson:"clusterTime"`
	FullDocument  bson.D              `bson:"fullDocument"`
	DocumentKey   documentKey         `bson:"documentKey"`
	Ns            namespace           `bson:"ns"`
}

type ErrorEvent struct {
	Event ChangeEvent
	Err   error
}

type MongoOptions struct {
	MongoURI string `yaml:"URI"`
	Database string `yaml:"database"`
}

type MongoDriver struct {
	Config *MongoOptions
	Driver *mongo.Client
}

func (d *MongoDriver) Init(config []byte) {
	mongoConfig := new(MongoOptions)
	if err := yaml.Unmarshal(config, mongoConfig); err != nil {
		log.Fatal(err)
	}

	options := options.Client().
		ApplyURI(mongoConfig.MongoURI).
		SetConnectTimeout(1 * time.Second)

	client, err := mongo.Connect(context.Background(), options)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected successfully to MongoDB")

	d.Config = mongoConfig
	d.Driver = client
}

func (d *MongoDriver) Read(ctx context.Context, eventCh chan ChangeEvent) error {
	// TODO: support custom aggregation pipeline watch
	options := options.ChangeStreamOptions{}
	options.SetBatchSize(100)
	// options.SetFullDocument(options.WhenAvailable)

	// We ignore other operation types because the community edition does not support more than one database
	// changeStream, err := d.Driver.Database(d.Config.Database).Watch(
	// 	context.TODO(),
	// 	mongo.Pipeline{{{"$match", bson.D{{"operationType", bson.D{{"$in", []string{"insert", "delete", "update", "invalidate"}}}}}}}})

	changeStream, err := d.Driver.Watch(ctx, mongo.Pipeline{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start listening to events...")

	defer changeStream.Close(ctx)
	for changeStream.Next(ctx) {
		var event ChangeEvent
		if err := changeStream.Decode(&event); err != nil {
			log.Print(err)
		}

		log.Println("push event", event)
		eventCh <- event
	}

	return changeStream.Err()
}
