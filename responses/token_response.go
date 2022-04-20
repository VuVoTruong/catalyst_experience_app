package responses

import (
	"catalyst/models"
)

type TokenResponse struct {
	Token  string `json:"token" example:"ABCXYZ"`
	Active int16  `json:"active" example:"1"`
}

func NewTokensResponse(tokens []models.Token) *[]TokenResponse {
	tokenResponse := make([]TokenResponse, 0)

	for i := range tokens {
		tokenResponse = append(tokenResponse, TokenResponse{
			Token:  tokens[i].Token,
			Active: tokens[i].Active,
		})
	}

	return &tokenResponse
}

func NewTokenResponse(token models.Token) *TokenResponse {
	tokenResponse := TokenResponse{
		Token:  token.Token,
		Active: token.Active,
	}

	return &tokenResponse
}
