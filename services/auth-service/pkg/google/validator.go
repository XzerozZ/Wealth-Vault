package google

import (
	"context"

	"google.golang.org/api/idtoken"
)

type GoogleValidator struct {
	clientID string
}

func NewGoogleValidator(clientID string) *GoogleValidator {
	return &GoogleValidator{
		clientID: clientID,
	}
}

func (g *GoogleValidator) Validate(ctx context.Context, token string) (*idtoken.Payload, error) {
	return idtoken.Validate(ctx, token, g.clientID)
}
