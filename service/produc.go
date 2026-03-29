package service

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ilwaad12/product-project/constant"
	dtos "github.com/ilwaad12/product-project/dto"
	"github.com/ilwaad12/product-project/models"
	"github.com/ilwaad12/product-project/repository"
	"gorm.io/gorm"
)

type ProductService struct {
	productRepo *repository.ProductRepo
}

func RegisterProductService(productRepo *repository.ProductRepo) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (svc *ProductService) Create(req dtos.CreateProductRequest) (*models.Product, int, error) {
	name := strings.TrimSpace(req.Name)
	sku := strings.TrimSpace(req.SKU)
	if name == "" || sku == "" {
		return nil, http.StatusBadRequest, errors.New("name and sku are required")
	}
	if req.Quantity < 0 {
		return nil, http.StatusBadRequest, errors.New("quantity cannot be negative")
	}
	if req.SellingPrice < req.PurchasePrice {
		return nil, http.StatusBadRequest, errors.New("sellingPrice must be greater than or equal to purchasePrice")
	}

	exists, err := svc.productRepo.ExistsSKU(sku, 0)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New(constant.DefaultErrorMsg)
	}
	if exists {
		return nil, http.StatusConflict, errors.New("sku already exists")
	}

	p := models.Product{
		Name:          name,
		SKU:           sku,
		PurchasePrice: req.PurchasePrice,
		SellingPrice:  req.SellingPrice,
		Quantity:      req.Quantity,
		CategoryID:    req.CategoryID,
	}
	if err := svc.productRepo.Create(&p); err != nil {
		return nil, http.StatusInternalServerError, errors.New(constant.DefaultErrorMsg)
	}

	created, err := svc.productRepo.GetByID(p.ID, false)
	if err != nil {
		return &p, http.StatusCreated, nil
	}
	return &created, http.StatusCreated, nil
}

func (svc *ProductService) Update(id uint, req dtos.UpdateProductRequest) (*models.Product, int, error) {
	// ensure exists
	existing, err := svc.productRepo.GetByID(id, true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.New(constant.NotFoundError)
		}
		return nil, http.StatusInternalServerError, errors.New(constant.DefaultErrorMsg)
	}

	updates := map[string]interface{}{}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, http.StatusBadRequest, errors.New("name cannot be empty")
		}
		updates["name"] = name
	}
	if req.SKU != nil {
		sku := strings.TrimSpace(*req.SKU)
		if sku == "" {
			return nil, http.StatusBadRequest, errors.New("sku cannot be empty")
		}
		exists, e := svc.productRepo.ExistsSKU(sku, id)
		if e != nil {
			return nil, http.StatusInternalServerError, errors.New(constant.DefaultErrorMsg)
		}
		if exists {
			return nil, http.StatusConflict, errors.New("sku already exists")
		}
		updates["sku"] = sku
	}
	if req.PurchasePrice != nil {
		updates["purchase_price"] = *req.PurchasePrice
	}
	if req.SellingPrice != nil {
		updates["selling_price"] = *req.SellingPrice
	}
	if req.Quantity != nil {
		if *req.Quantity < 0 {
			return nil, http.StatusBadRequest, errors.New("quantity cannot be negative")
		}
		updates["quantity"] = *req.Quantity
	}

	// validate price rule with final values
	finalPurchase := existing.PurchasePrice
	finalSelling := existing.SellingPrice
	if req.PurchasePrice != nil {
		finalPurchase = *req.PurchasePrice
	}
	if req.SellingPrice != nil {
		finalSelling = *req.SellingPrice
	}
	if finalSelling < finalPurchase {
		return nil, http.StatusBadRequest, errors.New("sellingPrice must be greater than or equal to purchasePrice")
	}

	if err := svc.productRepo.UpdateFields(id, updates); err != nil {
		return nil, http.StatusInternalServerError, errors.New(constant.DefaultErrorMsg)
	}

	updated, err := svc.productRepo.GetByID(id, true)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New(constant.DefaultErrorMsg)
	}
	return &updated, http.StatusOK, nil
}

func (svc *ProductService) List(q dtos.ListProductsQuery) ([]models.Product, int64, int, int, int, error) {
	page := q.Page
	limit := q.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	items, total, err := svc.productRepo.List(repository.ListProductsParams{
		Page:        page,
		Limit:       limit,
		Q:           q.Q,
		CategoryID:  q.CategoryID,
		MinPurchase: q.MinPurchase,
		MaxPurchase: q.MaxPurchase,
		MinSelling:  q.MinSelling,
		MaxSelling:  q.MaxSelling,
		SortBy:      q.SortBy,
		Order:       q.Order,
		IncludeDel:  q.IncludeDeleted,
	})
	if err != nil {
		return nil, 0, 0, 0, http.StatusInternalServerError, errors.New(constant.DefaultErrorMsg)
	}

	return items, total, page, limit, http.StatusOK, nil
}

func (svc *ProductService) GetByID(id uint, includeDeleted bool) (*models.Product, int, error) {
	p, err := svc.productRepo.GetByID(id, includeDeleted)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.New(constant.NotFoundError)
		}
		return nil, http.StatusInternalServerError, errors.New(constant.DefaultErrorMsg)
	}
	return &p, http.StatusOK, nil
}

func (svc *ProductService) SoftDelete(id uint) (int, error) {
	// ensure exists
	if _, err := svc.productRepo.GetByID(id, false); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New(constant.NotFoundError)
		}
		return http.StatusInternalServerError, errors.New(constant.DefaultErrorMsg)
	}

	if err := svc.productRepo.SoftDelete(id); err != nil {
		return http.StatusInternalServerError, errors.New(constant.DefaultErrorMsg)
	}
	return http.StatusOK, nil
}

func (svc *ProductService) Restore(id uint) (int, error) {
	// ensure exists even if deleted
	if _, err := svc.productRepo.GetByID(id, true); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New(constant.NotFoundError)
		}
		return http.StatusInternalServerError, errors.New(constant.DefaultErrorMsg)
	}

	if err := svc.productRepo.Restore(id); err != nil {
		return http.StatusInternalServerError, errors.New(constant.DefaultErrorMsg)
	}
	return http.StatusOK, nil
}
