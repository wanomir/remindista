package repository

import (
	"github.com/vedomirr/l"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresDB struct {
	conn *pgxpool.Pool
	log  *zap.Logger
}

func NewPostgresDB(pool *pgxpool.Pool) *PostgresDB {
	return &PostgresDB{
		conn: pool,
		log:  l.Logger(),
	}
}
