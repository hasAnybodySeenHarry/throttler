package data

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type UsersModel struct {
	client *redis.Client
}

func (m *UsersModel) GetUserIDByToken(authTkn string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	userID, err := m.client.Get(ctx, authTkn).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}

	return userID, nil
}

func (m *UsersModel) InsertTokenWithID(token string, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	expiration := 24 * time.Hour

	err := m.client.Set(ctx, token, userID, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}
