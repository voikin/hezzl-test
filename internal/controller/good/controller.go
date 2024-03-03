package good

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/voikin/hezzl-test/internal/service"
	"github.com/voikin/hezzl-test/internal/utils"
)

type GoodController struct {
	goodService service.GoodService
}

func NewGoodController(goodService service.GoodService) *GoodController {
	return &GoodController{goodService: goodService}
}

func (gc *GoodController) Create(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	req := RequestCreate{}
	err = c.ShouldBindJSON(&req)

	if err != nil || req.Name == "" {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	good, err := gc.goodService.CreateGood(c.Request.Context(), req.Name, projectId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, good)
}

func (gc *GoodController) Update(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	goodId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	req := RequestUpdate{}
	err = c.ShouldBindJSON(&req)
	if err != nil || req.Name == "" {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	good, err := gc.goodService.UpdateGood(c.Request.Context(), req.Name, req.Description, goodId, projectId)

	if errors.Is(err, utils.ErrGoodNotFound) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 3, "detail": "{}"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, good)
}

func (gc *GoodController) Delete(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	goodId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	good, err := gc.goodService.DeleteGood(c.Request.Context(), goodId, projectId)

	if errors.Is(err, utils.ErrGoodNotFound) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 3, "detail": "{}"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ResponseRemove{Removed: good.Removed, Id: goodId, ProjectId: projectId})
}

func (gc *GoodController) GetGood(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	goodId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	allGoods, err := gc.goodService.GetGood(c.Request.Context(), goodId, projectId)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, allGoods)
}

func (gc *GoodController) GetGoods(c *gin.Context) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	var limit, offset int
	var err error

	if limitStr == "" {
		limit = 10
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
	}
	if offsetStr == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}
	}

	allGoods, err := gc.goodService.GetGoods(c.Request.Context(), limit, offset)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, allGoods)
}

func (gc *GoodController) UpdateGoodPriority(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	goodId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	req := RequestReprioritize{}
	err = c.ShouldBindJSON(&req)
	if err != nil || req.NewPriority == 0 {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatedPriorities, err := gc.goodService.UpdateGoodPriority(c.Request.Context(), projectId, goodId, req.NewPriority)
	if err != nil {
		if errors.Is(err, utils.ErrGoodNotFound) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 3, "detail": "{}"})
			return
		}

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, updatedPriorities)
}
