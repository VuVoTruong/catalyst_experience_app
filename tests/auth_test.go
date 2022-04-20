package tests

import (
	"catalyst/requests"
	"catalyst/responses"
	"catalyst/server"
	"catalyst/server/handlers"
	"catalyst/tests/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestWalkAuth(t *testing.T) {
	request := helpers.Request{
		Method: http.MethodPost,
		Url:    "/login",
	}
	handlerFunc := func(s *server.Server, c echo.Context) error {
		return handlers.NewAuthHandler(s).Login(c)
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	commonMock := &helpers.QueryMock{
		Query: `SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((email = name@test.com))`,
		Reply: helpers.MockReply{{"id": helpers.UserId, "email": "name@test.com", "name": "User Name", "password": encryptedPassword}},
	}

	cases := []helpers.TestCase{
		{
			"Auth success",
			request,
			requests.LoginRequest{
				BasicAuth: requests.BasicAuth{
					Email:    "name@test.com",
					Password: "password",
				},
			},
			handlerFunc,
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "",
			},
		},
		{
			"Login attempt with incorrect password",
			request,
			requests.LoginRequest{
				BasicAuth: requests.BasicAuth{
					Email:    "name@test.com",
					Password: "incorrectPassword",
				},
			},
			handlerFunc,
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 401,
				BodyPart:   "Invalid credentials",
			},
		},
		{
			"Login attempt as non-existent user",
			request,
			requests.LoginRequest{
				BasicAuth: requests.BasicAuth{
					Email:    "user.not.exists@test.com",
					Password: "password",
				},
			},
			handlerFunc,
			commonMock,
			helpers.ExpectedResponse{
				StatusCode: 401,
				BodyPart:   "Invalid credentials",
			},
		},
	}

	s := helpers.NewServer()

	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			c, recorder := helpers.PrepareContextFromTestCase(s, test)

			if assert.NoError(t, test.HandlerFunc(s, c)) {
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
				if assert.Equal(t, test.Expected.StatusCode, recorder.Code) {
					if recorder.Code == http.StatusOK {
						assertTokenResponse(t, recorder)
					}
				}
			}
		})
	}
}

func assertTokenResponse(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	var authResponse responses.LoginResponse
	_ = json.Unmarshal([]byte(recorder.Body.String()), &authResponse)

	assert.Equal(t, float64(helpers.UserId), getUserIdFromToken(authResponse.AccessToken))
	assert.Equal(t, float64(helpers.UserId), getUserIdFromToken(authResponse.RefreshToken))
}

func getUserIdFromToken(tokenToParse string) float64 {
	token, _ := jwt.Parse(tokenToParse, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}
		var hmacSampleSecret []byte
		return hmacSampleSecret, nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)

	return claims["id"].(float64)
}
