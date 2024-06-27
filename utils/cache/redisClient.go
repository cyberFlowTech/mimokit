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

func (c *RedisClient) HgetCtx(ctx context.Context, key, field string) (val string, err error) {
	val, err = c.rc.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (c *RedisClient) ZrevrangeWithScoresCtx(ctx context.Context, key string, start, stop int64) (
	val []Pair, err error) {

	v, err := c.rc.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return
	}

	val = toPairs(v)
	return
}

func (c *RedisClient) ZRevRangeByScore(ctx context.Context, key string, min, max string, offset, count int64) (
	[]string, error) {

	return c.rc.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}).Result()
}

func (c *RedisClient) ZcardCtx(ctx context.Context, key string) (val int, err error) {
	v, err := c.rc.ZCard(ctx, key).Result()
	if err != nil {
		return
	}
	val = int(v)
	return
}

func (c *RedisClient) ZaddsCtx(ctx context.Context, key string, ps ...Pair) (val int64, err error) {
	var zs []redis.Z
	for _, p := range ps {
		z := redis.Z{Score: float64(p.Score), Member: p.Key}
		zs = append(zs, z)
	}
	v, err := c.rc.ZAdd(ctx, key, zs...).Result()
	if err != nil {
		return
	}
	val = v
	return
}

func (c *RedisClient) SaddCtx(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return c.rc.SAdd(ctx, key, members...).Result()
}

func (c *RedisClient) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return c.rc.SRem(ctx, key, members...).Result()
}
func (c *RedisClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return c.rc.SIsMember(ctx, key, member).Result()
}
func (c *RedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.rc.SMembers(ctx, key).Result()
}
func (c *RedisClient) SCard(ctx context.Context, key string) (int64, error) {
	return c.rc.SCard(ctx, key).Result()
}

func (c *RedisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	return c.rc.Del(ctx, keys...).Result()
}

func (c *RedisClient) Zrem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return c.rc.ZRem(ctx, key, members).Result()
}

func (c *RedisClient) IncrCtx(ctx context.Context, key string) (int64, error) {
	return c.rc.Incr(ctx, key).Result()
}

func (c *RedisClient) ZRevRangeByScoreWithScoresAndLimitCtx(ctx context.Context, key string, start,
	stop int64, page, size int) (val []Pair, err error) {
	if size <= 0 {
		return
	}
	v, err := c.rc.ZRevRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
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

func (c *RedisClient) ZCount(ctx context.Context, key string, min, max string) (val int64, err error) {
	v, err := c.rc.ZCount(ctx, key, min, max).Result()
	if err != nil {
		return
	}

	return v, nil
}

func (c *RedisClient) Ttl(ctx context.Context, key string) (val int, err error) {
	duration, err := c.rc.TTL(ctx, key).Result()
	if err != nil {
		return
	}

	val = int(duration / time.Second)
	return
}

func (c *RedisClient) ZRemCtx(ctx context.Context, key string, member ...interface{}) (int64, error) {
	return c.rc.ZRem(ctx, key, member...).Result()
}

func (c *RedisClient) ZScoreCtx(ctx context.Context, key string, member string) (float64, error) {
	s, err := c.rc.ZScore(ctx, key, member).Result()
	if err == redis.Nil {
		return -1, nil
	}

	return s, err
}

func (c *RedisClient) SetNX(ctx context.Context, key string, val interface{}, timeout int) (bool, error) {
	return c.rc.SetNX(ctx, key, val, time.Duration(timeout)*time.Second).Result()
}

func (c *RedisClient) ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	return c.rc.ZRemRangeByScore(ctx, key, min, max).Result()
}

func (c *RedisClient) XAdd(ctx context.Context, Stream string, Id string, Values interface{}) (string, error) {
	a := &redis.XAddArgs{
		Stream:     Stream,
		NoMkStream: false,
		ID:         Id,
		Values:     Values,
	}
	return c.rc.XAdd(ctx, a).Result()
}

type XReadGroupItem struct {
	Group    string
	Consumer string
	Streams  []string // list of streams and ids, e.g. stream1 stream2 id1 id2
	Count    int64
	Block    time.Duration
	NoAck    bool
}

type XReadGroupResp struct {
	Streams  string
	Messages []XMessageItem
}

type XMessageItem struct {
	ID     string
	Values map[string]interface{}
}

func (c *RedisClient) XReadGroup(ctx context.Context, data *XReadGroupItem) ([]XReadGroupResp, error) {
	a := &redis.XReadGroupArgs{
		Group:    data.Group,
		Consumer: data.Consumer,
		Streams:  data.Streams,
		Count:    data.Count,
		Block:    data.Block,
		NoAck:    data.NoAck,
	}
	resp, err := c.rc.XReadGroup(ctx, a).Result()
	if err != nil {
		return nil, err
	}
	msgData := make([]XReadGroupResp, 0)
	for _, v := range resp {
		r := XReadGroupResp{
			Streams:  v.Stream,
			Messages: make([]XMessageItem, 0),
		}
		for _, m := range v.Messages {
			r.Messages = append(r.Messages, XMessageItem{
				ID:     m.ID,
				Values: m.Values,
			})
		}
	}
	return msgData, nil
}

func (c *RedisClient) XAck(ctx context.Context, stream, group string, ids ...string) (int64, error) {
	return c.rc.XAck(ctx, stream, group, ids...).Result()
}
