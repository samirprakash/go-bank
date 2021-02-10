package api

import (
	"database/sql"
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

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	// validate request
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get account from db
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		// if account does not exist
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		// if there is an issue on the server
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// return the account
	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest

	// validate query params
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// return 400 if bad request
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	// get all accounts as per limit and offset
	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		// return 500 if internal error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// return the list of accounts
	ctx.JSON(http.StatusOK, accounts)
}
