package main

import (
	"common"
	"flag"
	"fmt"
	"github.com/drone/routes"
	"log"
	"msg"
	"net/http"
	"os"
	"response"
	"runtime"
	"time"
	"timer"
)

var redis *string = flag.String("redis", "127.0.0.1:6379", "redis host and post")
var enable *bool = flag.Bool("enable", false, "determine whether enable mutlti cores support or not")
var mode *string = flag.String("mode", "stage", "`stage` or `live`")
var help *bool = flag.Bool("help", false, "Display help info")

func main() {
	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	if *enable {
		cores := runtime.NumCPU()
		runtime.GOMAXPROCS(cores)
		log.Println(cores, "cores are enabled...")
	}

	msg.SetMode(*mode)

	msg.RC = common.RedisNew(*redis)

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
