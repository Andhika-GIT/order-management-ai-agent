package model

type UploadMessage struct {
	Filename string `json:"filename"`
	Filepath string `json:"filepath"`
}
