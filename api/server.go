package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/samirprakash/go-bank/db/sqlc"
	"github.com/samirprakash/go-bank/token"
	"github.com/samirprakash/go-bank/util"
)

// Server serves HTTP requests for the banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker : %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("customCurrencyValidator", validateCurrency)
	}

	// define routes
	router.POST("/users", server.createUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PATCH("/accounts/:id", server.updateAccountBalance)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server, nil
}

// Start starts a new server on the specified address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// returns custom error messages
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
