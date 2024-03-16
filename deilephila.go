package deilephila

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ofirbarak/deilephila/drivers/mongodb"
	"github.com/ofirbarak/deilephila/drivers/neo4jdb"
)

type Deilephila struct {
	config    *Config
	srcDriver Driver
	dstDriver Driver
	workersWg sync.WaitGroup
	eventCh   chan mongodb.ChangeEvent
	errorCh   chan mongodb.ErrorEvent
}

func (t *Deilephila) writerWorker(ctx context.Context) {
	defer t.workersWg.Done()

	for {
		select {
		case event := <-t.eventCh:
			// TODO: add context timeout for event
			if err := t.config.Sync.MapFunction(event, t.srcDriver, t.dstDriver); err != nil {
				t.errorCh <- mongodb.ErrorEvent{Event: event, Err: err}
			}
		case <-ctx.Done():
			log.Println("Closing writer worker", ctx.Err())
			return
		default:
			// TODO: Remove and find another way to iterate both channels
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (t *Deilephila) readerWorker(ctx context.Context) {
	defer t.workersWg.Done()
	defer close(t.eventCh)

	err := t.srcDriver.Read(ctx, t.eventCh)

	log.Println("Closing reader worker", err)
}

func (t *Deilephila) errorWorker(ctx context.Context) {
	defer t.workersWg.Done()
	defer close(t.errorCh) // TODO: Move to workers or upper function

	for {
		select {
		case event := <-t.errorCh:
			log.Println("error", event)
		case <-ctx.Done():
			log.Println("Closing error worker", ctx.Err())
			return
		default:
			// TODO: Remove and find another way to iterate both channels
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (t *Deilephila) Run(ctx context.Context) {
	log.Println("Connecting drivers...")
	t.srcDriver.Init(marshal(t.config.Src))
	t.dstDriver.Init(marshal(t.config.Dst))

	log.Println("Running...")
	t.workersWg.Add(t.config.Sync.Workers)
	for i := 0; i < t.config.Sync.Workers; i++ {
		go t.writerWorker(ctx)
	}

	t.workersWg.Add(1)
	go t.readerWorker(ctx)

	t.workersWg.Add(1)
	go t.errorWorker(ctx)

	t.workersWg.Wait()
}

func New(config *Config) *Deilephila {
	// TODO: add validation & defaults to config
	log.Printf("Using config: %+v\n", config)

	t := Deilephila{
		config:    config,
		eventCh:   make(chan mongodb.ChangeEvent),
		errorCh:   make(chan mongodb.ErrorEvent),
		srcDriver: new(mongodb.MongoDriver),
		dstDriver: new(neo4jdb.Neo4jDriver),
	}

	return &t
}
