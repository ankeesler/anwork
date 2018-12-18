package authenticator

// NullAuthenticator is the yes-woman or yes-man of authenticator's - it
// says that every token is valid. It should only be used for testing.
type NullAuthenticator struct{}

func (na NullAuthenticator) Authenticate(token string) error { return nil }
func (na NullAuthenticator) Token() (string, error)          { return "", nil }
