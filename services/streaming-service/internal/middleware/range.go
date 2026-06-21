package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// Range represents a parsed HTTP Range header (bytes=start-end)
type Range struct {
	Start int64
	End   int64
}

// ParseRange extracts the byte range from the Range header (supports single ranges only)
func ParseRange(r *http.Request, fileSize int64) *Range {
	header := r.Header.Get("Range")
	if header == "" {
		return nil
	}

	if !strings.HasPrefix(header, "bytes=") {
		return nil
	}

	parts := strings.SplitN(header[6:], "-", 2)
	if len(parts) != 2 {
		return nil
	}

	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil
	}

	if parts[1] == "" {
		return &Range{Start: start, End: fileSize - 1}
	}

	end, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil
	}

	return &Range{Start: start, End: end}
}
