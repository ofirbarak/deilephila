package benchmark

import (
	"context"
	"io/ioutil"
	"log"
	"testing"
	"time"

	. "github.com/ofirbarak/deilephila/test"
	. "github.com/ofirbarak/deilephila/test/drivers/mongodb"
	. "github.com/ofirbarak/deilephila/test/drivers/neo4jdb"

	"github.com/ofirbarak/deilephila"
)

type BenchmarkTestSuite struct {
	config  *deilephila.Config
	src     *TestMongoDriver
	dst     *TestNeo4jDriver
	program *deilephila.Deilephila
	ctx     context.Context
	cancel  context.CancelFunc
}

func setupSuite() *BenchmarkTestSuite {
	log.SetOutput(ioutil.Discard)
	s := BenchmarkTestSuite{
		config: &deilephila.Config{
			Src: MongoTestConfig(),
			Dst: Neo4jConfig(),
			Sync: deilephila.SyncOptions{
				Workers:     5,
				MapFunction: Map,
			},
		},
		src: new(TestMongoDriver),
		dst: new(TestNeo4jDriver),
	}

	s.src.Init(Marshal(s.config.Src))
	s.dst.Init(Marshal(s.config.Dst))

	ctx, cancel := context.WithCancel(context.Background())
	s.ctx, s.cancel = ctx, cancel

	go deilephila.New(s.config).Run(ctx)
	time.Sleep(time.Second)

	return &s
}

func tearDownSuite(s *BenchmarkTestSuite) {
	s.src.Base.Driver.Disconnect(context.TODO())
	s.dst.Driver.Close(context.TODO())
}

func tearDownTest(s *BenchmarkTestSuite) {
	s.cancel()
	s.src.CleanDatabase()
	s.dst.CleanDatabase()
}

func runTest(b *testing.B, s *BenchmarkTestSuite) {
	defer tearDownTest(s)

	numOfDocuments := 10
	docs := make([]interface{}, numOfDocuments)
	for i := 0; i < numOfDocuments; i++ {
		docs[i] = GenerateDoc()
	}

	b.StartTimer()
	if _, err := s.src.Collection.InsertMany(context.TODO(), docs); err != nil {
		b.Fatal(err)
	}

	err := WaitForCondition(func() bool {
		res, _ := s.dst.Count(COUNT_ALL_NODES)
		return res == int64(numOfDocuments)
	}, time.Second)

	if err != nil {
		b.Fatal(err)
	}

	b.StopTimer()
}

func BenchmarkMongoToNeo4j(b *testing.B) {
	suite := setupSuite()
	defer tearDownSuite(suite)

	for i := 0; i < b.N; i++ {
		runTest(b, suite)
	}
}
