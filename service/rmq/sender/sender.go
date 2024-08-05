package sender

import (
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
)

type Sender struct {
	ch *amqp.Channel
}

func New(ch *amqp.Channel) *Sender {
	return &Sender{
		ch: ch,
	}
}

func (s *Sender) Publish(rk, ex string, msg proto.Message) error {
	body, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	err = s.ch.Publish(
		ex,
		rk,
		false,
		false,
		amqp.Publishing{
			Body: body,
		})

	return err
}
