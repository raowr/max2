/*
* @desc:缓存处理

 */

package cache

import (
	"gate-service/internal/consts"
	"gate-service/internal/service"

	"github.com/tiger1103/gfast-cache/cache"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func init() {
	service.RegisterCache(New())
}

func New() *sCache {
	var (
		ctx            = gctx.New()
		cacheContainer *cache.GfCache
	)
	prefix := g.Cfg().MustGet(ctx, "system.cache.prefix").String()
	model := g.Cfg().MustGet(ctx, "system.cache.model").String()
	if model == consts.CacheModelRedis {
		// redis
		cacheContainer = cache.NewRedis(prefix)
	} else {
		// memory
		cacheContainer = cache.New(prefix)
	}
	return &sCache{
		GfCache: cacheContainer,
		prefix:  prefix,
	}
}

type sCache struct {
	*cache.GfCache
	prefix string
}
