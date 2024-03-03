package project

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/voikin/hezzl-test/internal/service"
	"github.com/voikin/hezzl-test/internal/utils"
)

type ProjectController struct {
	projectService service.ProjectService
}

func NewProjectController(projectService service.ProjectService) *ProjectController {
	return &ProjectController{projectService: projectService}
}

func (pc *ProjectController) Create(c *gin.Context) {
	req := &RequestCreate{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	project, err := pc.projectService.CreateProject(c.Request.Context(), req.Name)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, project)
}

func (pc *ProjectController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	req := &RequestCreate{}
	err = c.ShouldBindJSON(req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	project, err := pc.projectService.UpdateProject(c.Request.Context(), req.Name, id)

	if errors.Is(err, utils.ErrProjectNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error(), "code": 3, "detail": "{}"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "", "detail": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, project)
}

func (pc *ProjectController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isDeleted, err := pc.projectService.DeleteProject(c.Request.Context(), id)

	if errors.Is(err, utils.ErrProjectNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error(), "code": 3, "detail": "{}"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "removed": isDeleted})
}

func (pc *ProjectController) GetProjects(c *gin.Context) {
	allGoods, err := pc.projectService.GetProjects(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, allGoods)
}

func (pc *ProjectController) GetProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	allGoods, err := pc.projectService.GetProject(c.Request.Context(), id)

	if errors.Is(err, utils.ErrProjectNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error(), "code": 3, "detail": "{}"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allGoods)
}
