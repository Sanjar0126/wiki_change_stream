package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Sanjar0126/wiki_change_stream/config"
	"github.com/Sanjar0126/wiki_change_stream/event"
	"github.com/Sanjar0126/wiki_change_stream/models"
	"github.com/Sanjar0126/wiki_change_stream/pkg/helper"
	"github.com/Sanjar0126/wiki_change_stream/storage"
	"github.com/Sanjar0126/wiki_change_stream/storage/db"
)

func processEvents[T any](ctx context.Context, storage storage.StorageI,
	eventChan <-chan T, processor func(storage.StorageI, T), wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Processor %d: Context cancelled, stopping processor")
			return
		case event, ok := <-eventChan:
			if !ok {
				log.Printf("Processor %d: Channel closed, stopping processor")
				return
			}

			processor(storage, event)
		}
	}
}

func pushToDB(storage storage.StorageI, e models.WikiRecentChanges) {
	e.ServerPrefix = helper.GetPrefixFromServerName(e.ServerName)
	_, err := storage.WikiChanges().Create(context.Background(), e)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var wg sync.WaitGroup

	cfg := config.Load()

	dbConn := db.NewConn(&cfg)
	defer func() {
		err := dbConn.MongoConn.Client().Disconnect(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	storageDB := storage.New(dbConn.MongoConn)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan := make(chan models.WikiRecentChanges)
	go processEvents(ctx, storageDB, eventChan, pushToDB, &wg)
	go event.ConsumeEvents(ctx, config.EventStreamURL, eventChan, &wg)

	sig := <-sigChan
	log.Printf("Received signal: %v", sig)
	log.Println("Starting shutdown...")
	cancel()

	//wait for all goroutines to finish
	wg.Wait()
	log.Println("Shutdown complete")
}
