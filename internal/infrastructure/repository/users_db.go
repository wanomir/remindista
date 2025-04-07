package repository

import (
	"context"
	"errors"
	"fmt"

	u "github.com/vedomirr/remindista/internal/entity/user"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (db *PostgresDB) CreateUser(ctx context.Context, user u.User) (id int, err error) {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `INSERT INTO data.users (telegram_id, chat_id, is_running, location, window_floor, window_ceil, is_deleted)
VALUES ($1, $2, $3, $4, $5, $6, FALSE)
RETURNING id;`

	if err = db.conn.QueryRow(ctx, query,
		user.TelegramId,
		user.ChatId,
		user.IsRunning,
		user.Location.String(),
		user.WindowFloor,
		user.WindowCeil,
	).Scan(&id); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			err = fmt.Errorf("failed to rollback: %w", err)
		}
		return 0, fmt.Errorf("failed to execute insert user query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

func (db *PostgresDB) GetUser(ctx context.Context, id int) (user u.User, err error) {
	query := `SELECT id, telegram_id, chat_id, is_running, location, window_floor, window_ceil
FROM data.users
WHERE id = $1 AND is_deleted = FALSE;`

	var locationName string
	if err = db.conn.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.TelegramId,
		&user.ChatId,
		&user.IsRunning,
		&locationName,
		&user.WindowFloor,
		&user.WindowCeil,
	); errors.Is(err, pgx.ErrNoRows) {
		return user, nil
	} else if err != nil {
		return user, fmt.Errorf("failed to execute select user query: %w", err)
	}

	if err := user.SetLocation(locationName); err != nil {
		db.log.Error("failed to load user location", zap.Int("user id", user.Id), zap.Error(err))
	}

	return user, nil
}

func (db *PostgresDB) GetUserByTelegramId(ctx context.Context, telegramId int64) (user u.User, err error) {
	query := `SELECT id, telegram_id, chat_id, is_running, location, window_floor, window_ceil
FROM data.users
WHERE telegram_id = $1 AND is_deleted = FALSE;`

	var locationName string
	if err = db.conn.QueryRow(ctx, query, telegramId).Scan(
		&user.Id,
		&user.TelegramId,
		&user.ChatId,
		&user.IsRunning,
		&locationName,
		&user.WindowFloor,
		&user.WindowCeil,
	); errors.Is(err, pgx.ErrNoRows) {
		return user, nil
	} else if err != nil {
		return user, fmt.Errorf("failed to execute select user query: %w", err)
	}

	if err := user.SetLocation(locationName); err != nil {
		db.log.Error("failed to load user location", zap.Int("user id", user.Id), zap.Error(err))
	}

	return user, nil
}

func (db PostgresDB) GetAllUsers(ctx context.Context, limit int, offset int) (users []u.User, err error) {
	query := `SELECT id, telegram_id, chat_id, location, window_floor, window_ceil
FROM data.users
WHERE is_deleted = FALSE
LIMIT $1
OFFSET $2;`

	rows, err := db.conn.Query(ctx, query, limit, offset)
	if errors.Is(err, pgx.ErrNoRows) {
		return users, nil
	} else if err != nil {
		return users, fmt.Errorf("failed to execute select all users query: %w", err)
	}

	for rows.Next() {
		var user u.User
		var locationName string

		if err := rows.Scan(
			&user.Id,
			&user.TelegramId,
			&user.ChatId,
			&locationName,
			&user.WindowFloor,
			&user.WindowCeil,
		); err != nil {
			return users, fmt.Errorf("failed to scan row when quering users: %w", err)
		}

		if err := user.SetLocation(locationName); err != nil {
			db.log.Error("failed to load user location", zap.Int("user id", user.Id), zap.Error(err))
		}

		users = append(users, user)
	}

	return users, nil
}

func (db *PostgresDB) UpdateUser(ctx context.Context, user u.User) (affected int, err error) {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `WITH rows AS (
	UPDATE data.users
	SET telegram_id = $2, chat_id = $3, is_running = $4, location = $5, window_floor = $6, window_ceil = $7
	WHERE id = $1 AND is_deleted = FALSE
	RETURNING 1
) SELECT COUNT(*) FROM rows;`

	if err = db.conn.QueryRow(ctx, query,
		user.Id,
		user.TelegramId,
		user.ChatId,
		user.IsRunning,
		user.Location.String(),
		user.WindowFloor,
		user.WindowCeil,
	).Scan(&affected); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			err = fmt.Errorf("failed to rollback: %w", err)
		}
		return affected, fmt.Errorf("failed to execute update user query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return affected, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return affected, nil
}

func (db *PostgresDB) DeleteUser(ctx context.Context, id int, telegramId int64) (affected int, err error) {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `WITH rows AS (
	UPDATE data.users
	SET is_deleted = TRUE, is_running = FALSE
	WHERE (id = $1 OR telegram_id = $2) AND is_deleted = FALSE
	RETURNING 1
) SELECT COUNT(*) FROM rows;`

	if err = db.conn.QueryRow(ctx, query, id, telegramId).Scan(&affected); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			err = fmt.Errorf("failed to rollback: %w", err)
		}
		return affected, fmt.Errorf("failed to execute delete user query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return affected, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return affected, nil
}
