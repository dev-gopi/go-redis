package auth

type AuthManager struct {
	Enabled  bool
	Username string
	Password string
}

var Manager = &AuthManager{
	Enabled:  false,
	Username: "default",
	Password: "",
}

func (a *AuthManager) Authenticate(
	username string,
	password string,
) bool {

	if !a.Enabled {
		return true
	}

	return a.Username == username &&
		a.Password == password
}
