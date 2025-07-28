package biz

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrDataNotFound is data not found.
	ErrDataNotFound = errors.NotFound("DATA_NOT_FOUND", "data not found")
)

// CacheData represents a cache data model.
type CacheData struct {
	Key   string
	Value string
}

// CacheRepo is a cache repository interface.
type CacheRepo interface {
	// Cache operations
	GetFromCache(ctx context.Context, key string) (string, error)
	SetToCache(ctx context.Context, key, value string) error
	DeleteFromCache(ctx context.Context, key string) error

	// Database operations
	GetFromDB(ctx context.Context, key string) (string, error)
	SetToDB(ctx context.Context, key, value string) error

	// Health check operations
	CheckRedisHealth(ctx context.Context) bool
	CheckDBHealth(ctx context.Context) bool
}

// CacheUsecase is a cache usecase.
type CacheUsecase struct {
	repo CacheRepo
	log  *log.Helper
}

// NewCacheUsecase creates a new cache usecase.
func NewCacheUsecase(repo CacheRepo, logger log.Logger) *CacheUsecase {
	return &CacheUsecase{repo: repo, log: log.NewHelper(logger)}
}

// GetData gets data from cache first, if not found, get from database.
func (uc *CacheUsecase) GetData(ctx context.Context, key string) (*CacheData, string, error) {
	uc.log.WithContext(ctx).Infof("GetData: key=%s", key)

	// Try to get from cache first
	value, err := uc.repo.GetFromCache(ctx, key)
	if err == nil {
		uc.log.WithContext(ctx).Infof("Data found in cache: key=%s", key)
		return &CacheData{Key: key, Value: value}, "cache", nil
	}

	// If not found in cache, get from database
	uc.log.WithContext(ctx).Infof("Data not found in cache, querying database: key=%s", key)
	value, err = uc.repo.GetFromDB(ctx, key)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Data not found in database: key=%s, error=%v", key, err)
		return nil, "", ErrDataNotFound
	}

	// Set to cache for next time
	if cacheErr := uc.repo.SetToCache(ctx, key, value); cacheErr != nil {
		uc.log.WithContext(ctx).Warnf("Failed to set cache: key=%s, error=%v", key, cacheErr)
	}

	uc.log.WithContext(ctx).Infof("Data found in database: key=%s", key)
	return &CacheData{Key: key, Value: value}, "database", nil
}

// SetData sets data to database and deletes cache.
func (uc *CacheUsecase) SetData(ctx context.Context, key, value string) error {
	uc.log.WithContext(ctx).Infof("SetData: key=%s, value=%s", key, value)

	// Set to database
	if err := uc.repo.SetToDB(ctx, key, value); err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to set data to database: key=%s, error=%v", key, err)
		return fmt.Errorf("failed to set data to database: %w", err)
	}

	// Delete from cache to ensure consistency
	if err := uc.repo.DeleteFromCache(ctx, key); err != nil {
		uc.log.WithContext(ctx).Warnf("Failed to delete cache: key=%s, error=%v", key, err)
		// Don't return error here, as the main operation (database write) succeeded
	}

	uc.log.WithContext(ctx).Infof("Data set successfully: key=%s", key)
	return nil
}

// CheckRedisHealth checks if Redis is healthy.
func (uc *CacheUsecase) CheckRedisHealth(ctx context.Context) bool {
	return uc.repo.CheckRedisHealth(ctx)
}

// CheckDBHealth checks if database is healthy.
func (uc *CacheUsecase) CheckDBHealth(ctx context.Context) bool {
	return uc.repo.CheckDBHealth(ctx)
}
