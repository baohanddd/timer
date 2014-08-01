package main

import "net/http"
import "log"
import "msg"
import "timer"
import "time"
import "response"
import "github.com/drone/routes"

func main() {
	mux := routes.New()

	mux.Get("/notifications", index)
	mux.Post("/notifications", add)
	mux.Del("/notifications", remove)

	log.Println("Listen 8000...")

	http.Handle("/", mux)
	http.ListenAndServe(":8000", nil)
}

func init() {
	var c int = 0 // count number

	log.Println("Checking presistent items...")
	log.Println("Initialing persistent items...")

	items := msg.LoadAll()
	if len(items) == 0 {
		log.Println("There is nothing need to initialize...")
		return
	}
	now := time.Now().Unix()
	for _, item := range items {
		if item.ReBuild(now) {
			timer.Add(item)
			c += 1
		}
	}
	log.Println("Initialized ", c, "items...")
}

func index(w http.ResponseWriter, r *http.Request) {
	notis := msg.LoadAll()
	response.Results(w, notis, len(notis))
}

func add(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	noti, err := msg.NewForm(r.Form)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	log.Println(id, "cancalled.")
	timer.EchoSize()
	response.Success(w)
}
