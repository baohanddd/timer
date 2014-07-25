package main 

import "fmt"
import "net/http"
import "strconv"
import "msg"
import "log"
import "github.com/drone/routes"

var logfile string = "run.log"

func main() {
	mux := routes.New()
	
    mux.Get("/notifications", index)
    mux.Post("/notifications", add)
    mux.Patch("/notifications", edit)
    mux.Del("/notifications", remove)

    http.Handle("/", mux)
    http.ListenAndServe(":8000", nil)
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
    lastName := params.Get(":last")
    firstName := params.Get(":first")
    fmt.Fprintf(w, "you are %s %s", firstName, lastName)
}

func add(w http.ResponseWriter, r *http.Request) {
    untilRaw := r.FormValue("until")
    if untilRaw == "" {
    	fmt.Println("until is empty")
    	http.Error(w, "`until` is empty", http.StatusInternalServerError)
    }
    until, err := strconv.Atoi(untilRaw)
    if err != nil {
    	http.Error(w, "`until` is invalid", http.StatusInternalServerError)
    }
    
    message := r.FormValue("message")
    if message == "" {
    	fmt.Println("message is empty")
    	http.Error(w, "`message` is empty", http.StatusInternalServerError)
    }
    
    
    noti := &msg.Notification{&log.New()}
    
    fmt.Fprintf(w, "you are %s %s", first, last)
    
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