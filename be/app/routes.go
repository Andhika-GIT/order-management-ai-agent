package app

import (
	"github.com/Andhika-GIT/go-message-broker-monorepo/controller"
	"github.com/go-chi/chi/v5"
)

func registerOrderRoutes(r chi.Router, c *controller.OrderController) {
	r.Get("/order", c.GetAllOrders)
	r.Post("/order/upload", c.UploadOrder)
}

func registerUserRoutes(r chi.Router, c *controller.UserController) {
	r.Get("/user", c.GetAllUsers)
	r.Post("/user/upload", c.UploadUser)
}

func registerDashboardRoutes(r chi.Router, c *controller.DashboardController) {
	r.Get("/dashboard", c.GetDataSummary)
}
