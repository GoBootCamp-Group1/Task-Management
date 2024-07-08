package notification

import (
	"github.com/GoBootCamp-Group1/Task-Management/pkg/notification/push"
)

const FIREBASE_PROVIDER = "fcm"
const PUSHER_PROVIDER = "pusher"

type PushNotifierConf struct {
	NotifierConf

	Provider string
	AuthKey  string
}

func NewPushNotifier(cfg *PushNotifierConf) (push.IPusher, error) {
	if cfg.Provider == FIREBASE_PROVIDER {
		return &push.FirebasePushNotifier{
			FcmAuthKey: cfg.AuthKey,
		}, nil
	}

	if cfg.Provider == PUSHER_PROVIDER {
		return &push.PusherPushNotifier{
			AuthKey: cfg.AuthKey,
		}, nil
	}
	return nil, nil
}
