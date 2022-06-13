package data

import "go.mongodb.org/mongo-driver/mongo"

type Data struct {
	DB *mongo.Client
}

var MongoCN Data

func init() {
	MongoCN.DB = ConnectionDB()
}
