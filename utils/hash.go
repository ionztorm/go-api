package utils

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

const (
	timer   = 1
	memory  = 64 * 1024
	threads = 4
	keyLen  = 32
)

func HashPassword(password string) (string, error) {
	bytes, error := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), error
}

func CompareHash(password string, hashedPassword string) bool {
	error := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return error == nil
}

func HashPasswordArgon2id(password string) (string, string, error) {
	salt := make([]byte, 16)

	_, error := rand.Read(salt)
	if error != nil {
		return "", "", error
	}

	hash := argon2.IDKey([]byte(password), salt, timer, memory, uint8(threads), keyLen)

	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)

	return b64Hash, b64Salt, nil
}

func CompareHashArgon2id(password, b64Salt, hashedPassword string) bool {
	salt, err := base64.RawStdEncoding.DecodeString(b64Salt)
	if err != nil {
		return false
	}
	hash := argon2.IDKey([]byte(password), salt, timer, memory, uint8(threads), keyLen)

	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	return b64Hash == hashedPassword
}
