package invitation_token

import (
	"catalyst/models"
	"catalyst/requests"

	"github.com/jinzhu/gorm"
)

type ServiceWrapper interface {
	Create(token *models.Token)
	Update(token *models.Token, updateTokenRequest *requests.UpdateTokenRequest)
}

type Service struct {
	DB *gorm.DB
}

func NewInvitationTokenService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
