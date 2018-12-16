package client

import "fmt"

type badResponseError struct {
	code    string
	message string
}

func (bre *badResponseError) Error() string {
	return fmt.Sprintf("unexpected response: %s: %s", bre.code, bre.message)
}
