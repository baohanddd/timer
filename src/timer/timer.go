package timer

import "time"
import "fmt"
import "msg"

type TimerMap struct {
	timers		map[string]*time.Timer
	size		int
}

func Add(notice *msg.Notification) {
	timer := time.NewTimer(time.Duration(notice.Delay) * time.Second)
	
	go func(notice *msg.Notification) {
		<-timer.C
		if notice.Isok() {
			notice.Send()
		}
		fmt.Println("Finished")
	}(notice)

//	timer.Stop()
}