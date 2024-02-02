package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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

func (n *notifier) Create(ctx context.Context, obj *models.NotifierCU) (int64, error) {
	query := "INSERT INTO notifier (usr_id, chat_id, status_id, goal_id, notify,last_updated) VALUES (?, ?, ?, ?, ?, ?)"

	result, err := n.db.ExecContext(ctx, query,
		*obj.UsrID,
		*obj.ChatID,
		*obj.Status,
		*obj.GoalID,
		*obj.Notify,
		*obj.LastUpdated,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %v", err)
	}

	return id, nil
}

func (n *notifier) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM notifier WHERE id = ?"
	_, err := n.db.ExecContext(ctx, query, id)
	return err
}

func (n *notifier) Get(ctx context.Context, id int64) (models.Notifier, error) {
	query := "SELECT id, usr_id, chat_id, status_id, goal_id, notify, last_updated FROM notifier WHERE id = ?"
	var notif models.Notifier
	err := n.db.QueryRowContext(ctx, query, id).Scan(
		&notif.ID,
		&notif.UsrID,
		&notif.ChatID,
		&notif.Status,
		&notif.GoalID,
		&notif.Notify,
		&notif.LastUpdated,
	)
	if err != nil {
		return models.Notifier{}, err
	}

	return notif, nil
}

func (n *notifier) List(ctx context.Context, pars models.NotifierPars) ([]models.Notifier, int64, error) {
	query := "SELECT id, usr_id, chat_id, status_id, goal_id, notify, last_updated FROM notifier WHERE"
	var args []interface{}

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

	rows, err := n.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}

	var notifs []models.Notifier
	for rows.Next() {
		var notif models.Notifier

		err := rows.Scan(
			&notif.ID,
			&notif.UsrID,
			&notif.ChatID,
			&notif.Status,
			&notif.GoalID,
			&notif.Notify,
			&notif.LastUpdated,
		)
		if err != nil {
			return nil, 0, err
		}

		notifs = append(notifs, notif)
	}

	return notifs, 0, nil
}

func (n *notifier) Update(ctx context.Context, obj models.NotifierCU, id int64) error {
	updateValues := []interface{}{id}
	queryUpdate := ` UPDATE notifier`
	querySet := ` SET id = $1`
	queryWhere := ` WHERE id = $1`

	if obj.Notify != nil {
		updateValues = append(updateValues, time.Now().Add(*obj.Notify))
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

	if obj.GoalID != nil {
		updateValues = append(updateValues, *obj.GoalID)
		querySet += ` , goal_id = $` + strconv.Itoa(len(updateValues))
	}

	fmt.Println(queryUpdate, querySet, updateValues)
	_, err := n.db.ExecContext(ctx, queryUpdate+querySet+queryWhere, updateValues...)
	if err != nil {
		return err
	}

	return nil
}
