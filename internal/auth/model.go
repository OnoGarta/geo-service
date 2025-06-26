package auth

type User struct {
	Username string
	Password []byte // hash bcrypt
}
