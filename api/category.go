package api

import (
	"database/sql"
	db "inventory_management/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createCategoryRequest struct {
	CategoryName string `json:"category_name" binding:"required"`
	SectionName  string `json:"section_name" binding:"required"`
}

func (server *Server) createCategory(c *gin.Context) {
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateCategoryParams{
		CategoryName: req.CategoryName,
		SectionName:  req.SectionName,
	}

	category, err := server.store.CreateCategory(c, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, category)
}

type getCategoryRequest struct {
	ID int64 `uri:"id" binding:"required",min=1`
}

func (server *Server) getCategory(c *gin.Context) {
	var req getCategoryRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.GetCategory(c, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, category)
}

type listCategoryRequest struct {
	PageID   int32 `form:"page_id" binding:"required",min=1`
	PageSize int32 `form:"page_size" binding:"required",min=5,max=10`
}

func (server *Server) listCategory(c *gin.Context) {
	var req listCategoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListCategoriesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	categories, err := server.store.ListCategories(c, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, categories)
}

type updateCategoryRequest struct {
	ID           int64  `uri:"id" binding:"required"`
}

type updateCategoryRequestJson struct {
	CategoryName string `json:"category_name" binding:"required"`
	SectionName  string `json:"section_name" binding:"required"`
}


func (server *Server) updateCategory(c *gin.Context) {
	var req updateCategoryRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var reqUpdate updateCategoryRequestJson
	if err1 := c.ShouldBindJSON(&reqUpdate); err1 != nil {
		c.JSON(http.StatusBadRequest, err1)
		return
	}

	arg := db.UpdateCategoryParams{
		ID:           req.ID,
		CategoryName: reqUpdate.CategoryName,
		SectionName:  reqUpdate.SectionName,
	}

	category, err2 := server.store.UpdateCategory(c, arg)

	if err2 != nil {
		if err2 == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err2))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err2))
		return
	}

	c.JSON(http.StatusOK, category)
}

type deleteCategoryRequest struct {
	ID int64 `uri:"id" binding:"required",min=1`
}

func (server *Server) deleteCategory(c *gin.Context) {
	var req getCategoryRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteCategory(c, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "category deleted successfuly",
	})
}
