package application

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
)

const bearer = "Bearer"

// A JSON Web Token middleware
func getJwtAuthMiddleware(key string) echo.HandlerFunc {
	return func(c *echo.Context) error {

		// Skip WebSocket
		if (c.Request().Header.Get(echo.Upgrade)) == echo.WebSocket {
			return nil
		}

		auth := c.Request().Header.Get("Authorization")
		l := len(bearer)

		if len(auth) > l+1 && auth[:l] == bearer {
			_, err := jwt.Parse(auth[l+1:], func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(key), nil // key must not come from token itself
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, NewApiError(http.StatusUnauthorized, err.Error(), nil))
			}

			return nil
		}

		return c.JSON(http.StatusUnauthorized, ErrorUnautorized)
	}
}
