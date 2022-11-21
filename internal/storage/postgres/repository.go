package postgres

import (
	"context"

	"github.com/go-pg/pg/v10"
)

type Config struct {
	DB       string `envconfig:"APP_DB_NAME" default:"postgres"`
	Host     string `envconfig:"APP_DB_HOST" default:"localhost"`
	Port     string `envconfig:"APP_DB_PORT" default:"5432"`
	User     string `envconfig:"APP_DB_USER" default:"postgres"`
	Password string `envconfig:"APP_DB_PASS" default:"postgres"`
}

type DB struct {
	DB *pg.DB
}

func New(ctx context.Context, cfg Config) (DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		User:     cfg.User,
		Database: cfg.DB,
		Password: cfg.Password,
	})

	if err := db.Ping(ctx); err != nil {
		return DB{}, err
	}

	return DB{DB: db}, nil
}

type Currency struct {
	Symbol      string
	Description string
	Source      string
}

func (db DB) CountCurrencies(ctx context.Context) (int, error) {
	result, err := db.DB.ExecContext(ctx, "SELECT COUNT(*) FROM public.currencies")
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (db DB) CreateCurrency(ctx context.Context, symbol, description, source string) error {
	c := &Currency{Symbol: symbol, Description: description, Source: source}
	_, err := db.DB.Model(c).Context(ctx).Insert()
	return err
}
