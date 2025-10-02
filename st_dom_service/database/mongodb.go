package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB - predstavlja konekciju sa MongoDB bazom podataka
type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

// kreira novu MongoDB konekciju sa prosledjenim URI-jem i imenom baze
// testira konekciju i vraca MongoDB instancu
func NewMongoDB(uri, dbName string) (*MongoDB, error) {
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

	database := client.Database(dbName)

	return &MongoDB{
		client:   client,
		database: database,
	}, nil
}

// vraca kolekciju iz baze podataka po imenu
func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

// zatvara MongoDB konekciju
func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	return m.client.Disconnect(ctx)
}

// vraca instancu baze podataka
func (m *MongoDB) GetDatabase() *mongo.Database {
	return m.database
}
