package pool

import (
	"context"
	"net"

	"github.com/redis/go-redis/v9"
)

type HealthHook struct {
	*client
}

func newFailureHook(c *client) *HealthHook {
	return &HealthHook{
		client: c,
	}
}

func (hook *HealthHook) DialHook(next redis.DialHook) redis.DialHook { return next }

func (hook *HealthHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		err := next(ctx, cmd)
		if isNetworkError(err) {
			hook.onFailure()
		} else {
			hook.onSuccess()
		}
		return err
	}
}

func (hook *HealthHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		for _, cmd := range cmds {
			if isNetworkError(cmd.Err()) {
				hook.client.onFailure()
				return nil
			}
		}
		hook.onSuccess()

		return nil
	}
}

func isNetworkError(err error) bool {
	if err == nil {
		return false
	}
	// Network error
	_, ok := err.(net.Error)
	return ok
}
