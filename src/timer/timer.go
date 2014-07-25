package timer

import "time"
import "fmt"
import "msg"

func Add(notice *msg.Notification) {
	timer := time.NewTimer(time.Duration(notice.Until) * time.Second)
	
	go func(notice *msg.Notification) {
		<-timer.C
		if notice.Isok() {
			notice.Send()
		}
		fmt.Println("Finished")
	}(notice)

//	timer.Stop()
}