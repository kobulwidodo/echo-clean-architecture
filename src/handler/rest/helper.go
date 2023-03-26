package rest

import (
	"errors"
	"fmt"
	"go-clean/src/business/entity"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func (r *rest) httpRespSuccess(ctx echo.Context, code int, message string, data interface{}) error {
	resp := entity.Response{
		Meta: entity.Meta{
			Message: message,
			Code:    code,
			IsError: false,
		},
		Data: data,
	}
	return ctx.JSON(code, resp)
}

func (r *rest) httpRespError(ctx echo.Context, code int, err error) error {
	resp := entity.Response{
		Meta: entity.Meta{
			Message: err.Error(),
			Code:    code,
			IsError: true,
		},
		Data: nil,
	}
	r.log.Error(err)
	return ctx.JSON(code, resp)
}

func (r *rest) VerifyUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return r.httpRespError(c, http.StatusUnauthorized, errors.New("no authorization header"))
			}

			var tokenString string
			_, err := fmt.Sscanf(authHeader, "Bearer %v", &tokenString)
			if err != nil {
				return r.httpRespError(c, http.StatusUnauthorized, errors.New("invalid token"))
			}

			token, err := r.ValidateToken(tokenString)
			if err != nil {
				return r.httpRespError(c, http.StatusUnauthorized, err)
			}

			claim, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return r.httpRespError(c, http.StatusUnauthorized, errors.New("failed to claim token"))
			}

			user := entity.User{}
			user, err = r.uc.User.GetById(uint(claim["id"].(float64)))
			if err != nil {
				return r.httpRespError(c, http.StatusUnauthorized, errors.New("error while getting user"))
			}

			ctx := c.Request().Context()
			ctx = r.auth.SetUserAuthInfo(ctx, user.ConvertToAuthUser(), tokenString)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func (r *rest) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("token invalid")
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
