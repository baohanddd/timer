package send

import "jpush/api/push"
import "msg"
import "log"

const SECRET = "6cf473efa345e8aed996f39c"
const APPKEY = "2ae1ef727c1060c680ddde83"

var pf push.Platform
var ad push.Audience
var ios push.IosNotice
var ard push.AndroidNotice

func init() {
	pf.All()
}

func Push(noti *msg.Notification) {
	var (
		returns string
		err     error
	)

	users := noti.User
	ad.SetAlias(users)

	returns, err = send2Android(noti)
	if err != nil {
		log.Println("[Sent fails]:", returns)
	} else {
		log.Println("Sent", noti.Id)
		log.Println("[Success Android]:", returns)
	}

	returns, err = send2Ios(noti)
	if err != nil {
		log.Println("[Sent fails]:", returns)
	} else {
		log.Println("Sent", noti.Id)
		log.Println("[Success Ios]:", returns)
	}
}

func send2Android(noti *msg.Notification) (returns string, err error) {
	ard := msg.NewAndroid(noti)

	nb := push.NewNoticeBuilder()
	nb.Options.Apns_production = false
	nb.SetPlatform(&pf)
	nb.SetAudience(&ad)
	nb.SetNotice(ard)

	c := push.NewPushClient(SECRET, APPKEY)
	return c.Send(nb)
}

func send2Ios(noti *msg.Notification) (returns string, err error) {
	ios := msg.NewIos(noti)

	nb := push.NewNoticeBuilder()
	nb.Options.Apns_production = false
	nb.SetPlatform(&pf)
	nb.SetAudience(&ad)
	nb.SetNotice(ios)

	c := push.NewPushClient(SECRET, APPKEY)
	return c.Send(nb)
}
