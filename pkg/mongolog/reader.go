package mongolog

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogEntry struct {
	ID        interface{}            `bson:"_id" json:"id"`
	CreatedAt time.Time              `bson:"created_at" json:"-"`
	Details   map[string]interface{} `bson:"details" json:"details"`
	Level     string                 `bson:"level" json:"-"`
	Msg       string                 `bson:"msg" json:"-"`
	LogTime   string                 `bson:"time" json:"logTime"`
}

func (l *logStore) ReadLogsPaginated(ctx context.Context, limit, offset int64) ([]LogEntry, error) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if limit < 1 {
		limit = 10
	}

	findOptions := options.Find()
	findOptions.SetSkip(offset)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{Key: "timestamp", Value: -1}}) // Sort by timestamp descending

	cursor, err := l.collection.Find(ctxWithTimeout, bson.D{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch logs: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []LogEntry
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("failed to decode logs: %w", err)
	}

	return logs, nil
}
