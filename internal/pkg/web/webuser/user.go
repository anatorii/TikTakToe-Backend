package webuser

type User struct {
	Login    string
	Password string
}

func NewWebUser(login string, password string) User {
	return User{
		Login:    login,
		Password: password,
	}
}
