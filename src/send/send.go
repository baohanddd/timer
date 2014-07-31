package send

import "jpush/api/push"
import "msg"
import "fmt"

const SECRET = "6cf473efa345e8aed996f39c"
const APPKEY = "2ae1ef727c1060c680ddde83"

var pf push.Platform
var ad push.Audience
var content push.Notice

func init() {
	pf.All()
}

func Solo(noti *msg.Notification) (string, bool) {
	users := []string{noti.User}
	ad.SetAlias(users)
	content.Alert = noti.Msg

	nb := push.NewNoticeBuilder()
	nb.Options.Apns_production = false
	nb.SetPlatform(&pf)
	nb.SetAudience(&ad)
	nb.SetSimpleNotice(noti.Msg)

	// 	nb := `{
	//    "platform": "all",
	//    "audience" : {
	//       "alias" : ["stg_52a2e9149a0b8aea75956b88"]
	//    },
	//    "notification" : {
	//          "alert" : "Hi, JPush!"
	//    }
	// }`

	fmt.Println(nb)
	c := push.NewPushClient(SECRET, APPKEY)
	str, err := c.Send(nb)

	if err != nil {
		fmt.Println(err)
		return str, false
	}
	return str, true
}
