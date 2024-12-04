package types

type User struct {
	UserId   int
	Email    string
	Password string
}

type UserDTO struct {
	UserId   int    `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignInDTO struct {
	email    string
	password string
}
