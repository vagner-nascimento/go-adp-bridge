package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

func SubscribeConsumer(queueName string, consumerName string, handler func([]byte)) (err error) {
	var rbChan *amqp.Channel
	rbChan, err = getChannel()

	if err == nil {
		sub := newSubscriberInfo(queueName, consumerName, handler)
		err = processMessages(rbChan, sub)
	}

	return
}

func newSubscriberInfo(queueName string, consumerName string, handler func([]byte)) rabbitSubInfo {
	return rabbitSubInfo{
		queue: queueInfo{
			Name:         queueName,
			Durable:      false,
			DeleteUnused: false,
			Exclusive:    false,
			NoWait:       false,
		},
		message: messageInfo{
			Consumer:  consumerName,
			AutoAct:   true,
			Exclusive: false,
			Local:     false,
			NoWait:    false,
		},
		handler: handler,
	}
}

func processMessages(ch *amqp.Channel, sub rabbitSubInfo) (err error) {
	var q amqp.Queue
	q, err = ch.QueueDeclare(
		sub.queue.Name,
		sub.queue.Durable,
		sub.queue.DeleteUnused,
		sub.queue.Exclusive,
		sub.queue.NoWait,
		sub.queue.Args,
	)

	if err == nil {
		var msgs <-chan amqp.Delivery
		msgs, err = ch.Consume(
			q.Name,
			sub.message.Consumer,
			sub.message.AutoAct,
			sub.message.Exclusive,
			sub.message.Local,
			sub.message.NoWait,
			nil,
		)

		if err == nil {
			go func(queueName string, msgs <-chan amqp.Delivery) {
				for msg := range msgs {
					fmt.Println(fmt.Sprintf("message received from %s. body:\r\n %s", q.Name, string(msg.Body)))
					sub.handler(msg.Body)
				}
			}(q.Name, msgs)
		}
	}

	return
}
