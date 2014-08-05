package msg

import "log"
import "fmt"
import "crypto/rand"
import "bytes"
import "time"
import "errors"
import "net/url"
import "encoding/gob"
import "strconv"
import "strings"
import "common"

import "menteslibres.net/gosexy/redis"

const HOST = "127.0.0.1"
const PORT = 6379

var rc *redis.Client = common.RedisNew(HOST, uint(PORT))

type Notification struct {
	Id       string // uuid, uniqueness
	Delay    int    // utc, seconds
	SendTime int64  // a timestamp, seconds
	User     []string
	Msg      string
	Link     string
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

	return &notice, nil
}

func LoadOne(id string) (*Notification, error) {
	data, err := rc.HGet(KEY, id)
	if err != nil {
		return nil, err
	}
	o := decode([]byte(data))
	fmt.Printf("%+v\n", o)
	return o, nil
}

func LoadAll() []*Notification {
	rows, err := rc.HVals(KEY)
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
	_, err := rc.HDel(KEY, id)
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

	_, err = rc.HSet(KEY, o.Id, buf.Bytes())
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
