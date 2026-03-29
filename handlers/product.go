package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	dtos "github.com/ilwaad12/product-project/dto"
	"github.com/ilwaad12/product-project/infra"
	"github.com/ilwaad12/product-project/repository"
	"github.com/ilwaad12/product-project/service"
)

type ProductHandler struct {
	svc *service.ProductService
}

func RegisterProductHandler() *ProductHandler {
	productRepo := repository.RegisterProductRepo(infra.DB)
	svc := service.RegisterProductService(productRepo)
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var body dtos.CreateProductRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "Failed bind body request",
			"is_success": false,
		})
		return
	}

	resp, status, err := h.svc.Create(body)
	if err != nil {
		c.JSON(status, gin.H{
			"message":    err.Error(),
			"is_success": false,
		})
		return
	}

	c.JSON(status, gin.H{
		"message":    "Product created successfully",
		"is_success": true,
		"data":       resp,
	})
}

func (h *ProductHandler) Update(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "Invalid product id",
			"is_success": false,
		})
		return
	}

	var body dtos.UpdateProductRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "Failed bind body request",
			"is_success": false,
		})
		return
	}

	resp, status, svcErr := h.svc.Update(uint(id64), body)
	if svcErr != nil {
		c.JSON(status, gin.H{
			"message":    svcErr.Error(),
			"is_success": false,
		})
		return
	}

	c.JSON(status, gin.H{
		"message":    "Product updated successfully",
		"is_success": true,
		"data":       resp,
	})
}

func (h *ProductHandler) List(c *gin.Context) {
	var q dtos.ListProductsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "Invalid query parameters",
			"is_success": false,
		})
		return
	}

	// Only allow includeDeleted when admin.
	role := c.GetString("role")
	if role != "ADMIN" {
		q.IncludeDeleted = false
	}

	items, total, page, limit, status, err := h.svc.List(q)
	if err != nil {
		c.JSON(status, gin.H{
			"message":    err.Error(),
			"is_success": false,
			"data":       nil,
		})
		return
	}

	c.JSON(status, gin.H{
		"message":    "success",
		"is_success": true,
		"data": gin.H{
			"items": items,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "Invalid product id",
			"is_success": false,
		})
		return
	}

	// Only admin can fetch deleted records
	includeDeleted := c.GetString("role") == "ADMIN"
	resp, status, svcErr := h.svc.GetByID(uint(id64), includeDeleted)
	if svcErr != nil {
		c.JSON(status, gin.H{
			"message":    svcErr.Error(),
			"is_success": false,
		})
		return
	}

	c.JSON(status, gin.H{
		"message":    "success",
		"is_success": true,
		"data":       resp,
	})
}

func (h *ProductHandler) SoftDelete(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "Invalid product id",
			"is_success": false,
		})
		return
	}

	status, svcErr := h.svc.SoftDelete(uint(id64))
	if svcErr != nil {
		c.JSON(status, gin.H{
			"message":    svcErr.Error(),
			"is_success": false,
		})
		return
	}

	c.JSON(status, gin.H{
		"message":    "Product deleted successfully",
		"is_success": true,
	})
}

func (h *ProductHandler) Restore(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "Invalid product id",
			"is_success": false,
		})
		return
	}

	status, svcErr := h.svc.Restore(uint(id64))
	if svcErr != nil {
		c.JSON(status, gin.H{
			"message":    svcErr.Error(),
			"is_success": false,
		})
		return
	}

	c.JSON(status, gin.H{
		"message":    "Product restored successfully",
		"is_success": true,
	})
}
