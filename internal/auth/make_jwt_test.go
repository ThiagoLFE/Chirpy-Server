package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	tokenSecret := "paoComBananinha"
	userID := uuid.New()

	token, err := MakeJWT(userID, tokenSecret, time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT() returned an error: %v", err)
	}

	recoveredUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT() returned an error: %v", err)
	}

	if userID != recoveredUserID {
		t.Fatalf("fail to recover the ID from ValidateJWT. %s != %s", userID, recoveredUserID)
	}
}
