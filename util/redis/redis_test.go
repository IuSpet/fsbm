package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetWithRetry(t *testing.T) {
	key := "test_key"
	value := "test_value"
	ctx := context.Background()
	res, err := GetWithRetry(ctx, key)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, value, res)
}

func TestSetWithRetry(t *testing.T) {
	key := "test_key"
	value := "test_value"
	ctx := context.Background()
	err := SetWithRetry(ctx, key, value, 30*time.Minute)
	if err != nil {
		panic(err)
	}
}
