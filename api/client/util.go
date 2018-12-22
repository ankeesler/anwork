package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/ankeesler/anwork/api"
)

func encodeBody(input interface{}) (io.Reader, error) {
	var body io.Reader
	if input != nil {
		data, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(data)
	}
	return body, nil
}

func decodeBody(body io.Reader, output interface{}) error {
	if output != nil {
		bytes, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(bytes, output); err != nil {
			return fmt.Errorf("cannot unmarshal response body (%s): '%s'", err.Error(), string(bytes))
		}
	}

	return nil
}

func decodeError(body io.Reader) string {
	bodyData, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Sprintf("??? (ReadAll: %s)", err.Error())
	}

	errMsg := api.Error{}
	if err := json.Unmarshal(bodyData, &errMsg); err != nil {
		return fmt.Sprintf("??? (Unmarshal: %s)", err.Error())
	}

	return errMsg.Message
}

func is4xxStatus(rsp *http.Response) bool {
	return rsp.StatusCode >= 400 && rsp.StatusCode < 500
}

func is5xxStatus(rsp *http.Response) bool {
	return rsp.StatusCode >= 500 && rsp.StatusCode < 600
}

func parseID(location string, id *int) bool {
	segments := strings.Split(location, "/")
	idS := segments[len(segments)-1]
	idN, err := strconv.Atoi(idS)
	if err != nil {
		return false
	}

	*id = idN
	return true
}
