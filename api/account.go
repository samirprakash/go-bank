package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/samirprakash/go-bank/db/sqlc"
	"github.com/samirprakash/go-bank/token"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,customCurrencyValidator"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	// validate request params
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// add auth to create account
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	// create db params
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}

	// save to db
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

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

	// add auth middleware
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err = errors.New("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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

	// add auth middleware
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
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

type updateAccountBalanceRequest struct {
	Amount int64 `json:"amount" binding:"required,min=1,max=10000"`
}

func (server *Server) updateAccountBalance(ctx *gin.Context) {
	var req updateAccountBalanceRequest

	// validate path param id
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		// return 400 if bad request
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Account ID must not be less than 1
	if id < 1 {
		// return 400 if bad request
		ctx.JSON(http.StatusBadRequest, "Invalid ID")
		return
	}

	// validate requets body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// return 400 if bad request
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAccountBalanceParams{
		Amount: req.Amount,
		ID:     int64(id),
	}

	// update account balance in db
	account, err := server.store.UpdateAccountBalance(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			// reeturn 404 if account not found
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		// return 500 if internal error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// return the account with updated balance
	ctx.JSON(http.StatusOK, account)
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Accoun deleted!")
}
