package data

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"cache-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

type cacheRepo struct {
	data     *Data
	log      *log.Helper
	memStore map[string]string // In-memory storage when DB is not available
	memMutex sync.RWMutex      // Mutex for in-memory storage
}

// NewCacheRepo creates a new cache repository.
func NewCacheRepo(data *Data, logger log.Logger) biz.CacheRepo {
	return &cacheRepo{
		data:     data,
		log:      log.NewHelper(logger),
		memStore: make(map[string]string),
	}
}

// GetFromCache gets data from Redis cache.
func (r *cacheRepo) GetFromCache(ctx context.Context, key string) (string, error) {
	if r.data.rdb == nil {
		return "", fmt.Errorf("redis client not initialized")
	}

	val, err := r.data.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found in cache")
	}
	if err != nil {
		r.log.WithContext(ctx).Errorf("Redis get error: %v", err)
		return "", err
	}

	return val, nil
}

// SetToCache sets data to Redis cache with TTL.
func (r *cacheRepo) SetToCache(ctx context.Context, key, value string) error {
	if r.data.rdb == nil {
		return fmt.Errorf("redis client not initialized")
	}

	// Set with 1 hour TTL
	err := r.data.rdb.Set(ctx, key, value, time.Hour).Err()
	if err != nil {
		r.log.WithContext(ctx).Errorf("Redis set error: %v", err)
		return err
	}

	return nil
}

// DeleteFromCache deletes data from Redis cache.
func (r *cacheRepo) DeleteFromCache(ctx context.Context, key string) error {
	if r.data.rdb == nil {
		return fmt.Errorf("redis client not initialized")
	}

	err := r.data.rdb.Del(ctx, key).Err()
	if err != nil {
		r.log.WithContext(ctx).Errorf("Redis delete error: %v", err)
		return err
	}

	return nil
}

// GetFromDB gets data from MySQL database or in-memory storage.
func (r *cacheRepo) GetFromDB(ctx context.Context, key string) (string, error) {
	if r.data.db == nil {
		// Use in-memory storage
		r.memMutex.RLock()
		value, exists := r.memStore[key]
		r.memMutex.RUnlock()

		if !exists {
			return "", fmt.Errorf("key not found in memory storage")
		}
		return value, nil
	}

	var value string
	query := "SELECT value FROM cache_data WHERE `key` = ?"
	err := r.data.db.QueryRowContext(ctx, query, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("key not found in database")
	}
	if err != nil {
		r.log.WithContext(ctx).Errorf("Database query error: %v", err)
		return "", err
	}

	return value, nil
}

// SetToDB sets data to MySQL database or in-memory storage.
func (r *cacheRepo) SetToDB(ctx context.Context, key, value string) error {
	if r.data.db == nil {
		// Use in-memory storage
		r.memMutex.Lock()
		r.memStore[key] = value
		r.memMutex.Unlock()
		return nil
	}

	// Use INSERT ... ON DUPLICATE KEY UPDATE for upsert
	query := `INSERT INTO cache_data (` + "`key`" + `, value, created_at, updated_at) 
			   VALUES (?, ?, NOW(), NOW()) 
			   ON DUPLICATE KEY UPDATE value = VALUES(value), updated_at = NOW()`

	_, err := r.data.db.ExecContext(ctx, query, key, value)
	if err != nil {
		r.log.WithContext(ctx).Errorf("Database insert/update error: %v", err)
		return err
	}

	return nil
}

// CheckRedisHealth checks if Redis is healthy by performing a ping operation.
func (r *cacheRepo) CheckRedisHealth(ctx context.Context) bool {
	if r.data.rdb == nil {
		r.log.WithContext(ctx).Warn("Redis client not initialized")
		return false
	}

	// Use a timeout context for health check
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.data.rdb.Ping(timeoutCtx).Result()
	if err != nil {
		r.log.WithContext(ctx).Errorf("Redis health check failed: %v", err)
		return false
	}

	return true
}

// CheckDBHealth checks if database is healthy by performing a simple query.
func (r *cacheRepo) CheckDBHealth(ctx context.Context) bool {
	if r.data.db == nil {
		r.log.WithContext(ctx).Warn("Database client not initialized")
		return false
	}

	// Use a timeout context for health check
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Perform a simple query to check database connectivity
	err := r.data.db.PingContext(timeoutCtx)
	if err != nil {
		r.log.WithContext(ctx).Errorf("Database health check failed: %v", err)
		return false
	}

	return true
}
