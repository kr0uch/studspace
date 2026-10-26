package postgres

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresConfig struct {
	Host     string `yaml:"postgres_host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `yaml:"postgres_port" env:"POSTGRES_PORT" env-default:"5432"`
	DBName   string `yaml:"postgres_db" env:"POSTGRES_DB" env-default:"studspace"`
	User     string `yaml:"postgres_user" env:"POSTGRES_USER" env-default:"postgres"`
	Password string `yaml:"postgres_password" env:"POSTGRES_PASSWORD" env-required:"true"`
	SSLMode  string `yaml:"postgres_sslmode" env:"POSTGRES_SSLMODE" env-default:"disable"`

	MaxOpenConns    int           `yaml:"postgres_max_open_conns" env:"POSTGRES_MAX_OPEN_CONNS" env-default:"25"`
	MaxIdleConns    int           `yaml:"postgres_max_idle_conns" env:"POSTGRES_MAX_IDLE_CONNS" env-default:"10"`
	ConnMaxLifetime time.Duration `yaml:"postgres_conn_max_lifetime" env:"POSTGRES_CONN_MAX_LIFETIME" env-default:"1h"`
	ConnMaxIdleTime time.Duration `yaml:"postgres_conn_max_idle_time" env:"POSTGRES_CONN_MAX_IDLE_TIME" env-default:"15m"`
	ConnectTimeout  time.Duration `yaml:"postgres_connect_timeout" env:"POSTGRES_CONNECT_TIMEOUT" env-default:"5s"`
}

func (c PostgresConfig) DSN() string {
	query := url.Values{}
	query.Set("sslmode", c.SSLMode)

	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, c.Port),
		Path:     c.DBName,
		RawQuery: query.Encode(),
	}
	return dsn.String()
}

type PostgresDB struct {
	*sqlx.DB
}

func NewPostgresDB(ctx context.Context, cfg PostgresConfig) (*PostgresDB, error) {
	pgxCfg, err := pgx.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}
	pgxCfg.ConnectTimeout = cfg.ConnectTimeout

	db := sqlx.NewDb(stdlib.OpenDB(*pgxCfg), "pgx")

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	pingCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	if err = db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres %s:%s/%s: %w", cfg.Host, cfg.Port, cfg.DBName, err)
	}

	return &PostgresDB{db}, nil
}
