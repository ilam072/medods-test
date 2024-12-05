package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"medods-test/internal/auth/service"
	"medods-test/internal/auth/types"
	"net/http"
)

type userSignIn struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

type responseToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) SignInHandler(c *gin.Context) {
	var input userSignIn
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid body request")
		return
	}

	ip := c.ClientIP()
	tokens, err := h.auth.User.SingIn(c.Request.Context(), types.UserDTO{
		Email:    input.Email,
		Password: input.Password,
	}, ip)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, responseToken{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
