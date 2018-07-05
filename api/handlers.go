package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func parseLastPathSegment(r *http.Request) (int, error) {
	segs := strings.Split(r.URL.EscapedPath(), "/")
	lastSeg := segs[len(segs)-1]
	return strconv.Atoi(lastSeg)
}

func respondWithError(w http.ResponseWriter, message string) error {
	errRsp := ErrorResponse{Message: message}
	errRspJson, err := json.Marshal(errRsp)
	if err != nil {
		return err
	}

	_, err = w.Write(errRspJson)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	return nil
}
