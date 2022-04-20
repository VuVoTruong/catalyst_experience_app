package repositories

import (
	"catalyst/models"

	"github.com/jinzhu/gorm"
)

type TokenRepositoryQ interface {
	GetTokens(Tokens *[]models.Token)
	GetToken(Token *models.Token, id int)
}

type TokenRepository struct {
	DB *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{DB: db}
}

func (TokenRepository *TokenRepository) GetTokens(tokens *[]models.Token) {
	TokenRepository.DB.Find(tokens)
}

func (TokenRepository *TokenRepository) GetToken(token *models.Token, tokenString string) {
	TokenRepository.DB.Where("token = ?", tokenString).Find(token)
}
