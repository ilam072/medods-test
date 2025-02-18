package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"medods-test/internal/auth/service"
	"medods-test/internal/auth/types"
	"medods-test/pkg/logger"
	"net/http"
)

type userSignUp struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

func (h *Handler) SignUpHandler(c *gin.Context) {
	var input userSignUp
	if err := c.BindJSON(&input); err != nil {
		logger.Errorf("failed to decode request body: %s", err.Error())
		newResponse(c, http.StatusBadRequest, "invalid body request")
		return
	}

	if err := h.auth.User.SignUp(c.Request.Context(), types.UserDTO{
		Email:    input.Email,
		Password: input.Password,
	}); err != nil {
		logger.Errorf("failed to sign up: %s", err.Error())
		if errors.Is(err, service.ErrUserAlreadyExists) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, "Something went wrong. Try again later!")
		return
	}

	c.Status(http.StatusCreated)
}
