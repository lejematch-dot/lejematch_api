package security

import "testing"

func TestHashAndVerifyPassword(t *testing.T) {
	password := "correct horse battery staple"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	match, err := VerifyPassword(password, hash)
	if err != nil {
		t.Fatalf("VerifyPassword() error = %v", err)
	}

	if !match {
		t.Fatal("VerifyPassword() = false, want true")
	}
}

func TestHashPasswordRejectsShortPassword(t *testing.T) {
	if _, err := HashPassword("short"); err != ErrPasswordTooShort {
		t.Fatalf("HashPassword() error = %v, want %v", err, ErrPasswordTooShort)
	}
}
