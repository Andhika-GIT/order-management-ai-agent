package model

type UploadMessage struct {
	Key    string `json:"key"`
	Bucket string `json:"bucket"`
}

type PresignRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

type OrderImagePresignRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	OrderID     int64  `json:"order_id"`
}

type OrderAttachmentConfirmRequest struct {
	Key     string `json:"key"`
	OrderID int64  `json:"order_id"`
}

type PresignResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}
