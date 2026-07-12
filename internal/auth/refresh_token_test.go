package auth

import "testing"

func TestRefreshToken(t *testing.T) {
	refresh_token := MakeRefreshToken()
	if refresh_token == "" {
		t.Fatalf("refresh token is empty: %v", refresh_token)
	}

}
