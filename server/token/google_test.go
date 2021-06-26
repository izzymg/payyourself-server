package token

import (
	"context"
	"testing"
)

func TestGoogleTokenChecker(t *testing.T) {
	checker := MakeGoogleTokenChecker("fake")
	isValid := checker.TokenIsValid(context.Background(), "a")
	if isValid {
		t.Fatal("expected invalid token")
	}
}
