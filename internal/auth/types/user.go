package types

type User struct {
	// UserId   int
	UserUUID string
	Email    string
	Password string
}

type UserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
