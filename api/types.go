package api

// Error is sent in any non 2XX responses.
type Error struct {
	Message string
}
