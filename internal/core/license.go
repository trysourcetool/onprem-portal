package core

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

type License struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Key       string    `db:"key"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func GenerateLicenseKey() (plainKey, hashedKey string, err error) {
	const randomBytesLen = 20 // 160 bits → 32 base32 chars → 8 groups of 4 after formatting

	b := make([]byte, randomBytesLen)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}

	encoded := strings.ToUpper(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b))

	var sb strings.Builder
	for i, r := range encoded {
		if i > 0 && i%4 == 0 {
			sb.WriteRune('-')
		}
		sb.WriteRune(r)
	}
	plainKey = sb.String()

	stripped := strings.ReplaceAll(plainKey, "-", "")
	hash := sha256.Sum256([]byte(stripped))
	hashedKey = hex.EncodeToString(hash[:])
	return plainKey, hashedKey, nil
}

func HashLicenseKey(plainKey string) string {
	stripped := strings.ReplaceAll(plainKey, "-", "")
	h := sha256.Sum256([]byte(stripped))
	return hex.EncodeToString(h[:])
}
