package invitation_token

import (
	"catalyst/models"
	"time"
)

func (invitationTokenService *Service) Validate(tokenString string) bool {
	var token models.Token
	invitationTokenService.DB.Where("token = ?", tokenString).Find(&token)
	//Check token
	checkToken := &token
	if checkToken == nil || checkToken.Active != 1 {
		return false
	}
	//Check token auto expired
	expiredDate := checkToken.CreatedAt.AddDate(0, 0, 7)
	if expiredDate.Before(time.Now()) {
		return false
	}
	return true
}
