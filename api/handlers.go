package api

import (
	"encoding/json"
	"log"
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

func respondWithError2(w http.ResponseWriter, code int, message string, log *log.Logger) {
	errRsp := ErrorResponse{Message: message}
	errRspJson, err := json.Marshal(errRsp)
	if err != nil {
		log.Printf("Unable to marshal error response: %s", err.Error())
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(errRspJson)
	if err != nil {
		log.Printf("Unable to write error response: %s", err.Error())
		return
	}

	log.Printf("Responding with error: %s", message)
}
