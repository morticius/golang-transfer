package db

import (
	"context"
	"testing"
	"time"

	"github.com/morticius/golang-transfer/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, currency Currency) Transfer {
	arg := CreateTransferParams{
		UserID:     util.RandomInt(1, 10),
		CurrencyID: currency.ID,
		Amount:     util.RandomInt(1000, 10000000),
		TimePlaced: time.Now().UTC(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.UserID, transfer.UserID)
	require.Equal(t, arg.CurrencyID, transfer.CurrencyID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.WithinDuration(t, arg.TimePlaced, transfer.TimePlaced, time.Second)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	currency := createRandomCurrency(t)
	createRandomTransfer(t, currency)
}

func TestListTransfers(t *testing.T) {
	currency := createRandomCurrency(t)

	check := createRandomTransfer(t, currency)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, currency)
	}

	arg := ListTransfersParams{
		UserID: check.UserID,
		Limit:  5,
		Offset: 0,
	}

	_, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
}
