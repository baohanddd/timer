package msg

import "log"
import "fmt"
import "crypto/rand"
import "github.com/fzzy/radix/redis"
import "bytes"
import "time"
import "errors"
import "net/url"
import "encoding/gob"
import "common"
import "strconv"

// import "send"

type Notification struct {
	Ok       bool
	Id       string // uuid, uniqueness
	Delay    int    // utc, seconds
	SendTime int64  // a timestamp, seconds
	User     string
	Msg      string
}

const KEY = "notification"

var rc *redis.Client = common.RedisClient("192.168.3.141", "6379")

// func NewLog(logfile string) *log.Logger {
// 	flog, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
// 	if err != nil {
// 		fmt.Println("Can not open log file", err)
// 		log.Fatal(err)
// 	}

// 	return log.New(flog, "[noti]", log.Ldate|log.Ltime|log.Lshortfile)
// }

func New() *Notification {
	ok := true
	uuid := uuid()
	delay := 0 // delay: 0 means send it immediately
	userId := ""
	st := time.Now().Unix() // send time
	msg := ""

	noti := &Notification{ok, uuid, delay, st, userId, msg}

	return noti
}

func NewForm(data url.Values) (*Notification, error) {
	var (
		err   error
		delay int
	)

	uid := data.Get("user_id")
	if uid == "" {
		return nil, errors.New("`user_id` is invalid")
	}

	raw := data.Get("delay")
	if raw == "" {
		return nil, errors.New("`delay` is empty")
	}
	delay, err = strconv.Atoi(raw)
	if err != nil || delay < 0 {
		return nil, errors.New("`delay` is invalid")
	}

	msg := data.Get("message")
	if msg == "" {
		return nil, errors.New("`message` is invalid")
	}

	o := &Notification{
		Ok:       true,
		Id:       uuid(),
		Delay:    delay,
		User:     uid,
		Msg:      msg,
		SendTime: time.Now().Unix() + int64(delay),
	}

	return o, nil
}

func LoadOne(id string) (*Notification, error) {
	data, err := rc.Cmd("hget", "notification", id).Bytes()
	if err != nil {
		return nil, err
	}
	o := decode(data)
	fmt.Printf("%+v\n", o)
	return o, nil
}

func LoadAll() []*Notification {
	rows, err := rc.Cmd("hgetall", KEY).ListBytes()
	if err != nil {
		fmt.Println(err)
	}
	var size int = len(rows) / 2
	ret := make([]*Notification, size)
	c := 0
	for i, data := range rows {
		if i%2 == 1 {
			ret[c] = decode(data)
			c += 1
		}
	}
	return ret
}

func Delete(id string) bool {
	r := rc.Cmd("hdel", KEY, id)
	if r.Err != nil {
		fmt.Println(r.Err)
		return false
	}
	return true
}

func (o *Notification) Isok() bool {
	return o.Ok
}

// func (o *Notification) Send() error {
// 	o.Println("Send msg to Jpush begin...")
// 	send.Solo(o)
// 	return nil
// }

func (o *Notification) String() string {
	return fmt.Sprintf("notification: \nid:%s\ndelay:%d\nok:%v\nmsg:%s\nuser:%s\nsend_time:%v\n",
		o.Id, o.Delay, o.Ok, o.Msg, o.User, o.SendTime)
}

func (o *Notification) Save() {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(o)
	if err != nil {
		log.Fatal("encode error: ", err)
	}

	r := rc.Cmd("hset", "notification", o.Id, buf.Bytes())
	if r.Err != nil {
		log.Fatal("save notification fails ", r.Err)
	}
}

func decode(data []byte) *Notification {
	noti := New()

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(&noti)
	if err != nil {
		log.Fatal("decode error 1:", err)
	}

	return noti
}

func uuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
