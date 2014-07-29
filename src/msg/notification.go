package msg

import "log"
import "os"
import "fmt"
import "crypto/rand"

type Notification struct {
	*log.Logger
	ok 		bool
	Id 		string	// uuid, uniqueness
	Delay	int		// utc, seconds
	Msg		string
}

func NewLog(logfile string) *log.Logger {
	flog, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
	    fmt.Println("Can not open log file", err)
	    log.Fatal(err)
	}
	
	return log.New(flog, "[noti]", log.Ldate|log.Ltime)
}

func New(l *log.Logger) *Notification {
	ok := true
	uuid := uuid()
	delay := 0 		// delay: 0 means send it immediately
	msg := ""

	noti := &Notification{l, ok, uuid, delay, msg}

	return noti
}

func (o *Notification) Isok() bool {
	return o.ok
}

func (o *Notification) Send() error {
    o.Println("Send msg to Jpush begin...")
	return nil
}

func uuid() string {
     b := make([]byte, 16)
     rand.Read(b)
     b[6] = (b[6] & 0x0f) | 0x40
     b[8] = (b[8] & 0x3f) | 0x80
     return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}