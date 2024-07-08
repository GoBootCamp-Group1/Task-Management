package notification

import (
	"errors"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/notification/push"
	"sync"
)

const EMAIL_NOTIFIER = "email"
const DB_NOTIFIER = "db"
const PUSH_NOTIFIER = "push"
const SMS_NOTIFIER = "sms"

type Notifier struct {
	InApp *DatabaseNotifierConfInternal
	Email *EmailNotifierConf
	Push  push.IPusher
	SMS   *SMSNotifierConf
}

var (
	instance *Notifier
	once     sync.Once
)

// NotifierConf General use for all configs
type NotifierConf any

func NewNotifier(notifiersConf map[string]NotifierConf) (*Notifier, error) {
	var err error

	once.Do(func() {
		var databaseNotifier *DatabaseNotifierConfInternal
		var emailNotifier *EmailNotifierConf
		var pushNotifier push.IPusher
		var smsNotifier *SMSNotifierConf

		if hasConfig(notifiersConf, DB_NOTIFIER) {
			dbNotifyCfg, ok := notifiersConf[DB_NOTIFIER].(*DatabaseNotifierConf)
			if !ok {
				err = errors.New("invalid configuration for database notifier")
				return
			}
			databaseNotifier, err = NewDatabaseNotifier(dbNotifyCfg)
			if err != nil {
				return
			}
		}

		if hasConfig(notifiersConf, EMAIL_NOTIFIER) {
			emailCfg, ok := notifiersConf[EMAIL_NOTIFIER].(*EmailNotifierConf)
			if !ok {
				err = errors.New("invalid configuration for email notifier")
				return
			}
			emailNotifier, err = NewEmailNotifier(emailCfg)
			if err != nil {
				return
			}
		}

		if hasConfig(notifiersConf, PUSH_NOTIFIER) {
			pushCfg, ok := notifiersConf[PUSH_NOTIFIER].(*PushNotifierConf)
			if !ok {
				err = errors.New("invalid configuration for push notifier")
				return
			}
			pushNotifier, err = NewPushNotifier(pushCfg)
			if err != nil {
				return
			}
		}

		if hasConfig(notifiersConf, SMS_NOTIFIER) {
			smsCfg, ok := notifiersConf[SMS_NOTIFIER].(*SMSNotifierConf)
			if !ok {
				err = errors.New("invalid configuration for sms notifier")
				return
			}
			smsNotifier, err = NewSMSNotifier(smsCfg)
			if err != nil {
				return
			}
		}

		instance = &Notifier{
			InApp: databaseNotifier,
			Email: emailNotifier,
			Push:  pushNotifier,
			SMS:   smsNotifier,
		}
	})

	if instance == nil && err != nil {
		return nil, err
	}
	return instance, nil
}

func hasConfig(notifiersConf map[string]NotifierConf, item string) bool {
	for k, _ := range notifiersConf {
		if k == item {
			return true
		}
	}
	return false
}
