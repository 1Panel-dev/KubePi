package server

import (
	"context"
)

// RunOptions holds a set of configurations for Server.Run().
type RunOptions struct {
	gracefulCtx context.Context
}

// RunOption is an option of Server.Run().
type RunOption func(*RunOptions)

// WithGracefulContext accepts a context to shutdown a Server
// with care for existing client connections.
func WithGracefulContext(ctx context.Context) RunOption {
	return func(options *RunOptions) {
		options.gracefulCtx = ctx
	}
}
