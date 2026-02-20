package google

import (
	"context"

	"google.golang.org/api/idtoken"
)

type GoogleTokenValidator interface {
	Validate(ctx context.Context, token string) (*idtoken.Payload, error)
}
