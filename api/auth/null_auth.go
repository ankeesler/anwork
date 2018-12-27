package auth

// NullAuth is the yes-woman or yes-man of authenticator's - it says that
// every token is valid. It should only be used for testing.
type NullAuth struct{}

func (na NullAuth) Authenticate(token string) error { return nil }
func (na NullAuth) Token() (string, error)          { return "", nil }
