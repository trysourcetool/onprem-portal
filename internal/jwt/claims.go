package jwt

import "github.com/golang-jwt/jwt/v5"

const (
	Issuer = "portal.trysourcetool.com"

	UserSignatureSubjectMagicLink = "magic_link"
)

type UserMagicLinkRegistrationClaims struct {
	Email string
	jwt.RegisteredClaims
}
