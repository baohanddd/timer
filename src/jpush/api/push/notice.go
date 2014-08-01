package push

// type Notification struct {
// 	Android AndroidNotice
// 	Ios     IosNotice

// }

type Notice struct {
	Alert string `json:"alert"`
}

type AndroidNotice struct {
	Object NoticeAndroid `json:"android"`
}

type IosNotice struct {
	Wrapper NoticeIos `json:"ios"`
}

type NoticeAndroid struct {
	Alert     string            `json:"alert"`
	Title     string            `json:"title"`
	BuilderId int               `json:"builder_id"`
	Extras    map[string]string `json:"extras"`
}

type NoticeIos struct {
	Alert            string            `json:"alert"`
	Sound            string            `json:"sound"`
	Badge            int               `json:"badge"`
	ContentAvailable int               `json:"content-available"`
	Extras           map[string]string `json:"extras"`
}

func (this *AndroidNotice) SetAlert(alert string) {
	this.Object.Alert = alert
}

func (this *AndroidNotice) SetTitle(title string) {
	this.Object.Title = title
}

func (this *AndroidNotice) SetBuilderId(id int) {
	this.Object.BuilderId = id
}

func (this *AndroidNotice) SetExtras(key, value string) {
	if this.Object.Extras == nil {
		this.Object.Extras = make(map[string]string)
	}
	this.Object.Extras[key] = value
}

func (o *IosNotice) SetAlert(alert string) {
	o.Wrapper.Alert = alert
}

func (o *IosNotice) SetSound(sound string) {
	o.Wrapper.Sound = sound
}

func (o *IosNotice) SetBadge(badge int) {
	o.Wrapper.Badge = badge
}

func (o *IosNotice) SetContentAvailable(available int) {
	o.Wrapper.ContentAvailable = available
}

func (o *IosNotice) SetExtras(key, value string) {
	if o.Wrapper.Extras == nil {
		o.Wrapper.Extras = make(map[string]string)
	}
	o.Wrapper.Extras[key] = value
}
