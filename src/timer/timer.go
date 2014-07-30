package timer

import "time"
import "fmt"
import "msg"
import "send"

type TimerMap struct {
	timers map[string]*time.Timer
	size   int
}

var bucket *TimerMap // timer bucket

func init() {
	bucket = &TimerMap{make(map[string]*time.Timer, 50), 0}
}

func Add(notice *msg.Notification) {
	timer := time.NewTimer(time.Duration(notice.Delay) * time.Second)

	go func(notice *msg.Notification) {
		<-timer.C
		if notice.Isok() {
			// notice.Send()
			returns, ret := send.Solo(notice)
			fmt.Println(notice.User)
			if ret == false {
				fmt.Println("Sent fails, ", returns)
			} else {
				fmt.Println("Jpush says: ", returns)
			}
		}
		remove(notice.Id)
		echoSize()
		fmt.Println("Finished")
	}(notice)

	bucket.timers[notice.Id] = timer
	bucket.size += 1

	echoSize()
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
}

func echoSize() {
	fmt.Printf("Current timers number: %v\n", bucket.size)
}
