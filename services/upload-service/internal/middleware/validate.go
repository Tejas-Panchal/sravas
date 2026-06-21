package middleware

import (
	"mime/multipart"
	"path/filepath"
	"strings"
)

var allowedExtensions = map[string]bool{
	".mp4":  true,
	".webm": true,
	".avi":  true,
}

var allowedMIMETypes = map[string]bool{
	"video/mp4":        true,
	"video/webm":       true,
	"video/avi":        true,
	"video/x-msvideo":  true,
}

// MaxUploadSize is the maximum allowed file size (2 GB)
const MaxUploadSize int64 = 2 * 1024 * 1024 * 1024

// ValidateFile checks the file extension, MIME type, and size
func ValidateFile(header *multipart.FileHeader) string {
	if header.Size > MaxUploadSize {
		return "file exceeds maximum size of 2GB"
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		return "file format not supported (allowed: mp4, webm, avi)"
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedMIMETypes[contentType] {
		return "unsupported media type"
	}

	return ""
}
