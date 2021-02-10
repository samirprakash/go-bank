package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/samirprakash/go-bank/db/sqlc"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR GBP INR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	// validate request params
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// create db params
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	// save to db
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// return response
	ctx.JSON(http.StatusOK, account)
}
