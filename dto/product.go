package dtos

type CreateProductRequest struct {
	Name          string  `json:"name" binding:"required,min=2,max=200"`
	SKU           string  `json:"sku" binding:"required,min=2,max=64"`
	PurchasePrice float64 `json:"purchasePrice" binding:"required,gte=0"`
	SellingPrice  float64 `json:"sellingPrice" binding:"required,gte=0"`
	Quantity      int     `json:"quantity" binding:"required,gte=0"`
	CategoryID    uint    `json:"categoryId" binding:"required,gt=0"`
}

type UpdateProductRequest struct {
	Name          *string  `json:"name" binding:"omitempty,min=2,max=200"`
	SKU           *string  `json:"sku" binding:"omitempty,min=2,max=64"`
	PurchasePrice *float64 `json:"purchasePrice" binding:"omitempty,gte=0"`
	SellingPrice  *float64 `json:"sellingPrice" binding:"omitempty,gte=0"`
	Quantity      *int     `json:"quantity" binding:"omitempty,gte=0"`
	CategoryID    *uint    `json:"categoryId" binding:"omitempty,gt=0"`
}

type ListProductsQuery struct {
	Page int `form:"page" binding:"omitempty,min=1"`
	// Limit max 200
	Limit int `form:"limit" binding:"omitempty,min=1,max=200"`

	Q          string `form:"q" binding:"omitempty,max=200"` // name or sku
	CategoryID uint   `form:"categoryId" binding:"omitempty"`

	MinPurchase float64 `form:"minPurchase" binding:"omitempty,gte=0"`
	MaxPurchase float64 `form:"maxPurchase" binding:"omitempty,gte=0"`
	MinSelling  float64 `form:"minSelling" binding:"omitempty,gte=0"`
	MaxSelling  float64 `form:"maxSelling" binding:"omitempty,gte=0"`

	// Sorting
	SortBy string `form:"sortBy" binding:"omitempty,oneof=name sku quantity createdAt purchasePrice sellingPrice"`
	Order  string `form:"order" binding:"omitempty,oneof=asc desc"`

	// Admin-only usage: include soft-deleted rows
	IncludeDeleted bool `form:"includeDeleted"`
}
