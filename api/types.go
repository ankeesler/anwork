package api

// Error is sent in any non 2XX responses.
type Error struct {
	Message string
}

// Auth is sent in response to a successful POST to the /api/v1/auth endpoint.
type Auth struct {
	Token string
}
