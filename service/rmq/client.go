package rmq

import (
	"balance-service/service/rmq/consumer"
	"balance-service/service/rmq/handler"
	"balance-service/service/rmq/sender"
	gitProto "github.com/golang/protobuf/proto"
	"time"
)

type Routing struct {
	Queue      string
	RoutingKey string
	Exchange   string
}

type Client struct {
	C *consumer.Consumer
	H handler.Handler
	S *sender.Sender
}

func New(c *consumer.Consumer, h handler.Handler, s *sender.Sender) *Client {
	return &Client{
		C: c,
		H: h,
		S: s,
	}
}

func (c *Client) Send(message gitProto.Message, r *Routing) error {
	err := c.S.Publish(r.RoutingKey, r.Exchange, message)
	return err
}

func (c *Client) SendAndConsumeResponse(message gitProto.Message, respPoint gitProto.Message,
	condition func([]byte) bool, r *Routing, timeout time.Duration) error {

	err := c.Send(message, r)
	if err != nil {
		return err
	}

	del, err := c.C.ConsumeByCondition(condition, timeout)
	if err != nil {
		return err
	}

	err = gitProto.Unmarshal(del.Body, respPoint)
	if err != nil {
		return err
	}

	return nil
}
