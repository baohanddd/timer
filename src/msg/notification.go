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
	Id       string // uuid, uniqueness
	Delay    int    // utc, seconds
	SendTime int64  // a timestamp, seconds
	User     string
	Msg      string
}

const KEY = "notification"

var rc *redis.Client = common.RedisClient("192.168.3.141", "6379")

func New() *Notification {
	uuid := uuid()
	delay := 0 // delay: 0 means send it immediately
	userId := ""
	st := time.Now().Unix() // send time
	msg := ""

	noti := &Notification{uuid, delay, st, userId, msg}

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
	rows, err := rc.Cmd("hvals", KEY).ListBytes()
	if err != nil {
		log.Println("Can not read notification: ", err)
	}
	var size int = len(rows)
	ret := make([]*Notification, size)
	for i, data := range rows {
		ret[i] = decode(data)
	}
	return ret
}

func Delete(id string) bool {
	r := rc.Cmd("hdel", KEY, id)
	if r.Err != nil {
		log.Println(r.Err)
		return false
	}
	return true
}

func (o *Notification) ReBuild(Now int64) bool {
	if o.SendTime >= Now {
		o.Delay = int(o.SendTime - Now)
		return true
	}
	return false
}

func (o *Notification) Delete() bool {
	return Delete(o.Id)
}

// func (o *Notification) Send() error {
// 	o.Println("Send msg to Jpush begin...")
// 	send.Solo(o)
// 	return nil
// }

func (o *Notification) String() string {
	return fmt.Sprintf("notification: \nid:%s\ndelay:%d\nmsg:%s\nuser:%s\nsend_time:%v\n",
		o.Id, o.Delay, o.Msg, o.User, o.SendTime)
}

func (o *Notification) Save() {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(o)
	if err != nil {
		log.Println("Encode notification fails:", err)
	}

	r := rc.Cmd("hset", "notification", o.Id, buf.Bytes())
	if r.Err != nil {
		log.Println("Save notification fails:", r.Err)
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
