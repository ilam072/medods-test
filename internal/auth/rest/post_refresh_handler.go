package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"medods-test/internal/auth/service"
	"medods-test/pkg/logger"
	"net/http"
)

type refreshTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) RefreshTokensHandler(c *gin.Context) {
	var input refreshTokens
	if err := c.BindJSON(&input); err != nil {
		logger.Errorf("failed to decode request body: %s", err.Error())
		newResponse(c, http.StatusBadRequest, "invalid body request")
		return
	}

	ip := c.ClientIP()
	tokens, err := h.auth.User.RefreshTokens(c.Request.Context(), ip, input.AccessToken, input.RefreshToken)
	if err != nil {
		logger.Errorf("failed to refresh tokens: %s", err.Error())
		switch {
		case errors.Is(err, service.ErrSessionNotFound):
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, service.ErrInvalidRefreshToken):
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, service.ErrRefreshTokenAlreadyUsed):
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, service.ErrRefreshTokenExpired):
			newResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, "Something went wrong. Try again later!")
		return
	}

	c.JSON(http.StatusOK, responseToken{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
