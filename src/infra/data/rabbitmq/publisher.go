package rabbitmq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vagner-nascimento/go-adp-archref/src/infra/logger"
)

type rabbitPubInfo struct {
	queue   queueInfo
	message messageInfo
	data    amqp.Publishing
}

func Publish(data []byte, topic string) (err error) {
	pubInfo := newRabbitPubInfo(data, topic)

	if rbCh, err := getChannel(); err == nil {
		if qPub, err := rbCh.QueueDeclare(
			pubInfo.queue.Name,
			pubInfo.queue.Durable,
			pubInfo.queue.AutoDelete,
			pubInfo.queue.Exclusive,
			pubInfo.queue.NoWait,
			pubInfo.queue.Args,
		); err == nil {
			err = rbCh.Publish(
				pubInfo.message.Exchange,
				qPub.Name,
				pubInfo.message.Mandatory,
				pubInfo.message.Immediate,
				pubInfo.data,
			)
		}
	}

	if err != nil {
		msg := fmt.Sprintf("error on publish data into rabbit queue %s", topic)
		logger.Error(msg, err)
		err = errors.New(msg)
	}

	return err
}

func newRabbitPubInfo(data []byte, topic string) rabbitPubInfo {
	return rabbitPubInfo{
		queue: queueInfo{
			Name:       topic,
			Durable:    false,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		},
		message: messageInfo{
			Exchange:  "",
			Mandatory: false,
			Immediate: false,
		},
		data: amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	}
}