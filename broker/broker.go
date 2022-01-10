package broker

import (
	"encoding/json"

	"go-micro.dev/v4/broker"

	"github.com/blackdreamers/core/config"
	"github.com/blackdreamers/core/consts"
	log "github.com/blackdreamers/core/logger"
	"github.com/blackdreamers/core/retry"
)

var (
	subs        []Sub
	subscribers []broker.Subscriber
	nsq         broker.Broker
)

type Sub struct {
	Topic   string
	Handler broker.Handler
	Opts    []broker.SubscribeOption
}

type (
	Event           = broker.Event
	SubscribeOption = broker.SubscribeOption
	PublishOption   = broker.PublishOption
)

type Header map[string]string

type Body map[string]interface{}

func Init(b broker.Broker) error {
	nsq = b

	if err := b.Init(broker.Addrs(config.Broker.Addrs...)); err != nil {
		return err
	}

	if err := b.Connect(); err != nil {
		return err
	}

	for _, sub := range subs {
		sb, err := sub.subscribe()
		if err != nil {
			return err
		}
		subscribers = append(subscribers, sb)
	}

	return nil
}

func EmptyHeader() Header {
	return Header{}
}

func NewBody(keysAndValues ...interface{}) Body {
	body := make(Body)

	if len(keysAndValues) == 0 {
		return body
	}

	for i := 0; i < len(keysAndValues); {
		key := keysAndValues[i]
		if keyStr, ok := key.(string); ok {
			if i+1 < len(keysAndValues) {
				body[keyStr] = keysAndValues[i+1]
			} else {
				body[keyStr] = ""
			}
		}
		i += 2
	}

	return body
}

func SameQueueName() broker.SubscribeOption {
	return broker.Queue(config.Service.SrvName)
}

func DefaultSubOptions() []broker.SubscribeOption {
	return []broker.SubscribeOption{
		broker.DisableAutoAck(),
		SameQueueName(),
	}
}

func Publish(topic string, header Header, body Body, opts ...broker.PublishOption) error {
	bodyJson, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	msg := &broker.Message{
		Header: header,
		Body:   bodyJson,
	}

	err = retry.DefaultRetry.Do(
		func() error {
			return nsq.Publish(topic, msg, opts...)
		},
		retry.OnRetry(func(n uint, err error) {
			log.Field("times", n+1).Log(log.InfoLevel, "retry broker push")
		}),
	)
	if err != nil {
		log.Fields("header", header, "body", string(msg.Body), consts.ErrKey, err).Log(log.ErrorLevel, "broker push")
		return err
	}

	log.Fields("header", header, "body", string(msg.Body)).Log(log.InfoLevel, "broker push")

	return nil
}

func AddSubscribers(sbs ...Sub) {
	subs = append(subs, sbs...)
}

func Subscribers() []broker.Subscriber {
	return subscribers
}

func Broker() broker.Broker {
	return nsq
}

func (s *Sub) subscribe() (broker.Subscriber, error) {
	return nsq.Subscribe(s.Topic, s.Handler, s.Opts...)
}
