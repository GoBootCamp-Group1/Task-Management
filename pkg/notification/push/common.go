package push

type IPusher interface {
	Send(tokens []string, input *PushNotificationInput) error
}

type PushNotificationInput struct {
	Title     string
	Body      string
	ImageUrl  string
	IconUrl   string
	ActionUrl string
	Payload   any
}
