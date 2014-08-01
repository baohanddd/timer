package timer

import "time"
import "log"
import "msg"
import "send"

type TimerMap struct {
	timers map[string]*time.Timer
	size   int
}

var bucket *TimerMap // timer bucket

func init() {
	bucket = &TimerMap{make(map[string]*time.Timer, 1024), 0}
}

func Add(notice *msg.Notification) {
	timer := time.NewTimer(time.Duration(notice.Delay) * time.Second)

	go func(notice *msg.Notification) {
		<-timer.C
		send.Push(notice)
		notice.Delete()
		remove(notice.Id)
		EchoSize()
	}(notice)

	bucket.timers[notice.Id] = timer
	bucket.size += 1

	log.Println("New item arrival:", notice.Id)
	EchoSize()
}

func Size() int {
	return bucket.size
}

func Stop(id string) bool {
	timer, ok := bucket.timers[id]
	if ok {
		timer.Stop()
		remove(id)
		return true
	}
	return false
}

func remove(id string) {
	delete(bucket.timers, id)
	bucket.size -= 1
	msg.Delete(id)
}

func EchoSize() {
	log.Printf("Current timers number: %v\n", bucket.size)
}
