package model

type DashboardResponse struct {
	TotalOrders int64 `json:"total_orders"`
	TotalUsers  int64 `json:"total_users"`
}
