package validator

import "path/filepath"

func IsAllowedExtension(filename string) bool {
	ext := filepath.Ext(filename)

	if ext != ".xlsx" && ext != ".xls" && ext != ".pdf" {
		return false
	}

	return true
}

func IsAllowedImageExtension(filename string) bool {
	ext := filepath.Ext(filename)

	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return false
	}

	return true
}
