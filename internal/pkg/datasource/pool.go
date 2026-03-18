package datasource

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PoolConfig struct {
	dsn string
}

func NewPoolConfig(dsn string) *PoolConfig {
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	return &PoolConfig{
		dsn: dsn,
	}
}

type DbPool struct {
	config *PoolConfig
	ctx    context.Context
	cfg    *pgxpool.Config
	pool   *pgxpool.Pool
}

func NewDbPool(config *PoolConfig) *DbPool {
	ctx := context.Background()

	cfg, err := pgxpool.ParseConfig(config.dsn)
	if err != nil {
		log.Fatalf("Не удалось распарсить DSN: %v", err)
		return nil
	}
	cfg.MaxConns = 100
	return &DbPool{
		ctx:    ctx,
		cfg:    cfg,
		config: config,
	}
}

func (p *DbPool) GetPool() *pgxpool.Pool {
	return p.pool
}

func (p *DbPool) GetContext() context.Context {
	return p.ctx
}

func (p *DbPool) Connect() error {
	pool, err := pgxpool.ConnectConfig(p.ctx, p.cfg)
	if err != nil {
		log.Fatalf("Не удалось установить соединение: %v", err)
		return err
	}
	p.pool = pool
	return nil
}

func (p *DbPool) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

func (p *DbPool) Query(sql string, args ...interface{}) (pgx.Rows, error) {
	if p.pool == nil {
		return nil, errors.New("Отсутствует подключение к базе")
	}
	rows, err := p.pool.Query(p.ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (p *DbPool) QueryRow(sql string, args ...interface{}) (pgx.Row, error) {
	if p.pool == nil {
		return nil, errors.New("Отсутствует подключение к базе")
	}
	row := p.pool.QueryRow(p.ctx, sql, args...)
	return row, nil
}

func (p *DbPool) Exec(sql string, args ...interface{}) error {
	if p.pool == nil {
		return errors.New("Отсутствует подключение к базе")
	}
	_, err := p.pool.Exec(p.ctx, sql, args...)
	return err
}
