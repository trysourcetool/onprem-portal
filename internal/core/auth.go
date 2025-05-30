package core

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/trysourcetool/onprem-portal/internal/config"
)

const (
	tokenExpiration          = time.Duration(60) * time.Minute
	tokenExpirationDev       = time.Duration(365*24) * time.Hour
	RefreshTokenExpiration   = time.Duration(30*24) * time.Hour
	XSRFTokenExpiration      = time.Duration(30*24) * time.Hour
	RefreshTokenMaxAgeBuffer = time.Duration(7*24) * time.Hour
	TmpTokenExpiration       = time.Duration(30) * time.Minute
)

func TokenExpiration() time.Duration {
	if config.Config.Env == config.EnvLocal {
		return tokenExpirationDev
	}
	return tokenExpiration
}

func GenerateRefreshToken() (plainRefreshToken, hashedRefreshToken string, err error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", err
	}

	plainRefreshToken = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)

	hashedRefreshToken = HashRefreshToken(plainRefreshToken)

	return plainRefreshToken, hashedRefreshToken, nil
}

func HashRefreshToken(plainRefreshToken string) string {
	hash := sha256.Sum256([]byte(plainRefreshToken))
	return hex.EncodeToString(hash[:])
}
