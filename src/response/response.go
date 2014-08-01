package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResultSet struct {
	Total int
	Items interface{}
}

type CreatedResponse struct {
	Id string
}

func Success(w http.ResponseWriter) {
	var bytes []byte
	send(w, bytes, http.StatusOK)
}

func Created(w http.ResponseWriter, id string) {
	resp := CreatedResponse{id}
	bytes, _ := json.Marshal(resp)
	send(w, bytes, http.StatusCreated)
}

func Results(w http.ResponseWriter, data interface{}, size int) {
	resp := ResultSet{size, data}
	bytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("JSON encoding fails:", err.Error())
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
