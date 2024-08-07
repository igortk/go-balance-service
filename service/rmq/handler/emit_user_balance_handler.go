package handler

import (
	"balance-service/dto/proto"
	"balance-service/service/rmq/common"
	"balance-service/service/rmq/sender"
	"balance-service/storage/postgres"
	gitProto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

type EmitUserBalanceHandler struct {
	pgCl *postgres.Client
	s    *sender.Sender
}

func NewEmitUserBalanceHandler(pgCl *postgres.Client, s *sender.Sender) *EmitUserBalanceHandler {
	log.Info("creating new emit user balance handler...")
	return &EmitUserBalanceHandler{
		pgCl: pgCl,
		s:    s,
	}
}

func (h *EmitUserBalanceHandler) HandleMessage(msg []byte) error {
	log.Info("handling request for emit user balance...")
	req := proto.EmmitBalanceByUserIdRequest{}

	err := gitProto.Unmarshal(msg, &req)
	if err != nil {
		return err
	}
	log.Infof("request (id: %v) unmarshal successfully", req.Id)

	err = h.pgCl.UpdateUserBalance(req.UserId, req.Currency, req.Amount, 0)
	if err != nil {
		return err
	}

	balances, err := h.pgCl.GetBalanceByUserId(req.UserId)
	if err != nil {
		return err
	}

	balance := &proto.Balance{}
	for _, bal := range balances {
		if bal.Currency.Name == req.Currency {
			balance = &proto.Balance{
				Currency:      bal.Currency.Name,
				Balance:       float64(bal.Balance),
				LockedBalance: float64(bal.LockedBalance),
				UpdatedDate:   bal.UpdatedAt,
			}
			break
		}
	}

	resp := &proto.UserBalance{
		UserId: req.UserId,
		Balances: []*proto.Balance{
			balance,
		},
	}

	err = h.s.Publish(common.EmitRespRouting, common.BalanceExchange, resp)
	log.Info("handler emit user balance successfully")
	return nil
}
