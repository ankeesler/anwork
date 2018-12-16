package api

import (
	"net/http"

	"github.com/ankeesler/anwork/task2"
)

// returns statusCode, response body, any failure
type handler func(*handlerData, task2.Repo) (int, interface{}, error)

type handlerData struct {
	// for matching (and documentation)
	method string
	query  string // empty if there is not query to be handled

	// for performing the handling
	handler
	body interface{}

	// for documentation
	inputType  interface{}
	outputType interface{}
}

// map from path -> handlerData
var hData = map[string][]handlerData{
	"/api/v1/tasks": []handlerData{
		{
			method: http.MethodGet,

			handler: tasksGetHandler,

			inputType:  0,
			outputType: []*task2.Task{},
		},
	},
}

func findHData(path string) []handlerData {
	hData, ok := hData[path]
	if ok {
		return hData
	} else {
		return nil
	}
}

func findHDatum(hData []handlerData, method, query string) *handlerData {
	for i := range hData {
		hDatum := &hData[i]
		if hDatum.method == method && hDatum.query == query {
			return hDatum
		}
	}
	return nil
}
