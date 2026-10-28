package category

import (
	"net/http"

	"go-gin-boilerplate/internal/constants"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	service CategoryService
}

func NewCategoryController(service CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

func (ctrl *CategoryController) FindAll(c *gin.Context) {
	var params FindAllCategoryDTO

	if c.Request.Method == http.MethodGet {
		if err := c.ShouldBindQuery(&params); err != nil {
			_ = c.Error(err)
			return
		}
	} else if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&params); err != nil {
			_ = c.Error(err)
			return
		}
	}

	result, err := ctrl.service.FindAll(params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(constants.StatusOK, result)
}

func (ctrl *CategoryController) FindOne(c *gin.Context) {
	var uriParam FindByIDCategoryDTO
	if err := c.ShouldBindUri(&uriParam); err != nil {
		_ = c.Error(err)
		return
	}

	result, err := ctrl.service.FindByID(uriParam.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(constants.StatusOK, result)
}

func (ctrl *CategoryController) Create(c *gin.Context) {
	var payload CreateCategoryDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(err)
		return
	}

	category, err := ctrl.service.Create(payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(constants.StatusCreated, gin.H{
		"message":  constants.MessageSuccess,
		"category": category,
	})
}

func (ctrl *CategoryController) Update(c *gin.Context) {
	var uriParam FindByIDCategoryDTO
	if err := c.ShouldBindUri(&uriParam); err != nil {
		_ = c.Error(err)
		return
	}

	var payload UpdateCategoryDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(err)
		return
	}
	payload.ID = uriParam.ID

	category, err := ctrl.service.Update(payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(constants.StatusOK, gin.H{
		"message":  constants.MessageSuccess,
		"category": category,
	})
}

func (ctrl *CategoryController) Destroy(c *gin.Context) {
	var uriParam FindByIDCategoryDTO
	if err := c.ShouldBindUri(&uriParam); err != nil {
		_ = c.Error(err)
		return
	}

	if err := ctrl.service.Delete(uriParam.ID); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(constants.StatusOK, gin.H{
		"message": constants.MessageSuccess,
	})
}
