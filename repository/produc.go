package repository

import (
	"strings"

	"github.com/ilwaad12/product-project/models"
	"gorm.io/gorm"
)

type ProductRepo struct {
	DB *gorm.DB
}

func RegisterProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{DB: db}
}

func (repo *ProductRepo) Create(p *models.Product) error {
	return repo.DB.Create(p).Error
}

func (repo *ProductRepo) GetByID(id uint, includeDeleted bool) (models.Product, error) {
	var p models.Product
	q := repo.DB
	if includeDeleted {
		q = q.Unscoped()
	}
	err := q.Preload("Category").Where("id = ?", id).First(&p).Error
	if err != nil {
		return models.Product{}, err
	}
	return p, nil
}

func (repo *ProductRepo) ExistsSKU(sku string, excludeID uint) (bool, error) {
	var count int64
	q := repo.DB.Model(&models.Product{}).Where("sku = ?", sku)
	if excludeID != 0 {
		q = q.Where("id <> ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

type ListProductsParams struct {
	Page, Limit int
	Q           string
	CategoryID  uint
	MinPurchase float64
	MaxPurchase float64
	MinSelling  float64
	MaxSelling  float64
	SortBy      string
	Order       string
	IncludeDel  bool
}

func (repo *ProductRepo) List(params ListProductsParams) ([]models.Product, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 20
	}
	if params.Limit > 200 {
		params.Limit = 200
	}

	qb := repo.DB.Model(&models.Product{}).Preload("Category")
	if params.IncludeDel {
		qb = qb.Unscoped()
	}

	if params.Q = strings.TrimSpace(params.Q); params.Q != "" {
		like := "%" + strings.ToLower(params.Q) + "%"
		qb = qb.Where("LOWER(name) LIKE ? OR LOWER(sku) LIKE ?", like, like)
	}
	if params.CategoryID != 0 {
		qb = qb.Where("category_id = ?", params.CategoryID)
	}
	if params.MinPurchase > 0 {
		qb = qb.Where("purchase_price >= ?", params.MinPurchase)
	}
	if params.MaxPurchase > 0 {
		qb = qb.Where("purchase_price <= ?", params.MaxPurchase)
	}
	if params.MinSelling > 0 {
		qb = qb.Where("selling_price >= ?", params.MinSelling)
	}
	if params.MaxSelling > 0 {
		qb = qb.Where("selling_price <= ?", params.MaxSelling)
	}

	var total int64
	if err := qb.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortCol := "id"
	switch params.SortBy {
	case "name":
		sortCol = "name"
	case "sku":
		sortCol = "sku"
	case "quantity":
		sortCol = "quantity"
	case "createdAt":
		sortCol = "created_at"
	case "purchasePrice":
		sortCol = "purchase_price"
	case "sellingPrice":
		sortCol = "selling_price"
	}

	order := "desc"
	if strings.ToLower(params.Order) == "asc" {
		order = "asc"
	}

	offset := (params.Page - 1) * params.Limit
	var items []models.Product
	if err := qb.Order(sortCol + " " + order).Limit(params.Limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (repo *ProductRepo) UpdateFields(id uint, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	return repo.DB.Model(&models.Product{}).Where("id = ?", id).Updates(updates).Error
}

func (repo *ProductRepo) SoftDelete(id uint) error {
	return repo.DB.Delete(&models.Product{}, id).Error
}

func (repo *ProductRepo) Restore(id uint) error {
	// Unscoped update to clear deleted_at
	return repo.DB.Unscoped().Model(&models.Product{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (repo *ProductRepo) CountByCategoryID(categoryID uint) (int64, error) {
	var count int64
	err := repo.DB.Model(&models.Product{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}
