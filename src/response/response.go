package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResultSet struct {
	Total int
	Items interface{}
}

type CreatedResponse struct {
	Id string
}

type NotificationsRespoonse struct {
	Result ResultSet
}

func Success(w http.ResponseWriter) {
	var bytes []byte
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func Created(w http.ResponseWriter, id string) {
	resp := CreatedResponse{id}
	bytes, _ := json.Marshal(resp)
	send(w, bytes, http.StatusCreated)
}

func Results(w http.ResponseWriter, data interface{}, size int) {
	rs := ResultSet{size, data}
	resp := NotificationsRespoonse{rs}
	bytes, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err.Error())
		var empty []byte
		send(w, empty, http.StatusOK)
		return
	}
	send(w, bytes, http.StatusOK)
}

func send(w http.ResponseWriter, bytes []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bytes)
}
