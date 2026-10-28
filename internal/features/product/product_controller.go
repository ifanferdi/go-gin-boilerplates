package product

import (
	"net/http"

	"go-gin-boilerplate/internal/constants"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service ProductService
}

func NewProductController(service ProductService) *ProductController {
	return &ProductController{service: service}
}

func (ctrl *ProductController) FindAll(c *gin.Context) {
	var params FindAllProductDTO

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

func (ctrl *ProductController) FindOne(c *gin.Context) {
	var uriParam FindByIDProductDTO
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

func (ctrl *ProductController) Create(c *gin.Context) {
	var payload CreateProductDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(err)
		return
	}

	product, err := ctrl.service.Create(payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(constants.StatusCreated, gin.H{
		"message": constants.MessageSuccess,
		"product": product,
	})
}

func (ctrl *ProductController) Update(c *gin.Context) {
	var uriParam FindByIDProductDTO
	if err := c.ShouldBindUri(&uriParam); err != nil {
		_ = c.Error(err)
		return
	}

	var payload UpdateProductDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(err)
		return
	}
	payload.ID = uriParam.ID

	product, err := ctrl.service.Update(payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(constants.StatusOK, gin.H{
		"message": constants.MessageSuccess,
		"product": product,
	})
}

func (ctrl *ProductController) Destroy(c *gin.Context) {
	var uriParam FindByIDProductDTO
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
