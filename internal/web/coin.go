package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/miles0wu/meme-coin-api/internal/domain"
	"github.com/miles0wu/meme-coin-api/internal/service"
	"github.com/miles0wu/meme-coin-api/pkg/logger"
	"net/http"
	"strconv"
	"time"
)

type CoinHandler struct {
	svc service.CoinService
	l   logger.Logger
}

func NewCoinHandler(svc service.CoinService, l logger.Logger) *CoinHandler {
	return &CoinHandler{
		svc: svc,
		l:   l,
	}
}

func (h *CoinHandler) RegisterRoutes(server *gin.Engine) {
	cg := server.Group("/api/v1/meme-coins")
	// POST /meme-coins
	cg.POST("", h.Create)
	// GET /meme-coins/{id}
	cg.GET("/:id", h.Detail)
	// PUT /meme-coins/{id}
	cg.PUT("/:id", h.Update)
	// DELETE /meme-coins
	cg.DELETE("/:id", h.Delete)
	// POST /meme-coins/{id}/poke
	cg.POST("/:id/poke", h.Poke)
}

// Create is used to add a new meme coin
// @Summary Create meme coin
// @Description Add a new meme coin
// @Tags Coins
// @Accept json
// @Produce json
// @Param payload body CreateCoinReq true "coin"
// @Success 201 {object} Result{data=CoinVo}
// @Failure 400 {object} Result
// @Failure 500 {object} Result
// @Router /api/v1/meme-coins [post]
func (h *CoinHandler) Create(ctx *gin.Context) {
	var req CreateCoinReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Msg:  "invalid input",
			Code: 400,
		})
		h.l.Error("failed to create coin, invalid input",
			logger.Error(err))
		return
	}

	coin, err := h.svc.Create(ctx, domain.Coin{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, service.ErrDuplicateName) {
			ctx.JSON(http.StatusBadRequest, Result{
				Msg:  "coin name already exists",
				Code: 400,
			})
			h.l.Error("failed to create coin, duplicate coin name",
				logger.Error(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 500,
			Msg:  "internal server error",
		})
		h.l.Error("failed to create coin",
			logger.Error(err))
		return
	}

	ctx.JSON(http.StatusCreated, Result{
		Code: 200,
		Data: CoinVo{
			Id:              coin.Id,
			Name:            coin.Name,
			Description:     coin.Description,
			CreatedAt:       coin.CreatedAt.Format(time.DateTime),
			UpdatedAt:       coin.UpdatedAt.Format(time.DateTime),
			PopularityScore: coin.PopularityScore,
		},
	})
	ctx.Header("Location", fmt.Sprintf("/api/v1/meme-coins/%d", coin.Id))
}

// Detail is used to get a coin info by id.
// @Summary Get meme coin
// @Description Get a coin info by id.
// @Tags Coins
// @Accept json
// @Produce json
// @Param id path string true "Coin ID"
// @Success 200 {object} Result{data=CoinVo}
// @Failure 400 {object} Result
// @Failure 404 {object} Result
// @Failure 500 {object} Result
// @Router /api/v1/meme-coins/{id} [get]
func (h *CoinHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Msg:  "invalid id param",
			Code: 400,
		})
		h.l.Error("failed to get coin detail, invalid id",
			logger.Error(err),
			logger.String("id", idStr))
		return
	}

	coin, err := h.svc.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, Result{
				Code: 404,
				Msg:  "coin not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 500,
			Msg:  "internal server error",
		})
		h.l.Error("failed to get coin detail",
			logger.Error(err),
			logger.Int64("id", id))
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Data: CoinVo{
			Id:              coin.Id,
			Name:            coin.Name,
			Description:     coin.Description,
			CreatedAt:       coin.CreatedAt.Format(time.DateTime),
			UpdatedAt:       coin.UpdatedAt.Format(time.DateTime),
			PopularityScore: coin.PopularityScore,
		},
	})
}

// Update is used to modify the description of a meme coin by its ID
// @Summary Update meme coin
// @Description Modify the description of a meme coin by its ID
// @Tags Coins
// @Accept json
// @Produce json
// @Param id path string true "Coin ID"
// @Param payload body UpdateCoinReq true "coin"
// @Success 200 {object} Result
// @Failure 400 {object} Result
// @Failure 500 {object} Result
// @Router /api/v1/meme-coins/{id} [put]
func (h *CoinHandler) Update(ctx *gin.Context) {
	var req UpdateCoinReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Msg:  "invalid input",
			Code: 400,
		})
		h.l.Error("failed to create coin, invalid input",
			logger.Error(err))
		return
	}
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Msg:  "invalid id param",
			Code: 400,
		})
		h.l.Error("failed to update coin, invalid id",
			logger.Error(err),
			logger.String("id", idStr))
		return
	}

	c, err := h.svc.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			ctx.JSON(http.StatusBadRequest, Result{
				Msg:  "invalid id param",
				Code: 400,
			})
			h.l.Error("failed to update coin, input id not found",
				logger.Error(err),
				logger.Int64("id", id))
			return
		}
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 500,
			Msg:  "internal server error",
		})
		h.l.Error("failed to update coin, get coin by id error",
			logger.Error(err),
			logger.Int64("id", id))
		return
	}
	c.Description = req.Description

	err = h.svc.Update(ctx, c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 500,
			Msg:  "internal server error",
		})
		h.l.Error("failed to update coin",
			logger.Error(err),
			logger.Int64("id", id))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "OK",
	})
}

// Delete is used to remove a meme coin by its ID
// @Summary Delete meme coin
// @Description Remove a meme coin by its ID
// @Tags Coins
// @Accept json
// @Produce json
// @Param id path string true "Coin ID"
// @Success 204 {object} Result
// @Failure 400 {object} Result
// @Failure 500 {object} Result
// @Router /api/v1/meme-coins/{id} [delete]
func (h *CoinHandler) Delete(ctx *gin.Context) {
	var err error
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Msg:  "invalid id param",
			Code: 400,
		})
		h.l.Error("failed to delete coin, invalid id",
			logger.Error(err),
			logger.String("id", idStr))
		return
	}

	err = h.svc.DeleteById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 500,
			Msg:  "internal server error",
		})
		h.l.Error("failed to delete coin",
			logger.Error(err),
			logger.Int64("id", id))
		return
	}
	ctx.Status(http.StatusNoContent)
}

// Poke is used to poke a meme coin to show your interest in its ID
// @Summary Poke meme coin
// @Description Poke a meme coin to show your interest in its ID
// @Tags Coins
// @Accept json
// @Produce json
// @Param id path string true "Coin ID"
// @Success 200 {object} Result
// @Failure 400 {object} Result
// @Failure 500 {object} Result
// @Router /api/v1/meme-coins/{id}/poke [post]
func (h *CoinHandler) Poke(ctx *gin.Context) {
	var err error
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Msg:  "invalid id param",
			Code: 400,
		})
		h.l.Error("failed to poke coin, invalid id",
			logger.Error(err),
			logger.String("id", idStr))
		return
	}

	err = h.svc.IncrPopularityScore(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			ctx.JSON(http.StatusBadRequest, Result{
				Msg:  "invalid id param",
				Code: 400,
			})
			h.l.Error("failed to poke coin, input id not found",
				logger.Error(err),
				logger.Int64("id", id))
			return
		}
		ctx.JSON(http.StatusInternalServerError, Result{
			Code: 500,
			Msg:  "internal server error",
		})
		h.l.Error("failed to poke coin",
			logger.Error(err),
			logger.Int64("id", id))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "OK",
	})
}
