package api

import (
	"log"
	"net/http"
)

const homepage = `
<h1>ANWORK</h1>

<p>
Thank you for visiting ANWORK! There is currently no web frontend
for this app. Visit <a href="/api">/api</a> for information about
the API.
</p>
`

type frontPageHandler struct {
	log *log.Logger
}

func newFrontPageHandler(log *log.Logger) http.Handler {
	return &frontPageHandler{log: log}
}

func (fph *frontPageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fph.log.Printf("Handling %s /...", req.Method)

	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if _, err := w.Write([]byte(homepage)); err != nil {
		fph.log.Printf("Failed to write homepage: %s", err.Error())
	}
}
