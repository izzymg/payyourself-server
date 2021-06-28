package token

import (
	"context"

	"google.golang.org/api/idtoken"
)

// GoogleTokenChecker is a TokenChecker which validates the token
// against google's oauth2 api.
type GoogleTokenChecker struct {
	clientId string
}

// TokenIsValid checks the given token against google's oauth api,
// using the provided clientId if any is given.
func (c GoogleTokenChecker) TokenIsValid(ctx context.Context, token string) (string, bool) {
	payload, err := idtoken.Validate(ctx, token, c.clientId)
	if err != nil {
		return "", false
	}

	return payload.Subject, true
}

// MakeGoogleTokenChecker returns a new GoogleTokenChecker
// Pass an empty string to disable validating against a clientId
func MakeGoogleTokenChecker(clientId string) GoogleTokenChecker {
	return GoogleTokenChecker{
		clientId,
	}
}
