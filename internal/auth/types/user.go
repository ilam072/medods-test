package types

type User struct {
	UserUUID string
	Email    string
	Password string
}

type UserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
