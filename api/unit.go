package api

import (
	"database/sql"
	db "inventory_management/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createUnitRequest struct {
	UnitName  string `json:"unit_name" binding:"required"`
	UnitValue int64  `json:"unit_value" binding:"required"`
}

func (server *Server) createUnit(c *gin.Context) {
	var req createUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUnitParams{
		UnitName:  req.UnitName,
		UnitValue: req.UnitValue,
	}

	unit, err := server.store.CreateUnit(c, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, unit)
}

type listUnitRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUnit(c *gin.Context) {
	var req listUnitRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUnitsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	units, err := server.store.ListUnits(c, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, units)
}

type updateUnitRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type updateUnitRequestJson struct {
	UnitName  string `json:"unit_name" binding:"required"`
	UnitValue int64 `json:"unit_value" binding:"required"`
}

func (server *Server) updateUnit(c *gin.Context) {
	var req updateUnitRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var reqUpdate updateUnitRequestJson
	if err1 := c.ShouldBindJSON(&reqUpdate); err1 != nil {
		c.JSON(http.StatusBadRequest, err1)
		return
	}

	arg := db.UpdateUnitParams{
		ID:        req.ID,
		UnitName:  reqUpdate.UnitName,
		UnitValue: reqUpdate.UnitValue,
	}

	unit, err2 := server.store.UpdateUnit(c, arg)

	if err2 != nil {
		if err2 == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err2))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err2))
		return
	}

	c.JSON(http.StatusOK, unit)
}

type deleteUnitRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteUnit(c *gin.Context) {
	var req deleteUnitRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteUnit(c, req.ID)

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
