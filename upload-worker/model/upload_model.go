package model

type UploadMessage struct {
	Key    string `json:"key"`
	Bucket string `json:"bucket"`
}
