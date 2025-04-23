package jwt

import "github.com/golang-jwt/jwt/v5"

const issuer = "https://portal.trysourcetool.com"

type MagicLinkRegistrationClaims struct {
	jwt.RegisteredClaims
}
