package task

import (
	"fmt"
	"time"
)

const (
	maxRetries        = 10
	retryDelay        = 10 * time.Second
	connectionRefused = "connect: connection refused"
)

func retryWithBackoff[T any](operationName string, fetchFunc func() ([]T, error)) ([]T, error) {
	var result []T
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		result, err = fetchFunc()
		if err == nil {
			return result, nil
		}

		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to complete %s after %d attempts: %w", operationName, maxRetries, err)
}
