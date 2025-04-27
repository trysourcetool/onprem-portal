package jwt

import "github.com/golang-jwt/jwt/v5"

const issuer = "https://portal.trysourcetool.com"

type AuthClaims struct {
	XSRFToken string
	jwt.RegisteredClaims
}

type MagicLinkClaims struct {
	jwt.RegisteredClaims
}

type MagicLinkRegistrationClaims struct {
	jwt.RegisteredClaims
}
