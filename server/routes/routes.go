package routes

import (
	s "catalyst/server"
	"catalyst/server/handlers"
	"catalyst/services/auth_token"
	"fmt"

	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigureRoutes(server *s.Server) {
	invitationTokenHandler := handlers.NewTokenHandlers(server)
	authHandler := handlers.NewAuthHandler(server)
	registerHandler := handlers.NewRegisterHandler(server)

	server.Echo.Use(middleware.Logger())

	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	server.Echo.POST("/login", authHandler.Login)
	server.Echo.POST("/register", registerHandler.Register)
	server.Echo.POST("/refresh", authHandler.RefreshToken)

	fmt.Println(server.Config.Auth.AccessSecret)

	server.Echo.POST("/tokens/validate/:token", invitationTokenHandler.ValidateToken)
	r := server.Echo.Group("")
	config := middleware.JWTConfig{
		Claims:     &auth_token.JwtCustomClaims{},
		SigningKey: []byte(server.Config.Auth.AccessSecret),
	}
	r.Use(middleware.JWTWithConfig(config))

	r.GET("/tokens", invitationTokenHandler.GetTokens)
	r.POST("/tokens", invitationTokenHandler.CreateToken)
	r.PUT("/tokens/:token", invitationTokenHandler.UpdateToken)
}
