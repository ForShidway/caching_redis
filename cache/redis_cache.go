package cache

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"redis_caching_demo/config"
	"redis_caching_demo/domain"
)

// ErrCacheMiss adalah error yang dikembalikan saat key tidak ada di Redis.
var ErrCacheMiss = fmt.Errorf("cache miss: key tidak ditemukan")

// cacheTracker menyimpan statistik HIT dan MISS secara thread-safe.
// Menggunakan atomic agar aman saat banyak goroutine mengakses sekaligus.
type cacheTracker struct {
	hits   int64 // jumlah cache HIT (atomic)
	misses int64 // jumlah cache MISS (atomic)
}

// RecordHit mencatat satu cache HIT secara atomic (thread-safe).
func (t *cacheTracker) RecordHit() {
	atomic.AddInt64(&t.hits, 1)
}

// RecordMiss mencatat satu cache MISS secara atomic (thread-safe).
func (t *cacheTracker) RecordMiss() {
	atomic.AddInt64(&t.misses, 1)
}

// GetStats mengembalikan snapshot statistik saat ini.
func (t *cacheTracker) GetStats() (hits, misses int64) {
	return atomic.LoadInt64(&t.hits), atomic.LoadInt64(&t.misses)
}

// RedisCache adalah implementasi konkrit dari domain.CacheProvider.
// Bertanggung jawab untuk semua operasi Redis: Get, Set, Delete, Ping.
type RedisCache struct {
	client  *redis.Client  // koneksi ke Redis/Memurai
	ttl     time.Duration  // Time-To-Live untuk setiap key
	tracker *cacheTracker  // tracker HIT/MISS — "track method"
}

// NewRedisCache membuat instance RedisCache baru dengan konfigurasi dari Config.
// Ini adalah satu-satunya tempat redis.Client dibuat (single responsibility).
func NewRedisCache(cfg *config.Config) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	return &RedisCache{
		client:  client,
		ttl:     cfg.CacheTTL,
		tracker: &cacheTracker{},
	}
}

// Get mengambil nilai dari Redis berdasarkan key.
// TRACK METHOD: otomatis mencatat HIT atau MISS ke tracker.
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()

	if err == redis.Nil {
		// Key tidak ditemukan di Redis → CACHE MISS
		r.tracker.RecordMiss()
		return "", ErrCacheMiss
	}

	if err != nil {
		// Error koneksi atau lainnya → tetap catat sebagai MISS
		r.tracker.RecordMiss()
		return "", fmt.Errorf("redis get error: %w", err)
	}

	// Key ditemukan → CACHE HIT
	r.tracker.RecordHit()
	return val, nil
}

// Set menyimpan nilai ke Redis dengan TTL yang dikonfigurasi.
// Key akan otomatis terhapus dari Redis setelah durasi TTL.
func (r *RedisCache) Set(ctx context.Context, key string, value string) error {
	err := r.client.Set(ctx, key, value, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}
	return nil
}

// Delete menghapus key dari Redis (cache invalidation).
// Dipakai saat data produk diperbarui agar cache tidak stale.
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis delete error: %w", err)
	}
	return nil
}

// Ping mengecek apakah koneksi ke Redis server aktif.
// Dipakai saat startup aplikasi untuk validasi koneksi.
func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// GetStats mengembalikan statistik cache HIT/MISS dalam format domain.CacheStats.
// Menghitung hit_rate secara otomatis berdasarkan counter.
func (r *RedisCache) GetStats() *domain.CacheStats {
	hits, misses := r.tracker.GetStats()
	total := hits + misses

	hitRate := "0.00%"
	if total > 0 {
		hitRate = fmt.Sprintf("%.2f%%", float64(hits)/float64(total)*100)
	}

	return &domain.CacheStats{
		Hits:         hits,
		Misses:       misses,
		TotalRequest: total,
		HitRate:      hitRate,
	}
}
