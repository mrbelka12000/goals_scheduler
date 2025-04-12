package goals

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type (
	Repository interface {
		CreateNotification(ctx context.Context, n Goal) (int64, error)
		Delete(ctx context.Context, id int64) error
		Get(ctx context.Context, id int64) (Goal, error)
		List(ctx context.Context, pars GoalPars) ([]Goal, int64, error)
		DeleteAllUsersGoals(ctx context.Context, chatID string) error
		Update(ctx context.Context, n Goal, id int64) error
	}

	repo struct {
		db *sql.DB
	}
)

func NewRepo(db *sql.DB) Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) CreateNotification(ctx context.Context, n Goal) (int64, error) {
	query := `
INSERT INTO notifications (
	telegram_chat_id,
	schedule_type,
	interval_seconds,
	day,                    
	scheduled_time,
	message,
	status,
	next_execution,
	created_at,
	updated_at
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING id`

	var id int64
	err := r.db.QueryRowContext(ctx, query,
		n.TelegramChatID,
		n.ScheduleType,
		n.IntervalSeconds,
		n.Day,
		n.ScheduledTime,
		n.Message,
		n.Status,
		n.NextExecution,
		n.CreatedAt,
		n.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create notification: %w", err)
	}

	return id, nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM notifications WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repo) Get(ctx context.Context, id int64) (Goal, error) {
	query := `
SELECT id,
    telegram_chat_id,
	schedule_type,
	interval_seconds,
	day,
	scheduled_time,
	message,
	status,
	next_execution,
	created_at,
	updated_at
       FROM notifications WHERE id = $1`

	var goal Goal
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(
			&goal.ID,
			&goal.TelegramChatID,
			&goal.ScheduleType,
			&goal.IntervalSeconds,
			&goal.Day,
			&goal.ScheduledTime,
			&goal.Message,
			&goal.Status,
			&goal.NextExecution,
			&goal.CreatedAt,
			&goal.UpdatedAt,
		)
	if err != nil {
		return Goal{}, err
	}

	return goal, nil
}

func (r *repo) List(ctx context.Context, pars GoalPars) ([]Goal, int64, error) {
	query := `
SELECT id,
    telegram_chat_id,
	schedule_type,
	interval_seconds,
	day,
	scheduled_time,
	message,
	status,
	next_execution,
	created_at,
	updated_at
    FROM notifications WHERE`

	var args []interface{}

	if pars.ID != nil {
		args = append(args, *pars.ID)
		query += fmt.Sprintf(" id = $%v AND", len(args))
	}

	if pars.ChatID != nil {
		args = append(args, *pars.ChatID)
		query += fmt.Sprintf(" telegram_chat_id = $%v AND", len(args))
	}

	if pars.ScheduleType != nil {
		args = append(args, *pars.ScheduleType)
		query += fmt.Sprintf(" schedule_type = $%v AND", len(args))
	}

	if pars.Status != nil {
		args = append(args, *pars.Status)
		query += fmt.Sprintf(" status = $%v AND", len(args))
	}

	query = query[:len(query)-4] // Remove the trailing " AND"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []Goal
	for rows.Next() {
		var goal Goal
		err := rows.Scan(
			&goal.ID,
			&goal.TelegramChatID,
			&goal.ScheduleType,
			&goal.IntervalSeconds,
			&goal.Day,
			&goal.ScheduledTime,
			&goal.Message,
			&goal.Status,
			&goal.NextExecution,
			&goal.CreatedAt,
			&goal.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("scan goal: %w", err)
		}
		notifications = append(notifications, goal)
	}

	return notifications, 0, nil
}

func (r *repo) DeleteAllUsersGoals(ctx context.Context, chatID string) error {
	query := "DELETE FROM notifications where telegram_chat_id = $1"
	_, err := r.db.ExecContext(ctx, query, chatID)
	return err
}

func (r *repo) Update(ctx context.Context, n Goal, id int64) error {
	updateValues := []interface{}{id}
	queryUpdate := `UPDATE notifications`
	querySet := ` SET id = $1`
	queryWhere := ` WHERE id = $1`

	if n.TelegramChatID != 0 {
		updateValues = append(updateValues, n.TelegramChatID)
		querySet += `, telegram_chat_id = $` + strconv.Itoa(len(updateValues))
	}
	if n.ScheduleType != "" {
		updateValues = append(updateValues, n.ScheduleType)
		querySet += `, schedule_type = $` + strconv.Itoa(len(updateValues))
	}
	if n.IntervalSeconds != nil {
		updateValues = append(updateValues, *n.IntervalSeconds)
		querySet += `, interval_seconds = $` + strconv.Itoa(len(updateValues))
	}

	if n.ScheduledTime != nil {
		updateValues = append(updateValues, *n.ScheduledTime)
		querySet += `, scheduled_time = $` + strconv.Itoa(len(updateValues))
	}
	if n.Message != "" {
		updateValues = append(updateValues, n.Message)
		querySet += `, message = $` + strconv.Itoa(len(updateValues))
	}
	if n.Status != "" {
		updateValues = append(updateValues, n.Status)
		querySet += `, status = $` + strconv.Itoa(len(updateValues))
	}
	if n.NextExecution != nil {
		updateValues = append(updateValues, *n.NextExecution)
		querySet += `, next_execution = $` + strconv.Itoa(len(updateValues))
	}

	updateValues = append(updateValues, time.Now())
	querySet += `, updated_at = $` + strconv.Itoa(len(updateValues))

	_, err := r.db.ExecContext(ctx, queryUpdate+querySet+queryWhere, updateValues...)
	if err != nil {
		return fmt.Errorf("error updating notification: %w", err)
	}

	return nil
}
