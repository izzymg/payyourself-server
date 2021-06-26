package token

import (
	"context"
	"log"

	"google.golang.org/api/idtoken"
)

// GoogleTokenChecker is a TokenChecker which validates the token
// against google's oauth2 api.
type GoogleTokenChecker struct {
	clientId string
}

// TokenIsValid checks the given token against google's oauth api,
// using the provided clientId if any is given.
func (c GoogleTokenChecker) TokenIsValid(ctx context.Context, token string) bool {
	payload, err := idtoken.Validate(ctx, token, c.clientId)
	log.Println(payload)
	return err == nil
}

// MakeGoogleTokenChecker returns a new GoogleTokenChecker
// Pass an empty string to disable validating against a clientId
func MakeGoogleTokenChecker(clientId string) GoogleTokenChecker {
	return GoogleTokenChecker{
		clientId,
	}
}
