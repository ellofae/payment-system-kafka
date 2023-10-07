package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ellofae/payment-system-kafka/payment-system/processor/internal/domain/entity"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client *redis.Client
}

func NewStroage(rdb *redis.Client) *Storage {
	return &Storage{
		client: rdb,
	}
}

func orderIDKey(id string) string {
	return fmt.Sprintf("transaction:%s", id)
}

func (s *Storage) Insert(ctx context.Context, model *entity.TransactionData) error {
	log := logger.GetLogger()

	data, err := json.Marshal(model)
	if err != nil {
		log.Error("Unable to marshall data", "error", err.Error())
		return err
	}
	key := orderIDKey(model.TransactionID)

	res := s.client.Set(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		log.Error("Failed to set data in redis", "error", err.Error())
		return err
	}

	log.Info("Successfully connected to Redis")

	return nil
}
