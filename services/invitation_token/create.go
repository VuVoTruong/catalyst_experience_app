package invitation_token

import (
	"catalyst/models"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func randomCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (invitationTokenService *Service) Create() *models.Token {
	minLengthCode := 6
	maxLengthCode := 12
	ranCodeLength := minLengthCode + rand.Intn(maxLengthCode-minLengthCode)
	token := randomCode(ranCodeLength)
	//TODO: check if token unique here
	invitationToken := models.Token{
		Token:  token,
		Active: 1,
	}
	invitationTokenService.DB.Create(&invitationToken)
	return &invitationToken
}
