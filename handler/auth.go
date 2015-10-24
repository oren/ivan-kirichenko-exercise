package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
	"github.com/seesawlabs/ivan-kirichenko-exercise/lib"
	"golang.org/x/oauth2"
)

const defaultTokenExpiration time.Duration = 5 * time.Minute
const issuer string = "demoapp"
const bearer = "Bearer"

// TokenStorage defines some key-value storage for tokens by session id
type TokenStorage interface {
	Set(k string, x interface{}, d time.Duration)
	Get(k string) (interface{}, bool)
	Delete(k string)
}

// GetJwtAuthHandler creates a handler function that performs authorization
// based on JWT token in incoming request. Can be used as middleware
func GetJwtAuthHandler(jwtSecret string) echo.HandlerFunc {
	return func(c *echo.Context) error {

		// Skip WebSocket
		if (c.Request().Header.Get(echo.Upgrade)) == echo.WebSocket {
			return nil
		}

		auth := c.Request().Header.Get("Authorization")
		l := len(bearer)

		if len(auth) <= l+1 || auth[:l] != bearer {
			return echo.NewHTTPError(http.StatusUnauthorized,
				NewApiError("no or incorrect authorization token provided").String(),
			)
		}

		_, err := jwt.Parse(auth[l+1:], func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			if iss, ok := token.Claims["iss"].(string); !ok {
				return nil, errors.New("issuer not provided")
			} else if iss != issuer {
				return nil, errors.New("incorrect issuer provided")
			}

			accessToken, found := token.Claims["access_token"].(string)
			if !found {
				return nil, errors.New("access token not provided")
			}

			expirationTime, _ := token.Claims["exp"].(float64)
			if int64(expirationTime) < time.Now().Unix() {
				return nil, errors.New("access token has expired, try to authenticate again")
			}

			return getJwtSignature(jwtSecret, accessToken), nil
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized,
				NewApiError(err.Error()).String(),
			)
		}

		return nil
	}
}

// GetOAuthHandler creates a handler function that starts authorization process
// using Facebook as OAuth provider
func GetOAuthHandler(conf oauth2.Config, sessionSecret string, csrfStorage TokenStorage) echo.HandlerFunc {
	// oauthURLResponse is a type which is used only within oauth handler
	type oauthURLResponse struct {
		URL     string `json:"url"`
		Message string `json:"message"`
	}

	return func(c *echo.Context) error {
		// in general I must to ensure that sessionID is unique, but let's simplify
		// for test task
		sessionID, err := lib.GenerateRandomString(32)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		}

		csrfToken := generateCsrfToken(sessionID, sessionSecret)
		// in order to increase TTL of the cached value, let's save it as late as possible
		defer csrfStorage.Set(csrfToken, sessionID, defaultTokenExpiration)
		url := conf.AuthCodeURL(csrfToken, oauth2.AccessTypeOffline)

		return c.JSON(http.StatusOK, oauthURLResponse{url, "please use this url to authenticate with facebook"})
	}
}

// GetOAuthVerifyHandler creates a handler function that checks response of
// OAuth provider (Facebook), performs authorization of the user in our system
// and responds with JWT that should be used as access token to our API
func GetOAuthVerifyHandler(conf oauth2.Config, jwtSecret, sessionSecret string, csrfStorage TokenStorage) echo.HandlerFunc {
	// oauthVerifyResponse is a type which is used only within oauth verify handler
	type oauthVerifyResponse struct {
		Token   string `json:"jwt_token"`
		Expires int64  `json:"expires"`
	}

	return func(c *echo.Context) error {
		code := c.Query("code")
		errorMessage := c.Param("error_message")
		if code == "" {
			if errorMessage == "" {
				errorMessage = "no oauth code was provided"
			}
			return c.JSON(http.StatusBadRequest, NewApiError(errorMessage))
		}

		csrfToken := c.Query("state")
		if csrfToken == "" {
			return c.JSON(http.StatusBadRequest, NewApiError("no state was provided"))
		}
		defer csrfStorage.Delete(csrfToken)

		var sessionID string
		if cachedSessionID, ok := csrfStorage.Get(csrfToken); !ok {
			return c.JSON(http.StatusGone, NewApiError("oauth code has expired, try again"))
		} else if sessionID, ok = cachedSessionID.(string); !ok {
			return c.JSON(http.StatusGone, NewApiError("oauth code has expired, try again"))
		}
		if !isCsrfTokenMatchSession(csrfToken, sessionID, sessionSecret) {
			return c.JSON(http.StatusGone, NewApiError("CSRF attack detected"))
		}

		oauthToken, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		}

		jwtToken := jwt.New(jwt.SigningMethodHS256)
		jwtToken.Claims["iss"] = issuer
		jwtToken.Claims["iat"] = time.Now().Unix()
		jwtToken.Claims["exp"] = oauthToken.Expiry.Unix()
		jwtToken.Claims["access_token"] = oauthToken.AccessToken

		signature := getJwtSignature(jwtSecret, oauthToken.AccessToken)
		stringToken, err := jwtToken.SignedString(signature)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, NewApiError("could not sign token: "+err.Error()))
		}

		return c.JSON(http.StatusOK,
			oauthVerifyResponse{stringToken, oauthToken.Expiry.Unix()},
		)
	}
}

func getJwtSignature(jwtSecret, accessToken string) []byte {
	var buffer bytes.Buffer
	buffer.WriteString(accessToken)
	buffer.WriteString(jwtSecret)
	return buffer.Bytes()
}

func generateCsrfToken(sessionID, sessionSecret string) string {
	return lib.HMACSha1(sessionSecret, sessionID)
}

func isCsrfTokenMatchSession(csrfToken, sessionID, sessionSecret string) bool {
	return csrfToken == generateCsrfToken(sessionID, sessionSecret)
}
