package handler

import (
	"net/http"
	"time"

	"github.com/pmylund/go-cache"

	"github.com/labstack/echo"
	"github.com/seesawlabs/ivan-kirichenko-exercise/lib"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

const DefaultTokenExpiration time.Duration = 5 * time.Minute

type TokenStorage interface {
	Set(k string, x interface{}, d time.Duration)
	Get(k string) (interface{}, bool)
}

func GetOauthHandler(appid, secret, redirectUrl string, csrfStorage TokenStorage) func(c *echo.Context) error {
	// oauthUrlResponse is a type which is used only within oauth handlers
	type oauthUrlResponse struct {
		Url     string `json:"url"`
		Message string `json:"message"`
	}

	return func(c *echo.Context) error {
		conf := &oauth2.Config{
			ClientID:     appid,
			ClientSecret: secret,
			RedirectURL:  redirectUrl,
			Scopes:       []string{},
			Endpoint:     facebook.Endpoint,
		}
		token, err := lib.GenerateRandomString(32)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		}

		url := conf.AuthCodeURL(token, oauth2.AccessTypeOffline)

		csrfStorage.Set(token, true, cache.DefaultExpiration)

		return c.JSON(http.StatusOK, oauthUrlResponse{url, "please use this url to authenticate with facebook"})
	}
}

func GetOauthVerifyHandler(appid, secret, redirectUrl string, csrfStorage, tokenStorage TokenStorage) func(c *echo.Context) error {
	type oauthVerifyResponse struct {
		AccessToken string `json:"access_token"`
		Expires     uint64 `json:"expires"`
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
		if _, ok := csrfStorage.Get(csrfToken); !ok {
			return c.JSON(http.StatusGone, NewApiError("oauth code has expired, try again"))
		}

		conf := &oauth2.Config{
			ClientID:     appid,
			ClientSecret: secret,
			RedirectURL:  redirectUrl,
			Scopes:       []string{},
			Endpoint:     facebook.Endpoint,
		}

		token, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		}

		var expires time.Duration = DefaultTokenExpiration
		if !token.Expiry.IsZero() {
			expires = token.Expiry.Sub(time.Now())
		}

		tokenStorage.Set(token.AccessToken, true, expires)

		return c.JSON(http.StatusOK,
			oauthVerifyResponse{token.AccessToken, uint64(expires.Seconds())},
		)
	}
}
