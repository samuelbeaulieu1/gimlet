package gimlet

import "github.com/golang-jwt/jwt/v4"

type AuthTokenPayload struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func (payload AuthTokenPayload) Valid() error {
	return payload.StandardClaims.Valid()
}
