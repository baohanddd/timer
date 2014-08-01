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
	User     []string
	Msg      string
	Link     string
}

const KEY = "notification"

var rc *redis.Client = common.RedisClient("192.168.3.141", "6379")

func NewForm(data url.Values) (*Notification, error) {
	var (
		err    error
		delay  int
		users  []string
		notice Notification
	)

	users = make([]string, 1)

	uid := data.Get("user_id")
	if uid != "" {
		users[0] = uid
	}

	raw := data.Get("delay")
	if raw != "" {
		delay, err = strconv.Atoi(raw)
		if err != nil || delay < 0 {
			return nil, errors.New("`delay` is invalid")
		}
	}

	msg := data.Get("message")
	if msg == "" {
		return nil, errors.New("`message` is invalid")
	}

	notice.Delay = delay
	notice.Id = uuid()
	notice.Msg = msg
	notice.User = users
	notice.SendTime = time.Now().Unix() + int64(delay)
	notice.Link = data.Get("link")

	return &notice, nil
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
	delay := int(o.SendTime - Now)
	if delay > 0 {
		o.Delay = delay
	} else {
		o.Delay = 0
	}
	return true
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
	var noti Notification

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(&noti)
	if err != nil {
		log.Fatal("decode error 1:", err)
	}

	return &noti
}

func uuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
