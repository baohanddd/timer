package msg

import "log"

type Notification struct {
	*log.Logger
	ok		bool
	Until	int		// utc, seconds
}

func (o *Notification) Isok() bool {
	return o.ok
}

func (o *Notification) Send() error {
    o.Println("Send msg to Jpush begin...")
	return nil
}
