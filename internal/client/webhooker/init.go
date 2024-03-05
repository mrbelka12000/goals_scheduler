package webhooker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"goals_scheduler/pkg/config"
)

type (
	Client struct {
		cl           *http.Client
		webHookerURL string
		callbackUrl  string
	}
	CreateWebHookRequest struct {
		Method  string
		Params  map[string]string
		Body    string
		EndTime time.Time
	}
)

const (
	defaultCallbackPath = "/callback"
)

func NewClient(cfg config.Config) *Client {
	return &Client{
		cl:           http.DefaultClient,
		webHookerURL: cfg.WebHookerURL,
		callbackUrl:  cfg.ServiceName,
	}
}

func (c *Client) CreateWebHook(ctx context.Context, req CreateWebHookRequest) error {
	var (
		in = struct {
			CallbackURL string            `json:"callback_url,omitempty"`
			Method      string            `json:"method,omitempty"` // default POST
			Params      map[string]string `json:"params,omitempty"`
			Body        string            `json:"body,omitempty"`
			EndTime     time.Time         `json:"end_time,omitempty"`
		}{
			CallbackURL: fmt.Sprintf("%v%v", c.callbackUrl, defaultCallbackPath),
			Method:      req.Method,
			Params:      req.Params,
			Body:        req.Body,
			EndTime:     req.EndTime,
		}
	)

	reqBody, err := json.Marshal(in)
	if err != nil {
		return err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webHookerURL, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	_, err = c.cl.Do(r)
	if err != nil {
		return err
	}

	return nil
}
