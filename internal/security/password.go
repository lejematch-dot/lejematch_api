package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Memory      uint32 = 64 * 1024
	argon2Iterations  uint32 = 3
	argon2Parallelism uint8  = 2
	argon2SaltLength         = 16
	argon2KeyLength          = 32
)

var (
	ErrInvalidPasswordHash = errors.New("invalid password hash")
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
)

func HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", ErrPasswordTooShort
	}

	salt := make([]byte, argon2SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, argon2Iterations, argon2Memory, argon2Parallelism, argon2KeyLength)

	b64 := base64.RawStdEncoding
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argon2Memory,
		argon2Iterations,
		argon2Parallelism,
		b64.EncodeToString(salt),
		b64.EncodeToString(hash),
	), nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	params, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, uint32(len(hash)))
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

type hashParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
}

func decodeHash(encodedHash string) (hashParams, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return hashParams{}, nil, nil, ErrInvalidPasswordHash
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return hashParams{}, nil, nil, ErrInvalidPasswordHash
	}

	params := hashParams{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.memory, &params.iterations, &params.parallelism); err != nil {
		return hashParams{}, nil, nil, ErrInvalidPasswordHash
	}

	b64 := base64.RawStdEncoding

	salt, err := b64.DecodeString(parts[4])
	if err != nil {
		return hashParams{}, nil, nil, ErrInvalidPasswordHash
	}

	hash, err := b64.DecodeString(parts[5])
	if err != nil {
		return hashParams{}, nil, nil, ErrInvalidPasswordHash
	}

	return params, salt, hash, nil
}
