package msg

import "log"
import "os"
import "fmt"

type Notification struct {
	*log.Logger
	ok 		bool
	Delay	int		// utc, seconds
	Msg		string
}

func New(logfile string) *Notification {
	flog, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
	    fmt.Println("Can not open log file", err)
	    log.Fatal(err)
	}
	
	logger := log.New(flog, "[noti]", log.Ldate|log.Ltime)
	ok := true
	delay := 0 		// delay: 0 means send it immediately
	msg := ""

	noti := &Notification{logger, ok, delay, msg}

	return noti
}

func (o *Notification) Isok() bool {
	return o.ok
}

func (o *Notification) Send() error {
    o.Println("Send msg to Jpush begin...")
	return nil
}