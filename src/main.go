package main

import "net/http"
import "log"
import "msg"
import "timer"
import "time"
import "response"
import "fmt"
import "flag"
import "common"

// import "runtime"
import "github.com/drone/routes"

var RedisHost *string = flag.String("rh", "", "redis host, default value: 127.0.0.1")
var RedisPort *int = flag.Int("rp", 0, "redis port, default value: 6379")

func main() {
	// cores := runtime.NumCPU()
	// runtime.GOMAXPROCS(cores)
	// log.Println(cores, "cores are enabled...")

	flag.Parse()

	if *RedisHost == "" || *RedisPort == 0 {
		fmt.Println("Usage: ./main -rh [:host] -rp [:port]")
		fmt.Println("Examples: ./main -rh 127.0.0.1 -rp 6379")
		return
	}

	msg.RC = common.RedisNew(*RedisHost, uint(*RedisPort))

	recovery()

	mux := routes.New()

	mux.Get("/notifications", index)
	mux.Post("/notifications", add)
	mux.Del("/notifications", remove)
	mux.Get("/hello", hello)

	log.Println("Listen 8000...")

	http.Handle("/", mux)
	http.ListenAndServe(":8000", nil)
}

func recovery() {
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

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", "hello world")
}
