package server_push

import (
	"encoding/json"
	"log"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/SherClockHolmes/webpush-go"
)

var (
	Subscriber      = "test@example.com"
	VAPIDPublicKey  = ""
	VAPIDPrivateKey = ""
)

type Subscription struct {
	Endpoint string
	Auth     string
	P256dh   string
}

func init() {
	var err error
	VAPIDPrivateKey, VAPIDPublicKey, err = webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatal(err)
	}
}

func Push(subscription *Subscription, body interface{}) error {
	s := &webpush.Subscription{
		Endpoint: subscription.Endpoint,
		Keys: webpush.Keys{
			Auth:   subscription.Auth,
			P256dh: subscription.P256dh,
		},
	}

	message, err := json.Marshal(body)
	if err != nil {
		return errors.ErrCanNotUnmarshal
	}

	resp, err := webpush.SendNotification(message, s, &webpush.Options{
		Subscriber:      Subscriber,
		VAPIDPublicKey:  VAPIDPublicKey,
		VAPIDPrivateKey: VAPIDPrivateKey,
		TTL:             30,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
