package rcon

import "context"

type IRcon interface {
	Start(ctx context.Context) error
}
