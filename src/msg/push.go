package msg

import ()

func NewIos(notice *Notification) *map[string]interface{} {
	ios := make(map[string]interface{}, 5)

	ios["alert"] = notice.Msg
	ios["sound"] = "happy"
	ios["badge"] = 1
	if notice.Link != "" {
		extra := make(map[string]string, 1)
		extra["push_link"] = notice.Link
		ios["extras"] = extra
	}

	return &ios
}

func NewAndroid(notice *Notification) *map[string]interface{} {
	ard := make(map[string]interface{}, 5)

	ard["alert"] = notice.Msg
	ard["title"] = "Fishsaying"
	ard["builder_id"] = 3
	if notice.Link != "" {
		extra := make(map[string]string, 1)
		extra["push_link"] = notice.Link
		ard["extras"] = extra
	}

	return &ard
}
