package auth_token

import (
	"catalyst/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (authTokenService *Service) CreateRefreshToken(user *models.User) (t string, err error) {
	claimsRefresh := &JwtCustomRefreshClaims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * ExpireRefreshCount).Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)

	rt, err := refreshToken.SignedString([]byte(authTokenService.config.Auth.RefreshSecret))
	if err != nil {
		return "", err
	}
	return rt, err
}
