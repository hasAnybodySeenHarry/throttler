package data

import (
	"github.com/redis/go-redis/v9"
)

const (
	luaScript = `
		local tokens_key = KEYS[1]
		local timestamp_key = KEYS[2]

		local bucket_size = tonumber(ARGV[1])
		local refill_rate = tonumber(ARGV[2])
		local refill_period = tonumber(ARGV[3])

		local now = tonumber(ARGV[4])
		local tokens = tonumber(redis.call("get", tokens_key))
		local last_refill = tonumber(redis.call("get", timestamp_key))

		if tokens == nil then
			tokens = bucket_size
			last_refill = now
		end

		local elapsed = now - last_refill
		local refills = math.floor(elapsed / refill_period)
		tokens = math.min(bucket_size, tokens + refills * refill_rate)

		if tokens > 0 then
			tokens = tokens - 1
			redis.call("set", tokens_key, tokens)
			redis.call("set", timestamp_key, now)
			return 1
		else
			return 0
		end
	`
)

type Models struct {
	Buckets BucketsModel
}

func New(client *redis.Client) Models {
	return Models{
		Buckets: BucketsModel{
			client: client,
			script: redis.NewScript(luaScript),
		},
	}
}
