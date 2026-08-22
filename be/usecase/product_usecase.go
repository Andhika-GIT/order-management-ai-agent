package usecase

import (
	"context"
	"fmt"

	"github.com/Andhika-GIT/go-message-broker-monorepo/model"
	"github.com/Andhika-GIT/go-message-broker-monorepo/repository"
)

type ProductUseCase struct {
	Repository *repository.ProductRepository
}

func NewProductUseCase(Repository *repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		Repository: Repository,
	}
}

func (u *ProductUseCase) CountAllProducts(c context.Context) (*int64, error) {
	totalProducts, err := u.Repository.CountAll(c)

	if err != nil {
		return nil, model.WriteError(500, fmt.Sprintf("failed to count all products %s", err.Error()))
	}

	return totalProducts, nil
}

func (u *ProductUseCase) FindAllProducts(c context.Context, paginationReq *model.PaginationRequest, filter *model.ProductFilter) (*model.Paginated[model.ProductResponse], error) {
	paginated, err := u.Repository.FindAll(c, paginationReq, filter)

	if err != nil {
		return nil, model.WriteError(500, fmt.Sprintf("failed to find all products %s", err.Error()))
	}

	formatedProducts := convertToProductsResponse(paginated.Data)

	// return new paginated response with different type (ProductResponse)
	return &model.Paginated[model.ProductResponse]{
		Data:       formatedProducts,
		Total:      paginated.Total,
		TotalPages: paginated.TotalPages,
	}, nil

}

func (u *ProductUseCase) FindProductByID(c context.Context, productID int64) (*model.ProductResponse, error) {
	product, err := u.Repository.FindById(c, productID)

	if err != nil {
		return nil, err
	}

	resp := convertToProductResponse(*product)

	return &resp, nil

}

func (u *ProductUseCase) InsertProductImageURL(c context.Context, productID int64, image_url string) error {
	return u.Repository.UpdateProductImageURL(c, productID, image_url)
}

func convertToProductsResponse(products []model.Product) []model.ProductResponse {
	var productsResp []model.ProductResponse

	for _, product := range products {
		productsResp = append(productsResp, convertToProductResponse(product))
	}

	return productsResp
}

func convertToProductResponse(product model.Product) model.ProductResponse {

	return model.ProductResponse{
		ID:            product.ID,
		SKU:           product.SKU,
		Name:          product.Name,
		Slug:          product.Slug,
		Description:   product.Description,
		Price:         product.Price,
		DiscountPrice: product.DiscountPrice,
		Stock:         product.Stock,
		Weight:        product.Weight,
		ImageURL:      product.ImageURL,
		Status:        product.Status,
		IsFeatured:    product.IsFeatured,
	}
}
