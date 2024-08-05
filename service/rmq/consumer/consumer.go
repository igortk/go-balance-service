package consumer

import (
	"balance-service/service/rmq/handler"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Consumer struct {
	ch    *amqp.Channel
	h     handler.Handler
	qName string
	rk    string
	ex    string
}

func New(ch *amqp.Channel, handler handler.Handler, ex, q, rk string) (*Consumer, error) {
	err := ch.ExchangeDeclare(
		ex,
		amqp.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		q,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		queue.Name,
		rk,
		ex,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		ch:    ch,
		h:     handler,
		qName: q,
		ex:    ex,
		rk:    rk,
	}, err
}

func (c *Consumer) Consume() {
	msgs, _ := c.ch.Consume(
		c.qName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	go func() {
		for msg := range msgs {
			if c.h != nil {
				if err := c.h.HandleMessage(msg.Body); err != nil {
					log.Errorf("handler error: %s", err)
				}
			}
		}
	}()
}

func (c *Consumer) ConsumeByCondition(condition func([]byte) bool, timeout time.Duration) (*amqp.Delivery, error) {
	msgs, err := c.ch.Consume(
		c.qName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	timeoutChan := time.After(timeout)
	for {
		select {
		case d := <-msgs:
			if condition(d.Body) {
				return &d, nil
			}
		case <-timeoutChan:
			return nil, errors.New("timeout waiting for message")
		}
	}
}
