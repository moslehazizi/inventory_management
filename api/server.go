package api

import (
	db "inventory_management/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Server serve HTTP requests for inventory_management services
type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/categories", server.createCategory)
	router.GET("/categories/:id", server.getCategory)
	router.GET("/categories", server.listCategory)
	router.PUT("/categories/:id", server.updateCategory)
	router.DELETE("/categories/:id", server.deleteCategory)
	router.POST("/units", server.createUnit)
	router.GET("/units", server.listUnit)
	router.DELETE("/units/:id", server.deleteUnit)
	router.PUT("/units/:id", server.updateUnit)
	router.POST("/goods", server.createGood)
	router.GET("/goods/:id", server.getGood)
	router.GET("/goods", server.listGood)
	router.PUT("/goods/:id", server.updateGood)
	router.DELETE("/goods/:id", server.deleteGood)

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
