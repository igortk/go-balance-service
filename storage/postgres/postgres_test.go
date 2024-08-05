package postgres

import (
	"balance-service/config"
	"balance-service/dto"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testUserId = "1162cf26-202c-41de-86ad-244ea4ad24bd"

func TestCreateBalance(t *testing.T) {
	expected := []dto.Balance{
		dto.Balance{Id: uuid.MustParse("9493f8f9-fd33-4a80-90bf-5a71b8f00591"),
			UserId:     "1162cf26-202c-41de-86ad-244ea4ad24bd",
			CurrencyId: "ad9f2451-ed0a-4b14-b582-6a99c0cdecfa",
			Currency: dto.Currency{
				Id:   "ad9f2451-ed0a-4b14-b582-6a99c0cdecfa",
				Name: "UAH",
			},
			Balance:       13.45,
			LockedBalance: 3.48,
			UpdatedAt:     1719477726,
		},
		dto.Balance{Id: uuid.MustParse("4828d7cf-64f2-4604-a51b-5e19a24e4fdc"),
			UserId:     "1162cf26-202c-41de-86ad-244ea4ad24bd",
			CurrencyId: "f2b4a3f5-6017-448a-988f-c042a8bc8259",
			Currency: dto.Currency{
				Id:   "f2b4a3f5-6017-448a-988f-c042a8bc8259",
				Name: "USD",
			},
			Balance:       20255.4,
			LockedBalance: 12000.17,
			UpdatedAt:     1719218526,
		},
	}

	cfg, err := config.Read()
	assert.Nil(t, err)

	dbCl, err := New(&cfg.PostgresConfig)
	assert.Nil(t, err)

	actual := dbCl.GetBalanceByUserId(testUserId)
	log.Infof("Got user balances: %v", actual)
	assert.Equal(t, actual, expected)
}

func TestUpdateBalance(t *testing.T) {
	cfg, err := config.Read()
	assert.Nil(t, err)

	dbCl, err := New(&cfg.PostgresConfig)
	assert.Nil(t, err)

	err = dbCl.UpdateUserBalance(testUserId, "UAH", 20255.4, 12000.17)
	assert.Nil(t, err)

	actual := dbCl.GetBalanceByUserId(testUserId)
	log.Infof("get user balances: %v", actual)
}
