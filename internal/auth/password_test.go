package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "bah!1234"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() returned an error: %v", err)
	}

	got, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() returned an error: %v", err)
	}
	if !got {
		t.Fatal("CheckPasswordHash() returned false for the original password")
	}

	got, err = CheckPasswordHash("wrong-password", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() returned an error for a wrong password: %v", err)
	}
	if got {
		t.Fatal("CheckPasswordHash() returned true for a wrong password")
	}
}
