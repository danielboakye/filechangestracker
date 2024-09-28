package mongolog

import (
	"context"
	"fmt"
	"time"

	"log/slog"

	"github.com/danielboakye/filechangestracker/pkg/config"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -destination=../../mocks/mongolog/mock_mongolog.go -package=mongologmock -source=mongolog.go
type LogStore interface {
	Write(p []byte) (n int, err error)
	Close(ctx context.Context) error
	ReadLogsPaginated(ctx context.Context, page, pageSize int64) ([]LogEntry, error)
}

type logStore struct {
	collection *mongo.Collection
}

// NewLogger creates a new slog.Logger that writes to MongoDB
func NewLogger(mongoURI string) (*slog.Logger, LogStore, error) {
	mongoWriter, err := newMongoWriter(mongoURI, config.LogsDBName, config.LogsCollectionName)
	if err != nil {
		return nil, nil, err
	}
	handler := slog.NewJSONHandler(mongoWriter, &slog.HandlerOptions{})
	return slog.New(handler), mongoWriter, nil
}

func newMongoWriter(mongoURI, databaseName, collectionName string) (LogStore, error) {
	ctx := context.TODO()
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	db := client.Database(databaseName)
	collections, err := db.ListCollectionNames(ctx, map[string]interface{}{"name": collectionName})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}

	if len(collections) == 0 {
		err := db.CreateCollection(ctx, collectionName)
		if err != nil {
			return nil, fmt.Errorf("failed to create collection: %w", err)
		}
	}

	return &logStore{
		collection: db.Collection(collectionName),
	}, nil
}

// Write implements the io.Writer interface
func (l *logStore) Write(p []byte) (n int, err error) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	var logEntry LogEntry
	err = bson.UnmarshalExtJSON(p, true, &logEntry)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal log entry: %w", err)
	}

	logEntry.ID = uuid.NewString()
	logEntry.CreatedAt = time.Now()
	_, err = l.collection.InsertOne(ctxWithTimeout, logEntry)
	if err != nil {
		return 0, fmt.Errorf("failed to insert log entry into MongoDB: %w", err)
	}

	return len(p), nil
}

func (l *logStore) Close(ctx context.Context) error {
	return l.collection.Database().Client().Disconnect(ctx)
}
