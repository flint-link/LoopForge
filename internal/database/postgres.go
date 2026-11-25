package database

import (
	"context"
	"fmt"
	"time"

	"LoopForge/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Postgres 数据库连接池封装
type Postgres struct {
	Pool   *pgxpool.Pool
	logger *zap.Logger
}

// NewPostgres 创建数据库连接池
func NewPostgres(cfg *config.DatabaseConfig, logger *zap.Logger) (*Postgres, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	// 连接池配置
	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	// 测试连接
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	logger.Info("Database connected successfully",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.DBName),
	)

	return &Postgres{
		Pool:   pool,
		logger: logger,
	}, nil
}

// Close 关闭数据库连接池
func (p *Postgres) Close() {
	p.Pool.Close()
	p.logger.Info("Database connection pool closed")
}

// HealthCheck 健康检查
func (p *Postgres) HealthCheck(ctx context.Context) error {
	return p.Pool.Ping(ctx)
}

// Stats 返回连接池统计信息
func (p *Postgres) Stats() *pgxpool.Stat {
	return p.Pool.Stat()
}
