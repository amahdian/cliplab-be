package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/amahdian/cliplab-be/domain/contracts/resp"
	"github.com/gin-gonic/gin"
)

func VerifyRecaptcha(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if secret == "" {
			c.Next()
			return
		}

		token := c.GetHeader("X-Recaptcha-Token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, resp.NewErrorResponse(fmt.Errorf("recaptcha token is required")))
			return
		}

		verifyUrl := "https://www.google.com/recaptcha/api/siteverify"
		data := url.Values{}
		data.Set("secret", secret)
		data.Set("response", token)

		res, err := http.PostForm(verifyUrl, data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp.NewErrorResponse(err))
			return
		}
		defer res.Body.Close()

		var result struct {
			Success     bool      `json:"success"`
			Score       float64   `json:"score"`
			Action      string    `json:"action"`
			ChallengeTS time.Time `json:"challenge_ts"`
			Hostname    string    `json:"hostname"`
			ErrorCodes  []string  `json:"error-codes"`
		}

		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp.NewErrorResponse(err))
			return
		}

		if !result.Success {
			c.AbortWithStatusJSON(http.StatusForbidden, resp.NewErrorResponse(fmt.Errorf("recaptcha verification failed: %v", result.ErrorCodes)))
			return
		}

		if result.Score < 0.6 {
			c.AbortWithStatusJSON(http.StatusForbidden, resp.NewErrorResponse(fmt.Errorf("recaptcha score too low: %f", result.Score)))
			return
		}

		c.Next()
	}
}
