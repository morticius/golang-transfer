package db

import (
	"context"
	"testing"

	"github.com/morticius/golang-transfer/util"
	"github.com/stretchr/testify/require"
)

func createRandomCurrency(t *testing.T) Currency {
	checked, _ := testQueries.GetCurrencyByID(context.Background(), 1)

	if checked.ID != 0 {
		return checked
	}

	code := util.RandomCurrency()
	currency, err := testQueries.CreateCurrency(context.Background(), code)

	require.NoError(t, err)
	require.NotEmpty(t, currency)

	require.Equal(t, code, currency.Code)

	require.NotZero(t, currency.ID)
	require.NotZero(t, currency.Code)

	return currency
}

func TestCreateCurrency(t *testing.T) {
	createRandomCurrency(t)
}

func TestGetCurrencyByID(t *testing.T) {
	currencyCreated := createRandomCurrency(t)
	currencySelected, err := testQueries.GetCurrencyByID(context.Background(), currencyCreated.ID)

	require.NoError(t, err)
	require.NotEmpty(t, currencySelected)

	require.Equal(t, currencyCreated.Code, currencySelected.Code)
}

func TestGetAccountByCode(t *testing.T) {
	currencyCreated := createRandomCurrency(t)
	currencySelected, err := testQueries.GetCurrencyByCode(context.Background(), currencyCreated.Code)

	require.NoError(t, err)
	require.NotEmpty(t, currencySelected)

	require.Equal(t, currencyCreated.ID, currencySelected.ID)
}
