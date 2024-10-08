package data

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Bucket struct {
	bucketSize   int
	refillRate   int
	refillPeriod time.Duration
}

type BucketsModel struct {
	client *redis.Client
	script *redis.Script
}

func NewBucket(hasPrivileges bool) *Bucket {
	if hasPrivileges {
		return &Bucket{
			bucketSize:   8,
			refillRate:   1,
			refillPeriod: 3 * time.Second,
		}
	}

	return &Bucket{
		bucketSize:   5,
		refillRate:   1,
		refillPeriod: 4 * time.Second,
	}
}

func (m *BucketsModel) Allow(key string, activated bool) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tokensKey := key + ":tokens"
	timestampKey := key + ":timestamp"

	now := time.Now().Unix()
	expiration := 24 * 3600

	b := NewBucket(activated)

	args := []interface{}{b.bucketSize, b.refillRate, b.refillPeriod.Seconds(), now, expiration}

	result, err := m.script.Run(ctx, m.client, []string{tokensKey, timestampKey}, args...).Int64()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}
