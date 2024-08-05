package rmq

import (
	"balance-service/config"
	"balance-service/service/rmq/consumer"
	"balance-service/service/rmq/handler"
	"github.com/streadway/amqp"
)

type Server struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	handlers  map[string]handler.Handler
	consumers map[string]consumer.Consumer
}

func New(cfg *config.RabbitMqConfig) (*Server, error) {
	conn, err := amqp.Dial(cfg.Url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Server{
		conn: conn,
		ch:   ch,
	}, nil
}

func (s *Server) Run() {

}

func (s *Server) initHandlers() {
	s.handlers["EmitUserBalanceHandler"] = handler.NewEmitUserBalanceHandler()
}

func (s *Server) initConsumers() {

}
