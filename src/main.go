package main 

import "fmt"
import "net/http"
import "strconv"
import "strings"
import "msg"
import "log"
//import "os"
import "github.com/drone/routes"

var logger *log.Logger = msg.NewLog("run.log")

func main() {
	mux := routes.New()
	
    mux.Get("/notifications", index)
    mux.Post("/notifications", add)
    mux.Patch("/notifications", edit)
    mux.Del("/notifications", remove)

    http.Handle("/", mux)
    http.ListenAndServe(":8000", nil)
    
}

func init() {
    // todo, recover from log file when start...
}

//func rootHandler(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "rootHandler: %s\n", r.URL.Path)
//    fmt.Fprintf(w, "URL: %s\n", r.URL)
//    fmt.Fprintf(w, "Method: %s\n", r.Method)
//    fmt.Fprintf(w, "RequestURI: %s\n", r.RequestURI )
//    fmt.Fprintf(w, "Proto: %s\n", r.Proto)
//    fmt.Fprintf(w, "HOST: %s\n", r.Host) 
//}

func index(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
    lastName := params.Get("last")
    firstName := params.Get("first")
    fmt.Fprintf(w, "you are %s %s", firstName, lastName)
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
    noti.Send()
    
    fmt.Fprintf(w, "msg: %s will expire after %d", noti.Msg, noti.Delay)
}

func edit(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
    lastName := params.Get(":last")
    firstName := params.Get(":first")
    fmt.Fprintf(w, "you are %s %s", firstName, lastName)
}

func remove(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
    lastName := params.Get(":last")
    firstName := params.Get(":first")
    fmt.Fprintf(w, "you are %s %s", firstName, lastName)
}
