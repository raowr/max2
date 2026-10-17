/*
* @desc:缓存处理
 */

package service

import (
	"github.com/gogf/gf/v2/database/gredis" // ← 加这行
	"github.com/tiger1103/gfast-cache/cache"
)

type ICache interface {
	cache.IGCache
	Redis() *gredis.Redis // ← 加这个方法
}

var c ICache

func Cache() ICache {
	if c == nil {
		panic("implement not found for interface ICache, forgot register?")
	}
	return c
}

func RegisterCache(che ICache) {
	c = che
}
