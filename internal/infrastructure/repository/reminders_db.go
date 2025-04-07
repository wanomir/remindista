package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	r "github.com/vedomirr/remindista/internal/entity/reminder"

	"github.com/jackc/pgx/v5"
)

func (db *PostgresDB) CreateReminder(ctx context.Context, rmd r.Reminder) (id int, err error) {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `INSERT INTO data.reminders (user_id, text, tag, prompt, frequency, next_reminder, is_deleted)
VALUES ($1, $2, $3, $4, $5, $6, FALSE)
RETURNING id;`

	if err = db.conn.QueryRow(ctx, query,
		rmd.UserId,
		rmd.Text,
		rmd.Tag,
		rmd.Prompt,
		rmd.Frequency,
		rmd.NextReminder,
	).Scan(&id); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			err = fmt.Errorf("failed to rollback: %w", err)
		}
		return 0, fmt.Errorf("failed to execute insert reminder query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

func (db *PostgresDB) GetReminder(ctx context.Context, id int) (rmd r.Reminder, err error) {
	query := `SELECT id, user_id, text, tag, prompt, frequency, next_reminder
FROM data.reminders
WHERE id = $1 AND is_deleted = FALSE;`

	if err = db.conn.QueryRow(ctx, query, id).Scan(
		&rmd.Id,
		&rmd.UserId,
		&rmd.Text,
		&rmd.Tag,
		&rmd.Prompt,
		&rmd.Frequency,
		&rmd.NextReminder,
	); errors.Is(err, pgx.ErrNoRows) {
		return rmd, nil
	} else if err != nil {
		return rmd, fmt.Errorf("failed to execute select reminder query: %w", err)
	}

	return rmd, nil
}

func (db *PostgresDB) GetRemindersByUserId(ctx context.Context, userId int) (rmds []r.Reminder, err error) {
	rmds = make([]r.Reminder, 0)

	query := `SELECT id, user_id, text, tag, prompt, frequency, next_reminder
FROM data.reminders
WHERE user_id = $1 AND is_deleted = FALSE;`

	rows, err := db.conn.Query(ctx, query, userId)
	if errors.Is(err, pgx.ErrNoRows) {
		return rmds, nil
	} else if err != nil {
		return rmds, fmt.Errorf("failed to execute select reminders query: %w", err)
	}

	for rows.Next() {
		var rmd r.Reminder

		if err := rows.Scan(
			&rmd.Id,
			&rmd.UserId,
			&rmd.Text,
			&rmd.Tag,
			&rmd.Prompt,
			&rmd.Frequency,
			&rmd.NextReminder,
		); err != nil {
			return rmds, fmt.Errorf("failed to scan row when quering reminders: %w", err)
		}

		rmds = append(rmds, rmd)
	}

	return rmds, nil
}

func (db *PostgresDB) GetRemindersByUserIdAndTime(ctx context.Context, userId int, userTime time.Time) (rmds []r.Reminder, err error) {
	rmds = make([]r.Reminder, 0)

	query := `SELECT id, user_id, text, tag, prompt, frequency, next_reminder
FROM data.reminders
WHERE user_id = $1 AND next_reminder < $2 AND is_deleted = FALSE;`

	rows, err := db.conn.Query(ctx, query, userId, userTime)
	if errors.Is(err, pgx.ErrNoRows) {
		return rmds, nil
	} else if err != nil {
		return rmds, fmt.Errorf("failed to execute select reminders query: %w", err)
	}

	for rows.Next() {
		var rmd r.Reminder

		if err := rows.Scan(
			&rmd.Id,
			&rmd.UserId,
			&rmd.Text,
			&rmd.Tag,
			&rmd.Prompt,
			&rmd.Frequency,
			&rmd.NextReminder,
		); err != nil {
			return rmds, fmt.Errorf("failed to scan row when quering reminders: %w", err)
		}

		rmds = append(rmds, rmd)
	}

	return rmds, nil
}

func (db *PostgresDB) UpdateReminder(ctx context.Context, rmd r.Reminder) (affected int, err error) {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `WITH rows AS (
	UPDATE data.reminders
	SET user_id = $2, text = $3, tag = $4, prompt = $5, frequency = $6, next_reminder = $7
	WHERE id = $1 AND is_deleted = FALSE
	RETURNING 1
) SELECT COUNT(*) FROM rows;`

	if err = db.conn.QueryRow(ctx, query,
		rmd.Id,
		rmd.UserId,
		rmd.Text,
		rmd.Tag,
		rmd.Prompt,
		rmd.Frequency,
		rmd.NextReminder,
	).Scan(&affected); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			err = fmt.Errorf("failed to rollback: %w", err)
		}
		return 0, fmt.Errorf("failed to execute update reminder query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return affected, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return affected, nil
}

func (db *PostgresDB) DeleteReminder(ctx context.Context, id int) (affected int, err error) {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `WITH rows AS (
	UPDATE data.reminders
	SET is_deleted = true
	WHERE id = $1 AND is_deleted = FALSE
	RETURNING 1
) SELECT COUNT(*) FROM rows;`

	if err = db.conn.QueryRow(ctx, query, id).Scan(&affected); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			err = fmt.Errorf("failed to rollback: %w", err)
		}
		return affected, fmt.Errorf("failed to execute update reminders query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return affected, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return affected, nil
}

func (db *PostgresDB) DeleteRemindersByTag(ctx context.Context, userId int, tag string) (affected int, err error) {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `WITH rows AS (
	UPDATE data.reminders
	SET is_deleted = true
	WHERE user_id = $1 AND tag = $2 AND is_deleted = FALSE
	RETURNING 1
) SELECT COUNT(*) FROM rows;`

	if err = db.conn.QueryRow(ctx, query, userId, tag).Scan(&affected); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			err = fmt.Errorf("failed to rollback: %w", err)
		}
		return 0, fmt.Errorf("failed to execute update reminders query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return affected, nil
}

func (db *PostgresDB) DeleteRemindersByUserId(ctx context.Context, userId int) (affected int, err error) {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `WITH rows AS (
	UPDATE data.reminders
	SET is_deleted = true
	WHERE user_id = $1 AND is_deleted = FALSE
	RETURNING 1
) SELECT COUNT(*) FROM rows;`

	if err = db.conn.QueryRow(ctx, query, userId).Scan(&affected); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			err = fmt.Errorf("failed to rollback: %w", err)
		}
		return 0, fmt.Errorf("failed to execute update reminders query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return affected, nil
}
