package invitation_token

import (
	"catalyst/models"
	"catalyst/requests"
)

func (invitationTokenService *Service) Update(invitationToken *models.Token, updatePostRequest *requests.UpdateTokenRequest) {
	invitationToken.Token = updatePostRequest.Token
	invitationToken.Active = updatePostRequest.Active
	invitationTokenService.DB.Save(invitationToken)
}
