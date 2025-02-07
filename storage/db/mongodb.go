package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log"

	"github.com/Sanjar0126/wiki_change_stream/config"
)

type DB struct {
	MongoConn *mongo.Database
}

func NewConn(cfg *config.Config) *DB {
	var (
		mongoConn *mongo.Client
		err       error
	)

	credential := options.Credential{
		AuthSource:    cfg.MongoDBDatabase,
		Username:      cfg.MongoDBUser,
		Password:      cfg.MongoDBPassword,
		AuthMechanism: "SCRAM-SHA-256",
	}

	mongoString := fmt.Sprintf("mongodb://%s:%d", cfg.MongoDBHost, cfg.MongoDBPort)

	log.Println("Connecting:", mongoString)

	ctx := context.Background()

	if cfg.Environment == "develop" {
		mongoConn, err = mongo.Connect(
			ctx,
			options.Client().ApplyURI(mongoString),
		)
	} else {
		mongoConn, err = mongo.Connect(
			ctx,
			options.Client().ApplyURI(mongoString).SetAuth(credential),
		)
	}

	if err != nil {
		log.Printf("failed to connecto to mongodb, %v", err)
		panic(err)
	}

	if err := mongoConn.Ping(context.TODO(), nil); err != nil {
		log.Printf("failed to connecto to database, %v", err)
		panic(err)
	}

	connDB := mongoConn.Database(cfg.MongoDBDatabase)
	log.Printf("Connected to MongoDBdatabase: %v", connDB.Name())

	return &DB{
		MongoConn: connDB,
	}
}
