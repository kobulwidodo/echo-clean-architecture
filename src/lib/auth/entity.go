package auth

type UserAuthInfo struct {
	User  User
	Token string
}

type User struct {
	ID       uint
	GuestId  string
	Username string
	Name     string
	IsAdmin  bool
}
