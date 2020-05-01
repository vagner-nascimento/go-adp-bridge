package integration

import (
	"github.com/vagner-nascimento/go-adp-archref/config"
	"github.com/vagner-nascimento/go-adp-archref/src/infra/repository"
)

type sellerSub struct {
	topic    string
	consumer string
	handler  func(data []byte)
}

func (es *sellerSub) GetTopic() string {
	return es.topic
}

func (es *sellerSub) GetConsumer() string {
	return es.consumer
}

func (es *sellerSub) GetHandler() func([]byte) {
	return es.handler
}

func newSellerSub() repository.Subscription {
	sellerConfig := config.Get().Integration.Amqp.Subs.Seller
	return &sellerSub{
		topic:    sellerConfig.Topic,
		consumer: sellerConfig.Consumer,
		handler: func(data []byte) {
			createAccount(data)
		},
	}
}