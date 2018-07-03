// This package contains the HTTP handlers used by the anwork API.
//
// See api package for API definition.
package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

func parseLastPathSegment(r *http.Request) (int, error) {
	segs := strings.Split(r.URL.EscapedPath(), "/")
	lastSeg := segs[len(segs)-1]
	return strconv.Atoi(lastSeg)
}
