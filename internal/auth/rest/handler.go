package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"medods-test/internal/auth/types"
)

type UserService interface {
	SignUp(ctx context.Context, input types.UserDTO) error
	SingIn(ctx context.Context, input types.UserSignInDTO, IP string) (types.Tokens, error)
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

	// init endpoints

	return &Handler{
		api:  api,
		auth: auth,
	}
}
