package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ishu17077/code_runner_backend/helpers"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("token")
		clientToken := cookie.Value
		if err != nil {
			clientToken = c.Request.Header.Get("token")
			if clientToken == "" {
				// c.JSON(http.StatusUnauthorized, gin.H{"error": "No authentication token or cookie provided"})
				// c.Abort()
				c.Next()
				return
			}
			if cookie.HttpOnly != true {
				c.JSON(http.StatusNotAcceptable, gin.H{"error": "Cookie not valid"})
				c.Abort()
				return
			}

		}
		claims, tokErr := helpers.ValidateToken(clientToken)
		if tokErr != nil && claims != nil && claims.Username != "" {
			refreshCookie, err := c.Request.Cookie("refresh_token")
			refreshToken := refreshCookie.Value
			if err != nil || refreshToken == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Cookie not valid"})
				ClearCookie(c, "token", "refresh_token")
				c.Abort()
				return
			}
			if refreshCookie.HttpOnly != true {
				c.JSON(http.StatusNotAcceptable, gin.H{"error": "Cookie not valid"})
				ClearCookie(c, "token", "refresh_token")
				c.Abort()
				return
			}

			if errors.Is(tokErr, jwt.ErrTokenExpired) {
				refreshClaims, refreshTokErr := helpers.ValidateToken(refreshToken)
				if refreshTokErr != nil || refreshClaims == nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Tokens expired"})
					return
				}
				signedToken, signedRefreshToken, err := helpers.GenerateTokens(claims.Username)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Sign in again!"})
					return
				}
				SetCookie(c, "token", signedToken)
				SetCookie(c, "refresh_token", signedRefreshToken)
				c.Set("username", claims.Username)
				c.Set("admin", true)
				c.Next()
				return
			}
			ClearCookie(c, "token", "refresh_token")
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("admin", true)
		c.Next()
	}
}

func SetCookie(c *gin.Context, name, value string) {
	//TODO: Set domain as cors hosts && set secure for production environments
	c.SetCookie(
		name,
		value,
		2160000,
		"/",
		"",
		false,
		true,
	)
}

func ClearCookie(c *gin.Context, names ...string) {
	for _, name := range names {
		c.SetCookie(name, "", -1, "/", "", false, true)
	}
}
