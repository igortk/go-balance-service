package sender

import (
	"balance-service/config"
	"balance-service/dto/proto"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testExchange   = "e.balance.exchange"
	testRoutingKey = "r.balance-service.test"
	testUserId     = "1162cf26-202c-41de-86ad-244ea4ad24bd"
	testCurrency   = "USD"
	testAmount     = 17.25
)

func TestSendMessage(t *testing.T) {
	cfg, err := config.Read()
	assert.Nil(t, err)

	conn, err := amqp.Dial(cfg.RabbitMqConfig.Url)
	assert.Nil(t, err)

	channel, err := conn.Channel()
	assert.Nil(t, err)

	sender := New(channel)

	mes := &proto.EmmitBalanceByUserIdRequest{
		Id:       uuid.NewString(),
		UserId:   testUserId,
		Currency: testCurrency,
		Amount:   testAmount,
	}

	err = sender.Publish(testRoutingKey, testExchange, mes)
	assert.Nil(t, err)
}
