package model

type UploadMessage struct {
	Key    string `json:"key"`
	Bucket string `json:"bucket"`
}

type PresignRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

type ProductImagePresignRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	ProductID   int64  `json:"product_id"`
}

type ProductImageConfirmRequest struct {
	Key       string `json:"key"`
	ProductID int64  `json:"product_id"`
}

type PresignResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}
