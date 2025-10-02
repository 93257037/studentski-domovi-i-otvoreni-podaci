package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB - drzi konekciju sa bazom podataka
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// kreira novu MongoDB konekciju sa prosledjenim URI-jem i imenom baze
// testira konekciju i vraca MongoDB instancu
func NewMongoDB(uri, databaseName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB!")

	database := client.Database(databaseName)

	return &MongoDB{
		Client:   client,
		Database: database,
	}, nil
}

// zatvara MongoDB konekciju
func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return m.Client.Disconnect(ctx)
}

// vraca kolekciju iz baze podataka po imenu
func (m *MongoDB) GetCollection(collectionName string) *mongo.Collection {
	return m.Database.Collection(collectionName)
}

