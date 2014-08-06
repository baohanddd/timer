package send

import "jpush/api/push"
import "msg"
import "log"

const SECRET = "6cf473efa345e8aed996f39c"
const APPKEY = "2ae1ef727c1060c680ddde83"

var pf push.Platform
var ad push.Audience

func init() {
	pf.All()
}

func Push(noti *msg.Notification) {
	var (
		returns string
		err     error
	)

	if noti.IsEmptyUser() {
		ad.All()
	} else {
		ad.SetAlias(noti.User)
	}

	returns, err = send(noti)
	if err != nil {
		log.Println("[Fails]:", returns)
	} else {
		log.Println("Sent", noti.Id)
		log.Println("[Success]:", returns)
	}
}

func send(noti *msg.Notification) (returns string, err error) {
	notice := make(map[string]interface{}, 2)
	notice["android"] = msg.NewAndroid(noti)
	notice["ios"] = msg.NewIos(noti)

	nb := push.NewNoticeBuilder()
	nb.Options.Apns_production = noti.ProductMode
	// log.Fatalf("nb.Options.Apns_production = %v\n", nb.Options.Apns_production)
	nb.SetPlatform(&pf)
	nb.SetAudience(&ad)
	nb.SetNotice(notice)

	c := push.NewPushClient(SECRET, APPKEY)
	return c.Send(nb)
}
