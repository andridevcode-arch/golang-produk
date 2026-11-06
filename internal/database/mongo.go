package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go-produk/internal/config"
)

// Mongo adalah struct untuk koneksi ke MongoDB
type Mongo struct {
	client *mongo.Client
	db     *mongo.Database
	cfg    config.MongoConfig
}

// NewMongo membuat koneksi baru ke MongoDB
func NewMongo(ctx context.Context, cfg config.MongoConfig) (*Mongo, error) {
	clientOpts := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetServerSelectionTimeout(cfg.Timeout)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("mongo ping: %w", err)
	}

	db := client.Database(cfg.Database)

	// Buat indeks yang dibutuhkan collection (idempoten).
	if err := EnsureIndexes(context.Background(), db); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("mongo ensure indexes: %w", err)
	}

	return &Mongo{
		client: client,
		db:     db,
		cfg:    cfg,
	}, nil
}

// Close menutup koneksi
func (m *Mongo) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// Database mengembalikan database untuk digunakan
func (m *Mongo) Database() *mongo.Database {
	return m.db
}

// Client mengembalikan MongoDB client untuk digunakan
func (m *Mongo) Client() *mongo.Client {
	return m.client
}

// Ping test koneksi
func (m *Mongo) Ping(ctx context.Context) error {
	return m.client.Ping(ctx, readpref.Primary())
}
