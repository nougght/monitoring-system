package util

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHashString = errors.New("invalid hash string")
)

type HashParams struct {
	Memory      int
	Iterations  int
	Parallelism int
	SaltLength  int
	KeyLength   int
}

func defaultHashParams() *HashParams {
	return &HashParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}
func makeHashString(hash, salt string, params *HashParams) string {
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		salt,
		hash,
	)
}
func parseHashString(hashString string) (hash, salt string, params *HashParams, err error) {
	err = fmt.Errorf("can't parse hash string: %w", ErrInvalidHashString)
	params = &HashParams{}

	parts := strings.Split(hashString, "$")
	if len(parts) != 6 {
		return "", "", nil, err
	}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return "", "", nil, err
	}
	salt = parts[4]
	hash = parts[5]
	return hash, salt, params, nil
}

func Hash(value string, params *HashParams) (hashString string, err error) {
	if params == nil {
		params = defaultHashParams()
	}
	salt := make([]byte, params.SaltLength)
	if _, err = rand.Read(salt); err != nil {
		return
	}
	hashBytes := argon2.IDKey([]byte(value), salt, uint32(params.Iterations), uint32(params.Memory), uint8(params.Parallelism), uint32(params.KeyLength))
	hashB64 := base64.StdEncoding.EncodeToString(hashBytes)
	saltB64 := base64.StdEncoding.EncodeToString(salt)
	return makeHashString(hashB64, saltB64, params), nil
}

func CompareHash(hashString, value string) (bool, error) {
	expectedHashB64, saltB64, params, err := parseHashString(hashString)
	if err != nil {
		return false, err
	}

	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.StdEncoding.DecodeString(expectedHashB64)
	if err != nil {
		return false, err
	}

	params.KeyLength = len(expectedHash)

	hash := argon2.IDKey([]byte(value), []byte(salt), uint32(params.Iterations), uint32(params.Memory), uint8(params.Parallelism), uint32(params.KeyLength))

	// compare hash in constant time
	return subtle.ConstantTimeCompare(expectedHash, hash) == 1, nil
}
