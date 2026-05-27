package auth

type ACL struct {
	CanRead  bool
	CanWrite bool
	CanAdmin bool
}
