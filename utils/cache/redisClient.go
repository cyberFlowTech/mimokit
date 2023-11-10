package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/mapping"
	"strconv"
	"time"
)

type RedisClient struct {
	rc redis.UniversalClient
}
type (
	Pair struct {
		Key   string
		Score int64
	}
)

// 基于go-redis/v9 封装支持node/cluster mode的客户端
func NewRedisClient(adds []string) *RedisClient {
	return &RedisClient{
		rc: redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs: adds,
		}),
	}
}

func (c *RedisClient) Get(key string) (string, error) {
	return c.rc.Get(context.Background(), key).Result()
}

func (c *RedisClient) Mget(keys ...string) ([]string, error) {
	res := c.rc.MGet(context.Background(), keys...)

	if res.Err() != nil {
		return nil, res.Err()
	}

	return toStrings(res.Val()), res.Err()
}

func (c *RedisClient) Setex(key, value string, expire time.Duration) error {
	return c.rc.Set(context.Background(), key, value, expire).Err()
}

func toStrings(vals []any) []string {
	ret := make([]string, len(vals))

	for i, val := range vals {
		if val == nil {
			ret[i] = ""
			continue
		}

		switch val := val.(type) {
		case string:
			ret[i] = val
		default:
			ret[i] = mapping.Repr(val)
		}
	}

	return ret
}

func (c *RedisClient) HgetallCtx(ctx context.Context, key string) (map[string]string, error) {
	res := c.rc.HGetAll(ctx, key)
	if res.Err() != nil {
		return nil, res.Err()
	}
	return res.Val(), nil
}

func (c *RedisClient) HmsetCtx(ctx context.Context, key string, fieldsAndValues map[string]string) error {
	vals := make(map[string]any, len(fieldsAndValues))
	for k, v := range fieldsAndValues {
		vals[k] = v
	}
	return c.rc.HMSet(ctx, key, vals).Err()
}

func (c *RedisClient) ExpireCtx(ctx context.Context, key string, seconds int) error {
	return c.rc.Expire(ctx, key, time.Duration(seconds)*time.Second).Err()
}

func (c *RedisClient) HmgetCtx(ctx context.Context, key string, fields ...string) (val []string, err error) {
	v, err := c.rc.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return
	}
	val = toStrings(v)
	return
}
func (c *RedisClient) HdelCtx(ctx context.Context, key string, fields ...string) (val bool, err error) {
	v, err := c.rc.HDel(ctx, key, fields...).Result()
	if err != nil {
		return
	}

	val = v >= 1
	return
}

func (c *RedisClient) ScriptRun(script *redis.Script, keys []string, args ...any) (any, error) {
	return c.ScriptRunCtx(context.Background(), script, keys, args...)
}

func (c *RedisClient) ScriptRunCtx(ctx context.Context, script *redis.Script, keys []string, args ...any) (val any, err error) {

	val, err = script.Run(ctx, c.rc, keys, args...).Result()
	return
}

func (c *RedisClient) DelCtx(ctx context.Context, keys ...string) (val int, err error) {
	v, err := c.rc.Del(ctx, keys...).Result()
	if err != nil {
		return
	}

	val = int(v)
	return
}

func (c *RedisClient) GetCtx(ctx context.Context, key string) (val string, err error) {
	val, err = c.rc.Get(ctx, key).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (c *RedisClient) SetexCtx(ctx context.Context, key, value string, seconds int) error {

	return c.rc.Set(ctx, key, value, time.Duration(seconds)*time.Second).Err()
}

func (c *RedisClient) ZrangebyscoreWithScoresAndLimitCtx(ctx context.Context, key string, start,
	stop int64, page, size int) (val []Pair, err error) {
	if size <= 0 {
		return
	}
	v, err := c.rc.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min:    strconv.FormatInt(start, 10),
		Max:    strconv.FormatInt(stop, 10),
		Offset: int64(page * size),
		Count:  int64(size),
	}).Result()
	if err != nil {
		return
	}

	val = toPairs(v)
	return
}
func toPairs(vals []redis.Z) []Pair {
	pairs := make([]Pair, len(vals))
	for i, val := range vals {
		switch member := val.Member.(type) {
		case string:
			pairs[i] = Pair{
				Key:   member,
				Score: int64(val.Score),
			}
		default:
			pairs[i] = Pair{
				Key:   mapping.Repr(val.Member),
				Score: int64(val.Score),
			}
		}
	}
	return pairs
}

// ZaddCtx is the implementation of redis zadd command.
func (c *RedisClient) ZaddCtx(ctx context.Context, key string, score int64, value string) (
	val bool, err error) {
	return c.ZaddFloatCtx(ctx, key, float64(score), value)
}

func (c *RedisClient) ZaddFloatCtx(ctx context.Context, key string, score float64, value string) (
	val bool, err error) {
	v, err := c.rc.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: value,
	}).Result()
	if err != nil {
		return
	}
	val = v == 1
	return

}

func (c *RedisClient) RpushCtx(ctx context.Context, key string, values ...any) (val int, err error) {
	v, err := c.rc.RPush(ctx, key, values...).Result()
	if err != nil {
		return
	}
	val = int(v)
	return
}

func (c *RedisClient) ExistsCtx(ctx context.Context, key string) (val bool, err error) {
	v, err := c.rc.Exists(ctx, key).Result()
	if err != nil {
		return
	}
	val = v == 1
	return
}

func (c *RedisClient) HsetCtx(ctx context.Context, key, field, value string) error {
	return c.rc.HSet(ctx, key, field, value).Err()
}
