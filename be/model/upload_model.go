package model

type UploadMessage struct {
	Filename string `json:"filename"`
	Filepath string `json:"filepath"`
}

type PresignRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

type PresignResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}
