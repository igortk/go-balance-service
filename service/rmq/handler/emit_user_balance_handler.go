package handler

import (
	"balance-service/dto/proto"
	"balance-service/storage/postgres"
	gitProto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

type EmitUserBalanceHandler struct {
	pgCl *postgres.Client
}

func NewEmitUserBalanceHandler(pgCl *postgres.Client) *EmitUserBalanceHandler {
	return &EmitUserBalanceHandler{
		pgCl: pgCl,
	}
}

func (h *EmitUserBalanceHandler) HandleMessage(msg []byte) error {
	log.Info("handling request for emit user balance...")
	req := proto.EmmitBalanceByUserIdRequest{}

	err := gitProto.Unmarshal(msg, &req)
	if err != nil {
		return err
	}
	log.Info("request unmarshal successfully")

	err = h.pgCl.UpdateUserBalance(req.UserId, req.Currency, req.Amount, 0)
	if err != nil {
		return err
	}

	log.Info("handler emit user balance successfully")
	return nil
}
