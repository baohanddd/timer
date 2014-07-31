package main

// import "fmt"
import "net/http"
import "msg"
import "timer"
import "response"

// import "encoding/json"

//import "os"

import "github.com/drone/routes"

// var logger *log.Logger = msg.NewLog("run.log")

func main() {
	mux := routes.New()

	mux.Get("/notifications", index)
	mux.Post("/notifications", add)
	mux.Del("/notifications", remove)

	http.Handle("/", mux)
	http.ListenAndServe(":8000", nil)
}

func init() {
	// todo, recover from log when start...
}

func index(w http.ResponseWriter, r *http.Request) {
	notis := msg.LoadAll()
	response.Results(w, notis, len(notis))
}

func add(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "http body is invalid", http.StatusBadRequest)
		return
	}

	noti, ferr := msg.NewForm(r.Form)

	if ferr != nil {
		http.Error(w, ferr.Error(), http.StatusBadRequest)
		return
	}

	noti.Save()
	timer.Add(noti)
	response.Created(w, noti.Id)
}

func remove(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get("id")
	if id == "" {
		http.Error(w, "`id` is empty", http.StatusInternalServerError)
		return
	}

	timer.Stop(id)
	response.Success(w)
}
