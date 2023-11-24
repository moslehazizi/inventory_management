package api

import (
	"database/sql"
	db "inventory_management/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createGoodRequest struct {
	Category int64  `json:"category" binding:"required"`
	Model    string `json:"model" binding:"required"`
	Unit     int64  `json:"unit" binding:"required"`
	Amount   int64  `json:"amount" binding:"required"`
	GoodDesc string `json:"good_desc" binding:"required"`
}

func (server *Server) createGood(c *gin.Context) {
	var req createGoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateGoodParams{
		Category: req.Category,
		Model:    req.Model,
		Unit:     req.Unit,
		Amount:   req.Amount,
		GoodDesc: req.GoodDesc,
	}

	good, err := server.store.CreateGood(c, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, good)
}

type getGoodRequest struct {
	ID int64 `uri:"id" binding:"required",min=1`
}

func (server *Server) getGood(c *gin.Context) {
	var req getGoodRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	good, err := server.store.GetGood(c, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, good)
}

type listGoodRequest struct {
	PageID   int32 `form:"page_id" binding:"required",min=1`
	PageSize int32 `form:"page_size" binding:"required",min=5,max=10`
}

type listGoodRequestCategory struct {
	Category int64 `json:"category" binding:"required"`
}

func (server *Server) listGood(c *gin.Context) {
	var req listGoodRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqCat listGoodRequestCategory
	if err := c.ShouldBindJSON(&reqCat); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListGoodsParams{
		Category: reqCat.Category,
		Limit:    req.PageSize,
		Offset:   (req.PageID - 1) * req.PageSize,
	}
	goods, err := server.store.ListGoods(c, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, goods)
}

type updateGoodRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type updateGoodRequestJson struct {
	Unit   int64 `json:"unit" binding:"required"`
	Amount int64 `json:"amount" binding:"required"`
}

func (server *Server) updateGood(c *gin.Context) {
	var req updateGoodRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var reqUpdate updateGoodRequestJson
	if err1 := c.ShouldBindJSON(&reqUpdate); err1 != nil {
		c.JSON(http.StatusBadRequest, err1)
		return
	}

	arg := db.UpdateGoodParams{
		ID:     req.ID,
		Unit:   reqUpdate.Unit,
		Amount: reqUpdate.Amount,
	}

	good, err2 := server.store.UpdateGood(c, arg)

	if err2 != nil {
		if err2 == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err2))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err2))
		return
	}

	c.JSON(http.StatusOK, good)
}


type deleteGoodRequest struct {
	ID int64 `uri:"id" binding:"required",min=1`
}

func (server *Server) deleteGood(c *gin.Context) {
	var req deleteGoodRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteGood(c, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "unit deleted successfuly",
	})
}
