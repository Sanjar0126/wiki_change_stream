package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Sanjar0126/wiki_change_stream/config"
	"github.com/Sanjar0126/wiki_change_stream/discord"
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
			log.Printf("Processor : Context cancelled, stopping processor")
			return
		case event, ok := <-eventChan:
			if !ok {
				log.Printf("Processor : Channel closed, stopping processor")
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
		log.Println(err)
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

	eventConfig := event.ConsumerConfig{
		URL:            config.EventStreamURL,
		ReconnectDelay: time.Second * 5,
		GetLatestTimestamp: func() string {
			return storageDB.WikiChanges().GetLatest()
		},
		MaxRetries: 3, // 0 for infinite retries
	}

	wg.Add(2)

	go processEvents(ctx, storageDB, eventChan, pushToDB, &wg)
	go event.ConsumeEvents(ctx, eventConfig, eventChan, &wg)

	discordHander := discord.NewHandler(&discord.HandlerOptions{
		Config: &cfg,
		DB:     storageDB,
	})

	discord := discord.NewDiscord(&cfg, discordHander)

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := discord.Open(); err != nil {
			log.Printf("Error opening connection: %v", err)
			cancel()

			return
		}

		<-ctx.Done()

		if err := discord.Close(); err != nil {
			log.Printf("Error closing Discord connection: %v", err)
		}
	}()

	sig := <-sigChan
	log.Printf("Received signal: %v", sig)
	log.Println("Starting shutdown...")
	cancel()

	//wait for all goroutines to finish
	wg.Wait()
	log.Println("Shutdown complete")
}
