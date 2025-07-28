package data

import (
	"database/sql"
	"strings"
	"time"

	"cache-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewCacheRepo)

// Data .
type Data struct {
	rdb redis.Cmdable
	db  *sql.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(logger)

	// Initialize Redis client
	var rdb redis.Cmdable

	// 获取Redis代理地址
	addr := strings.TrimSpace(c.Redis.Addr)
	// 移除集群总线端口（如果有的话）
	if strings.Contains(addr, "@") {
		addr = strings.Split(addr, "@")[0]
	}
	helper.Info("使用Redis代理模式连接: " + addr)
	rdb = redis.NewClient(&redis.Options{
		Addr:         addr,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		DialTimeout:  10 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		// 增加重试配置
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	})

	// Initialize MySQL database
	db, err := sql.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		helper.Errorf("failed to open database: %v", err)
		return nil, nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test database connection
	if err := db.Ping(); err != nil {
		helper.Warnf("failed to ping database: %v, will use in-memory storage", err)
		db.Close()
		db = nil
	}

	// Create table if not exists (only if database is available)
	if db != nil {
		createTableSQL := `CREATE TABLE IF NOT EXISTS cache_data (
			id INT AUTO_INCREMENT PRIMARY KEY,
			` + "`key`" + ` VARCHAR(255) NOT NULL UNIQUE,
			value TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`

		if _, err := db.Exec(createTableSQL); err != nil {
			helper.Errorf("failed to create table: %v", err)
			return nil, nil, err
		}
	}

	cleanup := func() {
		helper.Info("closing the data resources")
		// Close Redis connection based on type
		switch client := rdb.(type) {
		case *redis.Client:
			if err := client.Close(); err != nil {
				helper.Errorf("failed to close redis client: %v", err)
			}
		case *redis.ClusterClient:
			if err := client.Close(); err != nil {
				helper.Errorf("failed to close redis cluster: %v", err)
			}
		}
		if err := db.Close(); err != nil {
			helper.Errorf("failed to close database: %v", err)
		}
	}

	return &Data{
		rdb: rdb,
		db:  db,
	}, cleanup, nil
}
