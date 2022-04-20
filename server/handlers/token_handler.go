package handlers

import (
	"catalyst/models"
	"catalyst/repositories"
	"catalyst/requests"
	"catalyst/responses"
	s "catalyst/server"
	invitationTokenService "catalyst/services/invitation_token"
	"fmt"

	"net/http"

	"github.com/labstack/echo/v4"
)

type TokenHandlers struct {
	server *s.Server
}

func NewTokenHandlers(server *s.Server) *TokenHandlers {
	return &TokenHandlers{server: server}
}

// CreateToken godoc
// @Summary Create token
// @Description Create token
// @Tags Tokens Actions
// @Accept json
// @Produce json
// @Success 201 {object} responses.Data
// @Failure 400 {object} responses.Error
// @Security ApiKeyAuth
// @Router /tokens [post]
func (p *TokenHandlers) CreateToken(c echo.Context) error {
	invitationTokenService := invitationTokenService.NewInvitationTokenService(p.server.DB)
	invitationToken := invitationTokenService.Create()
	response := responses.NewTokenResponse(*invitationToken)
	return responses.Response(c, http.StatusCreated, response)
}

// GetTokens godoc
// @Summary Get invitationTokens
// @Description Get the list of all invitationTokens
// @Tags Tokens Actions
// @Produce json
// @Success 200 {array} responses.TokenResponse
// @Security ApiKeyAuth
// @Router /tokens [get]
func (p *TokenHandlers) GetTokens(c echo.Context) error {
	var invitationTokens []models.Token

	tokenRepository := repositories.NewTokenRepository(p.server.DB)
	tokenRepository.GetTokens(&invitationTokens)

	response := responses.NewTokensResponse(invitationTokens)
	return responses.Response(c, http.StatusOK, response)
}

// UpdateToken godoc
// @Summary Update token
// @Description Update token
// @Tags Tokens Actions
// @Accept json
// @Produce json
// @Param token path string true "Token"
// @Param params body requests.UpdateTokenRequest true "Token body"
// @Success 200 {object} responses.Data
// @Failure 400 {object} responses.Error
// @Security ApiKeyAuth
// @Router /tokens/{token} [put]
func (p *TokenHandlers) UpdateToken(c echo.Context) error {
	updateTokenRequest := new(requests.UpdateTokenRequest)
	tokenStr := c.Param("token")
	if err := c.Bind(updateTokenRequest); err != nil {
		return err
	}
	if err := updateTokenRequest.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty")
	}
	invitationToken := models.Token{}
	tokenRepository := repositories.NewTokenRepository(p.server.DB)
	tokenRepository.GetToken(&invitationToken, tokenStr)
	fmt.Println(invitationToken)
	if invitationToken.Token == "" {
		return responses.ErrorResponse(c, http.StatusNotFound, "Token not found")
	}
	invitationTokenService := invitationTokenService.NewInvitationTokenService(p.server.DB)
	invitationTokenService.Update(&invitationToken, updateTokenRequest)

	return responses.MessageResponse(c, http.StatusOK, "Token successfully updated")
}

// ValidateToken godoc
// @Summary Validate token
// @Description Validate token
// @Tags Tokens Actions
// @Accept json
// @Produce json
// @Param token path string true "Token"
// @Success 200 {object} responses.Data
// @Failure 400 {object} responses.Error
// @Security ApiKeyAuth
// @Router /tokens/validate/{token} [post]
func (p *TokenHandlers) ValidateToken(c echo.Context) error {
	tokenStr := c.Param("token")
	invitationTokenService := invitationTokenService.NewInvitationTokenService(p.server.DB)
	isValid := invitationTokenService.Validate(tokenStr)
	if !isValid {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Token is invalid")
	}
	return responses.MessageResponse(c, http.StatusOK, "Token is valid")
}
