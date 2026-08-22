package app

import (
	"github.com/Andhika-GIT/go-message-broker-monorepo/controller"
	"github.com/go-chi/chi/v5"
)

func registerProductRoutes(r chi.Router, c *controller.ProductController) {
	r.Get("/product", c.GetAllProducts)
	r.Post("/product/upload", c.UploadProduct)
}

func registerDashboardRoutes(r chi.Router, c *controller.DashboardController) {
	r.Get("/dashboard", c.GetDataSummary)
}

func registerUploadRoutes(r chi.Router, c *controller.UploadController) {
	r.Post("/upload/presign", c.GetPresignedUploadURL)
	r.Post("/upload/product_image/presign", c.GetUploadProductImagePresignedURL)
	r.Post("/upload/product_image/confirm", c.ConfirmUploadProductImage)
}
