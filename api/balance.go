package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BalanceResponse struct {
	Code    string  `json:"user_id"`
	Balance float64 `json:"balance"`
}

func (server *Server) GetBalance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	balance, err := server.store.GetBalance(c, int64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var res []*BalanceResponse

	for _, val := range balance {
		res = append(res, &BalanceResponse{
			Code:    val.Code,
			Balance: float64(val.Balance) / 1000,
		})
	}

	c.JSON(http.StatusOK, res)
}
