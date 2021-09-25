package mock

import (
	"strings"
	"user-access-management/jwtparser"
)

//TokenManagerMock ...
type TokenManagerMock struct {
	Err error
}

//ExtractJWTClaims ...
func (t *TokenManagerMock) ExtractJWTClaims(bearerToken string) (*jwtparser.AuthTokenClaim, error) {
	if t.Err != nil {
		return nil, t.Err
	}
	claims := strings.Split(bearerToken, " ")
	return &jwtparser.AuthTokenClaim{
		User: jwtparser.User{
			Username: claims[0],
			Email:    claims[1],
			Role:     claims[2],
		},
	}, nil
}
