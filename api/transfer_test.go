package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mock_sqlc "github.com/morticius/golang-transfer/db/mock"
	db "github.com/morticius/golang-transfer/db/sqlc"
	"github.com/morticius/golang-transfer/util"
	"github.com/stretchr/testify/require"
)

func trasfer() db.Transfer {
	timePlaced, _ := time.Parse("02-Jan-06 15:04:05", "21-JAN-21 15:15:15")
	return db.Transfer{
		ID:         util.RandomInt(1, 1000),
		UserID:     util.RandomInt(1, 10),
		CurrencyID: 1,
		Amount:     util.RandomInt(100, 10000),
		TimePlaced: timePlaced,
		CreatedAt:  time.Now(),
	}
}

func currency() db.Currency {
	return db.Currency{
		ID:        1,
		Code:      "USD",
		CreatedAt: time.Now(),
	}
}

func TestCreateTransferAPI(t *testing.T) {
	tr := trasfer()
	cur := currency()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_sqlc.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"user_id":     strconv.Itoa(int(tr.UserID)),
				"currency":    "USD",
				"amount":      tr.Amount,
				"time_placed": tr.TimePlaced.Format("02-Jan-06 15:04:05"),
				"type":        "deposit",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {

				arg := db.CreateTransferParams{
					UserID:     tr.UserID,
					CurrencyID: 1,
					Amount:     tr.Amount * 1000,
					TimePlaced: tr.TimePlaced.UTC(),
				}

				store.EXPECT().
					GetCurrencyByCode(gomock.Any(), gomock.Eq(cur.Code)).
					Times(1).
					Return(cur, nil)

				store.EXPECT().
					CreateTransfer(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(tr, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"user_id":     strconv.Itoa(int(tr.UserID)),
				"currency":    "USD",
				"amount":      tr.Amount,
				"time_placed": tr.TimePlaced.Format("02-Jan-06 15:04:05"),
				"type":        "deposit",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.EXPECT().
					GetCurrencyByCode(gomock.Any(), gomock.Eq(cur.Code)).
					Times(1).
					Return(cur, nil)

				store.EXPECT().
					CreateTransfer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Transfer{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidUserID",
			body: gin.H{
				"user_id":     "invalid",
				"currency":    "USD",
				"amount":      tr.Amount,
				"time_placed": tr.TimePlaced.Format("02-Jan-06 15:04:05"),
				"type":        "deposit",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"user_id":     strconv.Itoa(int(tr.UserID)),
				"currency":    "invalid",
				"amount":      tr.Amount,
				"time_placed": tr.TimePlaced.Format("02-Jan-06 15:04:05"),
				"type":        "deposit",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidAmount",
			body: gin.H{
				"user_id":     strconv.Itoa(int(tr.UserID)),
				"currency":    "USD",
				"amount":      "invalid",
				"time_placed": tr.TimePlaced.Format("02-Jan-06 15:04:05"),
				"type":        "deposit",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidTimePlaced",
			body: gin.H{
				"user_id":     strconv.Itoa(int(tr.UserID)),
				"currency":    "USD",
				"amount":      tr.Amount,
				"time_placed": "invalid",
				"type":        "deposit",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidType",
			body: gin.H{
				"user_id":     strconv.Itoa(int(tr.UserID)),
				"currency":    "USD",
				"amount":      tr.Amount,
				"time_placed": tr.TimePlaced.Format("02-Jan-06 15:04:05"),
				"type":        "invalid",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		ts := testCases[i]

		t.Run(ts.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_sqlc.NewMockStore(ctrl)
			ts.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(ts.body)
			require.NoError(t, err)

			url := "/transfers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			ts.checkResponse(recorder)
		})
	}
}
