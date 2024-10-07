package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mrbelka12000/goals_scheduler/internal/models"
)

type notify struct {
	db *sql.DB
}

func newNotify(db *sql.DB) *notify {
	return &notify{
		db: db,
	}
}

func (n *notify) Create(ctx context.Context, obj models.NotifyCU) (int64, error) {
	query := `
INSERT INTO notify (
	hour, 
	minute, 
	weekday, 
	goal_id) VALUES (
	$1, 
	$2, 
	$3, 
	$4
) RETURNING id`

	var id int64
	err := n.db.QueryRowContext(ctx, query, *obj.Hour, *obj.Minute, *obj.Weekday, *obj.GoalID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create notify: %w", err)
	}

	return id, nil
}

func (n *notify) Get(ctx context.Context, pars models.NotifyPars) (empty models.Notify, err error) {
	query := `
SELECT 
	id,
	hour,
	minute,
	weekday,
	goal_id 
FROM notify WHERE`

	var args []interface{}

	if pars.ID != nil {
		args = append(args, *pars.ID)
		query += fmt.Sprintf(" id = $%v AND", len(args))
	}
	if pars.GoalID != nil {
		args = append(args, *pars.GoalID)
		query += fmt.Sprintf(" goal_id = $%v AND", len(args))
	}

	if pars.WeekDay != nil {
		args = append(args, *pars.WeekDay)
		query += fmt.Sprintf(" weekday = $%v AND", len(args))
	}

	query = query[:len(query)-4] // Remove the trailing " AND"

	var obj models.Notify

	row := n.db.QueryRowContext(ctx, query, args...)

	err = row.Scan(&obj.ID, &obj.Hour, &obj.Minute, &obj.Weekday, &obj.GoalID)
	if err != nil {
		return empty, fmt.Errorf("get notify: %w", err)
	}

	return obj, row.Err()
}
