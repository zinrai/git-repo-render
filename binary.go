package main

import (
	"net/http"
	"strings"
)

// isBinaryContent reports whether data looks like binary content.
// It uses http.DetectContentType on up to the first 512 bytes.
func isBinaryContent(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	sample := data
	if len(sample) > 512 {
		sample = sample[:512]
	}
	contentType := http.DetectContentType(sample)
	return !strings.HasPrefix(contentType, "text/") && contentType != "application/json"
}
