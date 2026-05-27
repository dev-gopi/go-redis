package auth

import "os"

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

func LoadFromEnv() {
	username := os.Getenv("REDIS_USERNAME")
	if username == "" {
		username = os.Getenv("REDIS_USER")
	}

	password := os.Getenv("REDIS_PASSWORD")
	if password == "" {
		password = os.Getenv("REDIS_AUTH_PASSWORD")
	}

	if username != "" {
		Manager.Username = username
	}

	Manager.Password = password
	Manager.Enabled = password != ""
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
