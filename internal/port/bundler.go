package port

import "context"

type Bundler interface {
	Bundle(ctx context.Context) (string, error)
}
