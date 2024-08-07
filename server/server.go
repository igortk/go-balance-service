package server

import (
	"balance-service/config"
	"balance-service/service/rmq"
	"balance-service/service/rmq/common"
	"balance-service/service/rmq/consumer"
	"balance-service/service/rmq/handler"
	"balance-service/service/rmq/sender"
	"balance-service/storage/postgres"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Server struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	rmqCls map[string]*rmq.Client
	pgCl   *postgres.Client
	s      *sender.Sender
}

func New(cfg *config.Config) (*Server, error) {
	log.Info("creating new server...")

	pgCl, err := postgres.New(&cfg.PostgresConfig)

	conn, err := amqp.Dial(cfg.RabbitMqConfig.Url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Server{
		conn:   conn,
		ch:     ch,
		pgCl:   pgCl,
		rmqCls: map[string]*rmq.Client{},
	}, nil
}

func (s *Server) Init() {
	s.initSender()
	s.initEmitUserBalanceClient()
}

func (s *Server) Run() {
	log.Info("starting consumers...")

	for _, csm := range s.rmqCls {
		csm.C.Consume()
	}

	select {}
}

func (s *Server) initEmitUserBalanceClient() {
	h := handler.NewEmitUserBalanceHandler(s.pgCl, s.s)
	c, err := consumer.New(s.ch,
		h,
		common.BalanceExchange,
		common.QueueEmitReq,
		common.EmitReqRouting,
	)

	s.rmqCls["EmitUserBalanceHandler"] = rmq.New(c, h, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) initSender() {
	s.s = sender.New(s.ch)
}
