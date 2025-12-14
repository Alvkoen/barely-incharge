package httpclient

import (
	"net/http"
	"time"
)

const (
	DefaultTimeout = 30 * time.Second
)

func New() *http.Client {
	return &http.Client{
		Timeout: DefaultTimeout,
	}
}
