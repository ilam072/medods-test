package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"medods-test/internal/auth/types"
	"net/http"
)

type UserService interface {
	SignUp(ctx context.Context, input types.UserDTO) error
	SingIn(ctx context.Context, input types.UserDTO, IP string) (types.Tokens, error)
	CreateSession(ctx context.Context, userId int, IP string) (types.Tokens, error)
}

type UseCase struct {
	User UserService
}
type Handler struct {
	api  *gin.Engine
	auth *UseCase
}

func New(auth *UseCase) *Handler {
	api := gin.Default()

	h := &Handler{
		api:  api,
		auth: auth,
	}

	// init endpoints
	api.POST("/auth/sign-up", h.SignUpHandler)
	api.POST("/auth/sign-in", h.SignInHandler)

	return h
}

func (h *Handler) Handler() http.Handler {
	return h.api.Handler()
}
