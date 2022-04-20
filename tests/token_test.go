package tests

import (
	"catalyst/requests"
	"catalyst/server"
	"catalyst/server/handlers"
	"catalyst/services/auth_token"
	"catalyst/tests/helpers"
	"net/http"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const token = "ABCXYZ"
const tokenNotExists = "ABCDEF"

func TestWalkPostsCrud(t *testing.T) {
	requestGenerate := helpers.Request{
		Method: http.MethodPost,
		Url:    "/tokens",
	}
	requestGet := helpers.Request{
		Method: http.MethodGet,
		Url:    "/tokens",
	}
	requestUpdate := helpers.Request{
		Method: http.MethodPut,
		Url:    "/tokens/" + token,
		PathParam: &helpers.PathParam{
			Name:  "token",
			Value: token,
		},
	}
	requestValidate := helpers.Request{
		Method: http.MethodPost,
		Url:    "/tokens/validate/" + token,
		PathParam: &helpers.PathParam{
			Name:  "token",
			Value: token,
		},
	}

	handlerFuncCreate := func(s *server.Server, c echo.Context) error {
		return handlers.NewTokenHandlers(s).CreateToken(c)
	}
	handlerFuncGet := func(s *server.Server, c echo.Context) error {
		return handlers.NewTokenHandlers(s).GetTokens(c)
	}
	handlerFuncUpdate := func(s *server.Server, c echo.Context) error {
		return handlers.NewTokenHandlers(s).UpdateToken(c)
	}
	handlerFuncValidate := func(s *server.Server, c echo.Context) error {
		return handlers.NewTokenHandlers(s).ValidateToken(c)
	}

	claims := &auth_token.JwtCustomClaims{
		Name: "user",
		ID:   helpers.UserId,
	}
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	cases := []helpers.TestCase{
		//Generate token successfully
		{
			"Generate token successfully",
			requestGenerate,
			nil,
			handlerFuncCreate,
			&helpers.QueryMock{
				Query: `INSERT INTO "tokens" ("token","active","created_at","updated_at") VALUES (?,?,?,?)`,
				Reply: nil,
			},
			helpers.ExpectedResponse{
				StatusCode: 201,
			},
		},
		//Get token
		{
			"Get tokens success",
			requestGet,
			"",
			handlerFuncGet,
			&helpers.QueryMock{
				Query: `SELECT * FROM "tokens"`,
				Reply: helpers.MockReply{{"token": token, "active": 1, "created_at": time.Now()}},
			},
			helpers.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "[{\"token\":\"ABCXYZ\",\"active\":1}]",
			},
		},
		//Update token
		{
			"Update token success",
			requestUpdate,
			requests.UpdateTokenRequest{
				BasicToken: requests.BasicToken{
					Token:  token,
					Active: 1,
				},
			},
			handlerFuncUpdate,
			&helpers.QueryMock{
				Query: `SELECT * FROM "tokens"  WHERE (token = ABCXYZ)`,
				Reply: helpers.MockReply{{"token": token, "active": 1, "created_at": time.Now()}},
			},
			helpers.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "Token successfully updated",
			},
		},
		//Update non-existed token
		{
			"Update non-existed token",
			helpers.Request{
				Method: http.MethodPut,
				Url:    "/tokens/" + tokenNotExists,
				PathParam: &helpers.PathParam{
					Name:  "token",
					Value: tokenNotExists,
				},
			},
			requests.UpdateTokenRequest{
				BasicToken: requests.BasicToken{
					Token:  tokenNotExists,
					Active: 0,
				},
			},
			handlerFuncUpdate,
			&helpers.QueryMock{
				Query: `SELECT * FROM "tokens"  WHERE (token = ABCDEF)`,
				Reply: helpers.MockReply{nil},
			},
			helpers.ExpectedResponse{
				StatusCode: 404,
				BodyPart:   "Token not found",
			},
		},
		//Validate token successfully
		{
			"Validate valid token",
			requestValidate,
			nil,
			handlerFuncValidate,
			&helpers.QueryMock{
				Query: `SELECT * FROM "tokens"  WHERE (token = ABCXYZ)`,
				Reply: helpers.MockReply{{"token": token, "active": 1, "created_at": time.Now()}},
			},
			helpers.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "Token is valid",
			},
		},
		//Validate inactive token
		{
			"Validate inactive token",
			requestValidate,
			nil,
			handlerFuncValidate,
			&helpers.QueryMock{
				Query: `SELECT * FROM "tokens"  WHERE (token = ABCXYZ)`,
				Reply: helpers.MockReply{{"token": token, "active": 0, "created_at": time.Now()}},
			},
			helpers.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "Token is invalid",
			},
		},
		//Validate expired token
		{
			"Validate inactive token",
			requestValidate,
			nil,
			handlerFuncValidate,
			&helpers.QueryMock{
				Query: `SELECT * FROM "tokens"  WHERE (token = ABCXYZ)`,
				Reply: helpers.MockReply{{"token": token, "active": 1, "created_at": time.Now().AddDate(0, 0, -8)}},
			},
			helpers.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "Token is invalid",
			},
		},
	}

	s := helpers.NewServer()

	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			c, recorder := helpers.PrepareContextFromTestCase(s, test)
			c.Set("user", validToken)

			if assert.NoError(t, test.HandlerFunc(s, c)) {
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
				assert.Equal(t, test.Expected.StatusCode, recorder.Code)
			}
		})
	}
}
