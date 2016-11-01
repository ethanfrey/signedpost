package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RenderQuery writes the output of any of the above issues to the client
func RenderQuery(rw http.ResponseWriter, res interface{}, err error) {
	var data []byte
	if err == nil {
		data, err = json.Marshal(res)
	}
	if err != nil {
		rw.WriteHeader(400)
		rw.Write([]byte(fmt.Sprintf("%+v", err)))
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}
