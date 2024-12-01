package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client     *mongo.Client
	Collection *mongo.Collection
)

func InitMongoDB() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb://mongoadmin:secret@localhost:27017/")
	Client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	Collection = Client.Database("coupon_management").Collection("coupons")
}
