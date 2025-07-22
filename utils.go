package metabase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"
)

// DecodeJSON decodes JSON response into target structure
func DecodeJSON(data []byte, out any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode(out)
}

// InferFieldsFromJSON returns the list of fields from JSON array of objects
func InferFieldsFromJSON(jsonData []byte) ([]string, error) {
	var records []map[string]any
	if err := json.Unmarshal(jsonData, &records); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	fieldSet := make(map[string]struct{})
	for _, rec := range records {
		for k := range rec {
			fieldSet[k] = struct{}{}
		}
	}
	var fields []string
	for f := range fieldSet {
		fields = append(fields, f)
	}
	sort.Strings(fields)
	return fields, nil
}

// WithRetry wraps a function with retry logic
func WithRetry(ctx context.Context, maxAttempts int, delay time.Duration, fn func() error) error {
	if maxAttempts <= 0 {
		return errors.New("invalid retry attempts")
	}
	var lastErr error
	for range maxAttempts {
		if err := fn(); err != nil {
			lastErr = err
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
				continue
			}
		} else {
			return nil
		}
	}
	return fmt.Errorf("after %d attempts, last error: %w", maxAttempts, lastErr)
}
