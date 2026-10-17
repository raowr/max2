package libBloom

import (
	"context"
	"hash/fnv"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
)

const (
	bloomFilterKey = "bloom:user_exists"
	bloomBitSize   = 2000000
	bloomHashCount = 14
)

func bloomHashes(key string) []int {
	positions := make([]int, bloomHashCount)
	for i := 0; i < bloomHashCount; i++ {
		h := fnv.New64a()
		h.Write([]byte(key + ":" + strconv.Itoa(i)))
		positions[i] = int(h.Sum64() % uint64(bloomBitSize))
	}
	return positions
}

func BloomAdd(ctx context.Context, username string) {
	redis := g.Redis()
	for _, pos := range bloomHashes(username) {
		redis.Do(ctx, "SETBIT", bloomFilterKey, pos, 1)
	}
}

func BloomMightExist(ctx context.Context, username string) bool {
	redis := g.Redis()
	for _, pos := range bloomHashes(username) {
		ret, err := redis.Do(ctx, "GETBIT", bloomFilterKey, pos)
		if err != nil {
			return true
		}
		if ret.Int() == 0 {
			return false
		}
	}
	return true
}

func BloomWarmUp(ctx context.Context) error {
	users, err := g.DB().Model("users").Fields("name").All()
	if err != nil {
		return err
	}
	for _, u := range users {
		BloomAdd(ctx, u["name"].String())
	}
	return nil
}
