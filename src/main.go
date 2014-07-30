package main

import "fmt"
import "net/http"
import "strconv"
import "strings"
import "msg"
import "log"
import "time"
import "timer"

//import "os"
import "github.com/drone/routes"
import "github.com/fzzy/radix/redis"

var logger *log.Logger = msg.NewLog("run.log")
var rc *redis.Client = RedisClient("192.168.3.141", "6379")

func main() {
	mux := routes.New()

	mux.Get("/notifications", index)
	mux.Post("/notifications", add)
	// mux.Patch("/notifications", edit)
	mux.Del("/notifications", remove)

	http.Handle("/", mux)
	http.ListenAndServe(":8000", nil)
}

func init() {
	// todo, recover from log when start...
}

func index(w http.ResponseWriter, r *http.Request) {
	// params := r.URL.Query()
	// id := params.Get("id")
	// firstName := params.Get("first")
	// fmt.Fprintf(w, "you are %s %s", firstName, lastName)
	// noti := msg.Load(id, logger, rc)
	// fmt.Println(noti.Delay)
	// fmt.Fprintf(w, "notification: \nid:%s\ndelay:%d\nis ok:%d\nmsg:%s", noti.Id, noti.Delay, noti.Ok, noti.Msg)

	notis := msg.ReadAll(logger, rc)
	for _, noti := range notis {
		fmt.Printf("notification: \nid:%s\ndelay:%d\nis ok:%d\nmsg:%s\n", noti.Id, noti.Delay, noti.Ok, noti.Msg)
		// fmt.Println(noti.String())
	}
}

func add(w http.ResponseWriter, r *http.Request) {
	delayRaw := r.FormValue("delay")
	if delayRaw == "" {
		http.Error(w, "`delay` is empty", http.StatusInternalServerError)
		return
	}
	delay, err := strconv.Atoi(delayRaw)
	if err != nil || delay <= 0 {
		http.Error(w, "`delay` is invalid", http.StatusInternalServerError)
		return
	}
	message := r.FormValue("message")
	message = strings.Trim(message, " ")
	if message == "" {
		http.Error(w, "`message` is empty", http.StatusInternalServerError)
		return
	}

	noti := msg.New(logger)
	fmt.Println(noti.Id)
	noti.Delay = delay
	noti.Msg = message
	msg.Save(noti, rc)

	timer.Add(noti)

	fmt.Fprintf(w, "msg: %s will expire after %d", noti.Msg, noti.Delay)
}

func remove(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get("id")
	if id == "" {
		http.Error(w, "`id` is empty", http.StatusInternalServerError)
		return
	}

	timer.Stop(id)
	fmt.Fprintf(w, "stopped")
}

func RedisClient(host string, port string) *redis.Client {
	c, err := redis.DialTimeout("tcp", host+":"+port, time.Duration(10)*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	// defer c.Close()

	r := c.Cmd("select", 2)
	if r.Err != nil {
		log.Fatal(err)
	}

	return c
}
