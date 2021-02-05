package logs

import (
	"context"
	"testing"
)

func TestCtxInfo(t *testing.T) {
	ctx := context.WithValue(context.Background(), "req_id", "123456")
	SetGlobalLogLevel(Warn)
	CtxInfo(ctx, "test infoLogger log")
}
