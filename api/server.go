package api

import (
	db "inventory_management/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serve HTTP requests for inventory_management services
type Server struct {
	store *db.Store
	router *gin.Engine
}

func  NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/categories", server.createCategory)
	router.GET("/categories/:id", server.getCategory)
	router.GET("/categories", server.listCategory)
	router.PUT("/categories/:id", server.updateCategory)
	router.DELETE("/categories/:id", server.deleteCategory)

	server.router = router
	return server
}

// Start run the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}