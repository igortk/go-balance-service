package consumer

import (
	"balance-service/config"
	"balance-service/dto/proto"
	"balance-service/service/rmq/sender"
	gitProto "github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	testQueue      = "q.test.balance-service"
	testExchange   = "e.balance.exchange"
	testRoutingKey = "r.balance-service.test"
	testUserId     = "1162cf26-202c-41de-86ad-244ea4ad24bd"
	testCurrency   = "USD"
	testAmount     = 17.25
)

func TestConsumer(t *testing.T) {
	cfg, err := config.Read()
	assert.Nil(t, err)

	conn, err := amqp.Dial(cfg.RabbitMqConfig.Url)
	assert.Nil(t, err)

	channel, err := conn.Channel()
	assert.Nil(t, err)

	clConsumer, err := New(channel, nil, testExchange, testQueue, testRoutingKey)
	assert.Nil(t, err)

	clConsumer.Consume()
}

func TestConsumerWithCondition(t *testing.T) {
	log.Info("start test")

	cfg, err := config.Read()
	assert.Nil(t, err)

	conn, err := amqp.Dial(cfg.RabbitMqConfig.Url)
	assert.Nil(t, err)

	channel, err := conn.Channel()
	assert.Nil(t, err)

	clConsumer, err := New(channel, nil, testExchange, testQueue, testRoutingKey)
	assert.Nil(t, err)

	sdr := sender.New(channel)

	go func() {
		time.Sleep(5 * time.Second)
		mes := &proto.EmmitBalanceByUserIdRequest{
			Id:       uuid.NewString(),
			UserId:   testUserId,
			Currency: testCurrency,
			Amount:   testAmount,
		}

		err = sdr.Publish(testRoutingKey, testExchange, mes)
		assert.Nil(t, err)
	}()

	condition := func(body []byte) bool {
		msg := &proto.EmmitBalanceByUserIdRequest{}
		err = gitProto.Unmarshal(body, msg)

		return msg.UserId == testUserId
	}

	del, err := clConsumer.ConsumeByCondition(condition, 10*time.Second)
	assert.Nil(t, err)

	msg := &proto.EmmitBalanceByUserIdRequest{}
	err = gitProto.Unmarshal(del.Body, msg)
	assert.Nil(t, err)

	log.Infof("recived message: %s", msg)
}
