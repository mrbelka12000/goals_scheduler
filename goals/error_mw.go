package goals

import (
	"context"
)

type (
	ErrorMW struct {
		srv Service
		es  errorSender
	}

	errorSender interface {
		SendError(error)
	}
)

func NewErrorMW(srv Service, es errorSender) *ErrorMW {
	return &ErrorMW{
		srv: srv,
		es:  es,
	}
}

func (e *ErrorMW) Create(ctx context.Context, obj Goal) (int64, error) {
	out, err := e.srv.Create(ctx, obj)
	if err != nil {
		e.es.SendError(err)
		return 0, err
	}

	return out, nil
}

func (e *ErrorMW) Get(ctx context.Context, id int64) (Goal, error) {
	out, err := e.srv.Get(ctx, id)
	if err != nil {
		e.es.SendError(err)
		return out, err
	}
	return out, nil
}

func (e *ErrorMW) List(ctx context.Context, pars GoalPars) ([]Goal, int64, error) {
	out, count, err := e.srv.List(ctx, pars)
	if err != nil {
		e.es.SendError(err)
		return nil, 0, err
	}
	return out, count, nil
}

func (e *ErrorMW) DeleteAllOfUsers(ctx context.Context, chatID string) error {
	err := e.srv.DeleteAllOfUsers(ctx, chatID)
	if err != nil {
		e.es.SendError(err)
		return err
	}
	return nil
}

func (e *ErrorMW) Update(ctx context.Context, obj Goal, id int64) error {
	err := e.srv.Update(ctx, obj, id)
	if err != nil {
		e.es.SendError(err)
		return err
	}
	return nil
}

func (e *ErrorMW) Delete(ctx context.Context, id int64) error {
	err := e.srv.Delete(ctx, id)
	if err != nil {
		e.es.SendError(err)
		return err
	}
	return nil
}
