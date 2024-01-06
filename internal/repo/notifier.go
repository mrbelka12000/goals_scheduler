package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"goals_scheduler/internal/models"
)

type notifier struct {
	db *sql.DB
}

func newNotifier(db *sql.DB) *notifier {
	return &notifier{
		db: db,
	}
}

func (r *notifier) Create(ctx context.Context, obj *models.NotifierCU) (int64, error) {
	query := "INSERT INTO notifier (usr_id, status_id, ticker, last_updated, expires) VALUES (?, ?, ?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, *obj.UsrID, obj.Status, *obj.Ticker, time.Now(), obj.Expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %v", err)
	}

	return id, nil
}

func (r *notifier) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM notifier WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *notifier) Get(ctx context.Context, id int64) (models.Notifier, error) {
	query := "SELECT id, usr_id, status_id, ticker, last_updated, expires FROM notifier WHERE id = ?"
	var notifier models.Notifier
	err := r.db.QueryRowContext(ctx, query, id).Scan(&notifier.ID, &notifier.UsrID, &notifier.Status, &notifier.Ticker, &notifier.LastUpdated, &notifier.Expires)
	if err != nil {
		return models.Notifier{}, err
	}

	return notifier, nil
}

func (r *notifier) List(ctx context.Context, pars models.NotifierPars) ([]models.Notifier, int64, error) {
	query := "SELECT id, usr_id, notifier_id, message, status_id, deadline FROM goals WHERE"
	args := []interface{}{}

	if pars.ID != nil {
		query += " id = ? AND"
		args = append(args, *pars.ID)
	}

	if pars.UsrID != nil {
		query += " usr_id = ? AND"
		args = append(args, *pars.UsrID)
	}

	if pars.Status != nil {
		query += " status_id = ? AND"
		args = append(args, *pars.Status)
	}

	query = query[:len(query)-4] // Remove the trailing " AND"
	fmt.Println(query, args)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifiers []models.Notifier
	for rows.Next() {
		var notifier models.Notifier
		err := rows.Scan(&notifier.ID, &notifier.UsrID, &notifier.Status, &notifier.Ticker, &notifier.LastUpdated, &notifier.Expires)
		if err != nil {
			return nil, 0, err
		}
		notifiers = append(notifiers, notifier)
	}

	countQuery := "SELECT COUNT(*) FROM notifier WHERE"
	countQuery += query[6:] // Remove the "SELECT" part from the count query
	var count int64
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return notifiers, count, nil
}
