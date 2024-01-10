package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"goals_scheduler/internal/models"
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
	query := "INSERT INTO goals (usr_id, chat_id, message, status_id, deadline) VALUES (?, ?, ?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, *obj.UsrID, *obj.ChatID, *obj.Text, obj.Status, *obj.Deadline, time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %v", err)
	}

	return id, nil
}

func (r *goal) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM goals WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *goal) Get(ctx context.Context, id int64) (models.Goal, error) {
	query := "SELECT id, usr_id, message, status_id, deadline FROM goals WHERE id = ?"
	var goal models.Goal
	err := r.db.QueryRowContext(ctx, query, id).Scan(&goal.ID, &goal.UsrID, &goal.Text, &goal.Status, &goal.Deadline)
	if err != nil {
		return models.Goal{}, err
	}

	return goal, nil
}

func (r *goal) List(ctx context.Context, pars models.GoalPars) ([]models.Goal, int64, error) {
	query := "SELECT id, usr_id, chat_id, message, status_id, deadline FROM goals WHERE"

	var args []interface{}

	if pars.ID != nil {
		query += " id = ? AND"
		args = append(args, *pars.ID)
	}

	if pars.UsrID != nil {
		query += " usr_id = ? AND"
		args = append(args, *pars.UsrID)
	}

	if pars.StatusID != nil {
		query += " status_id = ? AND"
		args = append(args, *pars.StatusID)
	}

	query = query[:len(query)-4] // Remove the trailing " AND"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var goals []models.Goal
	for rows.Next() {
		var goal models.Goal
		err := rows.Scan(&goal.ID, &goal.UsrID, &goal.ChatID, &goal.Text, &goal.Status, &goal.Deadline)
		if err != nil {
			return nil, 0, err
		}
		goals = append(goals, goal)
	}

	return goals, 0, nil
}
