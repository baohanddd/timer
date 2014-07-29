package msg

import "log"
import "os"
import "fmt"
import "crypto/rand"
import "github.com/fzzy/radix/redis"
import "bytes"
import "encoding/gob"

type Notification struct {
	*log.Logger
	Ok 		bool
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
	return o.Ok
}

func (o *Notification) Send() error {
    o.Println("Send msg to Jpush begin...")
	return nil
}

func Save(noti *Notification, client *redis.Client) {
	var buf bytes.Buffer        // Stand-in for a buf connection
	enc := gob.NewEncoder(&buf) // Will write to buf.
	

	err := enc.Encode(noti.Ok)
   	if err != nil {
        log.Fatal("encode error:", err)
    }
    err = enc.Encode(noti.Id)
   	if err != nil {
        log.Fatal("encode error:", err)
    }
    err = enc.Encode(noti.Delay)
   	if err != nil {
        log.Fatal("encode error:", err)
    }
    err = enc.Encode(noti.Msg)
   	if err != nil {
        log.Fatal("encode error:", err)
    }
    fmt.Println(buf)

    r := client.Cmd("hset", "notification", noti.Id, buf.Bytes())
    if r.Err != nil {
    	log.Fatal("save notification fails ", r.Err)
    } 
}

func Load(id string, l *log.Logger, client *redis.Client) *Notification {
	var ao *Notification = New(l)

	data, err := client.Cmd("hget", "notification", id).Bytes()
	if err != nil {
		log.Fatal("Get notification fails ", err)
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	err = dec.Decode(&ao.Ok)
	if err != nil {
	    log.Fatal("decode error 1:", err)
	}    
	err = dec.Decode(&ao.Id)
	if err != nil {
	    log.Fatal("decode error 1:", err)
	}    
	err = dec.Decode(&ao.Delay)
	if err != nil {
	    log.Fatal("decode error 1:", err)
	}    
	err = dec.Decode(&ao.Msg)
	if err != nil {
	    log.Fatal("decode error 1:", err)
	}
	fmt.Println(ao)

	return ao
}

func uuid() string {
     b := make([]byte, 16)
     rand.Read(b)
     b[6] = (b[6] & 0x0f) | 0x40
     b[8] = (b[8] & 0x3f) | 0x80
     return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}