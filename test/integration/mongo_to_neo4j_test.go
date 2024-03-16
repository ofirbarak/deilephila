package integration

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ofirbarak/deilephila"
	. "github.com/ofirbarak/deilephila/test"
	. "github.com/ofirbarak/deilephila/test/drivers/mongodb"
	. "github.com/ofirbarak/deilephila/test/drivers/neo4jdb"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	config  *deilephila.Config
	src     *TestMongoDriver
	dst     *TestNeo4jDriver
	program *deilephila.Deilephila
	ctx     context.Context
	cancel  context.CancelFunc
}

func (t *IntegrationTestSuite) SetupSuite() {
	t.src.Init(Marshal(t.config.Src))
	t.dst.Init(Marshal(t.config.Dst))
}

func (t *IntegrationTestSuite) TearDownSuite() {
	t.src.Base.Driver.Disconnect(context.TODO())
	t.dst.Driver.Close(context.TODO())
}

func (t *IntegrationTestSuite) SetupTest() {
	ctx, cancel := context.WithCancel(context.Background())
	t.ctx, t.cancel = ctx, cancel

	go deilephila.New(t.config).Run(ctx)
	time.Sleep(time.Second)
}

func (t *IntegrationTestSuite) TearDownTest() {
	t.cancel()
	t.src.CleanDatabase()
	t.dst.CleanDatabase()
}

func (t *IntegrationTestSuite) TestInsert() {
	{
		_, err := t.src.Collection.InsertOne(context.TODO(), GenerateDoc())
		t.Assertions.Nil(err)
	}

	err := WaitForCondition(func() bool {
		res, err := t.dst.Count(COUNT_ALL_NODES)
		return res == 1 && err == nil
	}, 2*time.Second)
	t.Assertions.Nil(err)
}

func (t *IntegrationTestSuite) TestInsertInParallel() {
	numOfItemsToInsert := 10
	var waitGg sync.WaitGroup
	errorCh := make(chan error)

	for i := 0; i < numOfItemsToInsert; i++ {
		waitGg.Add(1)
		go func() {
			if _, err := t.src.Collection.InsertOne(context.TODO(), GenerateDoc()); err != nil {
				errorCh <- err
			}
			waitGg.Done()
		}()
	}
	waitGg.Wait()

	t.Assertions.Empty(errorCh)

	err := WaitForCondition(func() bool {
		res, err := t.dst.Count(COUNT_ALL_NODES)
		return res == int64(numOfItemsToInsert) && err == nil
	}, time.Second)
	t.Assertions.Nil(err)
}

func TestMongoToNeo4jTestSuite(t *testing.T) {
	testSuite := IntegrationTestSuite{
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
	suite.Run(t, &testSuite)
}
