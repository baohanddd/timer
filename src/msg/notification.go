package msg

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"menteslibres.net/gosexy/redis"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var RC *redis.Client
var mode int

func SetMode(m string) {
	m = strings.ToLower(m)
	switch m {
	case "stage":
		mode = 1
	case "live":
		mode = 2
	default:
		mode = 1
	}
}

type Notification struct {
	Id          string // uuid, uniqueness
	Delay       int    // utc, seconds
	SendTime    int64  // a timestamp, seconds
	User        []string
	Msg         string
	Link        string
	ProductMode bool
}

const KEY = "notification"

func NewForm(data url.Values) (*Notification, error) {
	var (
		err    error
		id     string
		delay  int
		users  []string
		notice Notification
	)

	id = data.Get("id")
	if id == "" {
		id = uuid()
	}

	uid := data.Get("user_id")
	if uid != "" {
		users = strings.Split(strings.Trim(uid, " "), ",")
		for i, user := range users {
			users[i] = strings.Trim(user, " ")
		}
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
	notice.Id = id
	notice.Msg = msg
	notice.User = users
	notice.SendTime = time.Now().Unix() + int64(delay)
	notice.Link = data.Get("link")
	if mode == 1 {
		notice.ProductMode = false
	} else if mode == 2 {
		notice.ProductMode = true
	}

	// log.Printf("notice.ProductMode = %v\n", notice.ProductMode)
	// log.Printf("mode = %v\n", mode)

	return &notice, nil
}

func LoadOne(id string) (*Notification, error) {
	data, err := RC.HGet(KEY, id)
	if err != nil {
		return nil, err
	}
	o := decode([]byte(data))
	fmt.Printf("%+v\n", o)
	return o, nil
}

func LoadAll() []*Notification {
	rows, err := RC.HVals(KEY)
	if err != nil {
		log.Println("Can not read notification: ", err)
	}
	var size int = len(rows)
	ret := make([]*Notification, size)
	for i, data := range rows {
		ret[i] = decode([]byte(data))
	}
	return ret
}

func Delete(id string) bool {
	_, err := RC.HDel(KEY, id)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (o *Notification) IsEmptyUser() bool {
	return len(o.User) == 0
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

func (o *Notification) String() string {
	return fmt.Sprintf("notification: \nid:%s\ndelay:%d\nmsg:%s\nuser:%s\nsend_time:%v\n",
		o.Id, o.Delay, o.Msg, o.User, o.SendTime)
}

func (o *Notification) Save() {
	var (
		err error
	)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(o)
	if err != nil {
		log.Println("Encode notification fails:", err)
	}

	_, err = RC.HSet(KEY, o.Id, buf.Bytes())
	if err != nil {
		log.Println("Save notification fails:", err)
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
