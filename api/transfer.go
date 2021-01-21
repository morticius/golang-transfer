package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/morticius/golang-transfer/db/sqlc"
)

type createTransferRequest struct {
	UserID     string `json:"user_id" binding:"required"`
	Currency   string `json:"currency" binding:"required,oneof=USD EUR"`
	Amount     int64  `json:"amount" binding:"required,gt=0"`
	TimePlaced string `json:"time_placed" binding:"required"`
	Type       string `json:"type" binding:"oneof=deposit withdrawal"`
}

func (server *Server) createTransfer(c *gin.Context) {
	var r createTransferRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	uid, err := strconv.Atoi(r.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	amount := int64(r.Amount * 1000)
	if r.Type == "withdrawal" {
		amount = amount * (-1)
	}

	time, err := time.Parse("02-Jan-06 15:04:05", r.TimePlaced)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cur, err := server.store.GetCurrencyByCode(c, r.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateTransferParams{
		UserID:     int64(uid),
		CurrencyID: cur.ID,
		Amount:     amount,
		TimePlaced: time,
	}

	_, err = server.store.CreateTransfer(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	c.Status(http.StatusCreated)
}
