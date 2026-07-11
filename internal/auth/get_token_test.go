package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "bananaSplit"
	token, err := MakeJWT(userID, tokenSecret, time.Minute)
	if err != nil {
		t.Fatalf("fail to make JWT: %v", err)
	}

	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)

	_, err = GetBearerToken(header)
	if err != nil {
		t.Fatalf("Fail to get Bearer Token: %v", err)
	}

}
