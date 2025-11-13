package util

import (
	"sync"
	"time"

	"github.com/vearne/autotest/internal/config"
	slog "github.com/vearne/simplelog"
)

// CacheItem 缓存项
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// IsExpired 检查是否过期
func (item *CacheItem) IsExpired() bool {
	return time.Now().After(item.ExpiresAt)
}

// Cache 通用缓存
type Cache struct {
	items   map[string]*CacheItem
	mutex   sync.RWMutex
	ttl     time.Duration
	maxSize int
	enabled bool

	// 统计信息
	hits   int64
	misses int64
}

// NewCache 创建新的缓存实例
func NewCache(cfg config.AutoTestConfig) *Cache {
	cache := &Cache{
		items:   make(map[string]*CacheItem),
		ttl:     cfg.Global.Cache.TTL,
		maxSize: cfg.Global.Cache.MaxSize,
		enabled: cfg.Global.Cache.Enabled,
	}

	// 设置默认值
	if cache.ttl <= 0 {
		cache.ttl = 5 * time.Minute // 默认5分钟
	}
	if cache.maxSize <= 0 {
		cache.maxSize = 100 // 默认100个条目，适合gRPC描述符缓存
	}

	// 缓存默认开启策略：
	// 考虑到gRPC描述符缓存收益明显(节省50-200秒)且开销很小，默认开启缓存
	// 只有当配置文件中有完整的缓存配置且显式设置 enabled: false 时才禁用

	cache.enabled = true // 默认开启

	// 检查用户是否有完整的缓存配置且显式禁用
	hasFullCacheConfig := (cfg.Global.Cache.TTL > 0 || cfg.Global.Cache.MaxSize > 0) && !cfg.Global.Cache.Enabled

	if hasFullCacheConfig {
		// 用户有缓存配置但明确设置了 enabled: false，尊重用户选择
		cache.enabled = false
		slog.Info("Cache disabled by user configuration (enabled: false)")
	} else if cfg.Global.Cache.Enabled {
		// 用户显式开启
		cache.enabled = true
		slog.Info("Cache enabled by user configuration (enabled: true)")
	} else {
		// 其他情况：没有配置或只有部分配置，默认开启
		cache.enabled = true
		slog.Info("Cache enabled by default for better performance")
	}

	// 启动清理协程
	if cache.enabled {
		go cache.cleanup()
		slog.Info("Cache initialized: TTL=%v, MaxSize=%d", cache.ttl, cache.maxSize)
	} else {
		slog.Info("Cache disabled by configuration")
	}

	return cache
}

// Set 设置缓存项
func (c *Cache) Set(key string, value interface{}) {
	if !c.enabled {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 检查是否需要清理空间
	if len(c.items) >= c.maxSize {
		c.evictOldest()
	}

	c.items[key] = &CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}

	slog.Debug("Cache SET: %s", key)
}

// Get 获取缓存项
func (c *Cache) Get(key string) (interface{}, bool) {
	if !c.enabled {
		return nil, false
	}

	c.mutex.RLock()
	item, exists := c.items[key]
	c.mutex.RUnlock()

	if !exists {
		c.mutex.Lock()
		c.misses++
		c.mutex.Unlock()
		slog.Debug("Cache MISS: %s", key)
		return nil, false
	}

	if item.IsExpired() {
		c.mutex.Lock()
		delete(c.items, key)
		c.misses++
		c.mutex.Unlock()
		slog.Debug("Cache EXPIRED: %s", key)
		return nil, false
	}

	c.mutex.Lock()
	c.hits++
	c.mutex.Unlock()
	slog.Debug("Cache HIT: %s", key)
	return item.Value, true
}

// Delete 删除缓存项
func (c *Cache) Delete(key string) {
	if !c.enabled {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
	slog.Debug("Cache DELETE: %s", key)
}

// Clear 清空缓存
func (c *Cache) Clear() {
	if !c.enabled {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items = make(map[string]*CacheItem)
	c.hits = 0
	c.misses = 0
	slog.Info("Cache cleared")
}

// Size 获取缓存大小
func (c *Cache) Size() int {
	if !c.enabled {
		return 0
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.items)
}

// Stats 获取缓存统计信息
func (c *Cache) Stats() (hits, misses int64, hitRate float64) {
	if !c.enabled {
		return 0, 0, 0
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	hits = c.hits
	misses = c.misses
	total := hits + misses

	if total > 0 {
		hitRate = float64(hits) / float64(total) * 100
	}

	return hits, misses, hitRate
}

// evictOldest 淘汰最旧的缓存项
func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, item := range c.items {
		if first || item.ExpiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.ExpiresAt
			first = false
		}
	}

	if oldestKey != "" {
		delete(c.items, oldestKey)
		slog.Debug("Cache EVICT: %s", oldestKey)
	}
}

// cleanup 定期清理过期项
func (c *Cache) cleanup() {
	ticker := time.NewTicker(c.ttl / 2) // 每半个TTL清理一次
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()

		var expiredKeys []string
		for key, item := range c.items {
			if item.IsExpired() {
				expiredKeys = append(expiredKeys, key)
			}
		}

		for _, key := range expiredKeys {
			delete(c.items, key)
		}

		if len(expiredKeys) > 0 {
			slog.Debug("Cache cleanup: removed %d expired items", len(expiredKeys))
		}

		c.mutex.Unlock()
	}
}

// CacheManager 缓存管理器
type CacheManager struct {
	// 只保留真正有价值的缓存
	GrpcDescriptorCache *Cache // gRPC描述符缓存
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(cfg config.AutoTestConfig) *CacheManager {
	return &CacheManager{
		GrpcDescriptorCache: NewCache(cfg),
	}
}

// GetStats 获取缓存的统计信息
func (cm *CacheManager) GetStats() map[string]map[string]interface{} {
	stats := make(map[string]map[string]interface{})

	hits, misses, hitRate := cm.GrpcDescriptorCache.Stats()
	stats["grpc_descriptor"] = map[string]interface{}{
		"hits":     hits,
		"misses":   misses,
		"hit_rate": hitRate,
		"size":     cm.GrpcDescriptorCache.Size(),
	}

	return stats
}

// ClearAll 清空缓存
func (cm *CacheManager) ClearAll() {
	cm.GrpcDescriptorCache.Clear()
	slog.Info("gRPC descriptor cache cleared")
}
