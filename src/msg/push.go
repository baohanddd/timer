package msg

import ()

func NewIos(notice *Notification) *map[string]interface{} {
	wrap := make(map[string]interface{}, 1)
	ios := make(map[string]interface{}, 5)

	ios["alert"] = notice.Msg
	ios["sound"] = "happy"
	ios["badge"] = 5
	if notice.Link != "" {
		extra := make(map[string]string, 1)
		extra["push_link"] = notice.Link
		ios["extra"] = extra
	}

	wrap["ios"] = ios

	return &wrap
}

func NewAndroid(notice *Notification) *map[string]interface{} {
	wrap := make(map[string]interface{}, 1)
	ios := make(map[string]interface{}, 5)

	ios["alert"] = notice.Msg
	ios["title"] = "Fishsaying"
	ios["builder_id"] = 3
	if notice.Link != "" {
		extra := make(map[string]string, 1)
		extra["push_link"] = notice.Link
		ios["extra"] = extra
	}

	wrap["android"] = ios

	return &wrap
}
