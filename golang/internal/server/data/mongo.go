package data

import (
	"context"
	"log"

	"github.com/chipocrudos/microblog/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectionDB() *mongo.Client {

	clientOpts := options.Client().ApplyURI(config.Config.MONGO_URI)
	client, err := mongo.Connect(context.TODO(), clientOpts)

	if err != nil {
		log.Panic(err.Error())
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	return client
}

func (db *Data) PingDB() bool {

	err := db.DB.Ping(context.TODO(), nil)
	if err != nil {
		return false
	}
	return true
}

func (db *Data) GetCollection(collectionName string) *mongo.Collection {
	collection := db.DB.Database(config.Config.MONGO_DB).Collection(collectionName)
	return collection
}
