package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/mrbelka12000/goals_scheduler/internal/models"
)

type goal struct {
	db *sql.DB
}

func newGoal(db *sql.DB) *goal {
	return &goal{
		db: db,
	}
}

func (r *goal) Create(ctx context.Context, obj *models.GoalCU) (int64, error) {
	query := "INSERT INTO goals (usr_id, chat_id, message, status_id, deadline, timer, timer_enabled, last_updated) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"

	var id int64
	err := r.db.QueryRowContext(ctx, query, *obj.UsrID, *obj.ChatID, *obj.Text, *obj.Status, *obj.Deadline, *obj.Timer, obj.TimerEnabled, *obj.LastUpdated).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create goal: %w", err)
	}

	return id, nil
}

func (r *goal) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM goals WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *goal) Get(ctx context.Context, id int64) (models.Goal, error) {
	query := "SELECT id, usr_id, message, status_id, deadline FROM goals WHERE id = $1"
	var goal models.Goal
	err := r.db.QueryRowContext(ctx, query, id).Scan(&goal.ID, &goal.UsrID, &goal.Text, &goal.Status, &goal.Deadline)
	if err != nil {
		return models.Goal{}, err
	}

	return goal, nil
}

func (r *goal) List(ctx context.Context, pars models.GoalPars) ([]models.Goal, int64, error) {
	query := "SELECT id, usr_id, chat_id, message, status_id, deadline, timer, timer_enabled, last_updated FROM goals WHERE"

	var args []interface{}

	if pars.ID != nil {
		args = append(args, *pars.ID)
		query += fmt.Sprintf(" id = $%v AND", len(args))
	}

	if pars.UsrID != nil {
		args = append(args, *pars.UsrID)
		query += fmt.Sprintf(" usr_id = $%v AND", len(args))
	}

	if pars.StatusID != nil {
		args = append(args, *pars.StatusID)
		query += fmt.Sprintf(" status_id = $%v AND", len(args))
	}

	if pars.TimerEnabled != nil {
		args = append(args, *pars.TimerEnabled)
		query += fmt.Sprintf(" timer_enabled = $%v AND", len(args))
	}

	query = query[:len(query)-4] // Remove the trailing " AND"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list goals: %w", err)
	}
	defer rows.Close()

	var goals []models.Goal
	for rows.Next() {
		var goal models.Goal
		err := rows.Scan(
			&goal.ID,
			&goal.UsrID,
			&goal.ChatID,
			&goal.Text,
			&goal.Status,
			&goal.Deadline,
			&goal.Timer,
			&goal.TimerEnabled,
			&goal.LastUpdated,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("scan goal: %w", err)
		}
		goals = append(goals, goal)
	}

	return goals, 0, nil
}

func (r *goal) DeleteAllUsersGoals(ctx context.Context, usrID int) error {
	query := "DELETE FROM goals where usr_id = $1"
	_, err := r.db.ExecContext(ctx, query, usrID)
	return err
}

func (r *goal) Update(ctx context.Context, obj models.GoalCU, id int64) error {
	updateValues := []interface{}{id}
	queryUpdate := ` UPDATE goals`
	querySet := ` SET id = $1`
	queryWhere := ` WHERE id = $1`

	if obj.Timer != nil {
		updateValues = append(updateValues, time.Now().Add(*obj.Timer))
		querySet += ", last_updated = $" + strconv.Itoa(len(updateValues))
	}
	if obj.UsrID != nil {
		updateValues = append(updateValues, *obj.UsrID)
		querySet += ` , usr_id = $` + strconv.Itoa(len(updateValues))
	}
	if obj.ChatID != nil {
		updateValues = append(updateValues, *obj.ChatID)
		querySet += ` , chat_id = $` + strconv.Itoa(len(updateValues))
	}
	if obj.Status != nil {
		updateValues = append(updateValues, *obj.Status)
		querySet += ` , status_id = $` + strconv.Itoa(len(updateValues))
	}

	_, err := r.db.ExecContext(ctx, queryUpdate+querySet+queryWhere, updateValues...)
	if err != nil {
		return fmt.Errorf("error updating goals: %v", err)
	}

	return nil
}
