package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect se koristi za povezivanje sa MongoDB bazom
func Connect(uri string, dbName string) (*mongo.Database, error) {
	// kontekst sa timeout-om da ne blokira aplikaciju ako nema konekcije
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// pokušaj konekcije
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// testiraj konekciju
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// vrati bazu
	return client.Database(dbName), nil
}
