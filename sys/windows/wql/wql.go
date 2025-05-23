//go:build windows

package wql

import (
	"context"
	"time"

	"github.com/yusufpapurcu/wmi"
)

var (
	Timeout = 3 * time.Second
)

// Depracted
func WMIQuery(ctx context.Context, query string, dst any, connectServerArgs ...any) error {
	if _, ok := ctx.Deadline(); !ok {
		ctxTimeout, cancel := context.WithTimeout(ctx, Timeout)
		defer cancel()
		ctx = ctxTimeout
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- wmi.Query(query, dst, connectServerArgs...)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
