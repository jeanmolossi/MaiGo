package httpx

import (
	"bytes"
	"io"
	"net/http"
)

// CloneRequest avoid data races (copies r and Header)
func CloneRequest(r *http.Request) *http.Request {
	cr := r.Clone(r.Context())
	*cr.URL = *r.URL
	cr.Header = r.Header.Clone()

	return cr
}

func ReadAndRestoreBody(body io.ReadCloser, maxSize int) (raw []byte, restored io.ReadCloser, err error) {
	if body == nil {
		return nil, http.NoBody, nil
	}

	//nolint:errcheck // just close body
	defer body.Close()

	limited := io.LimitReader(body, int64(maxSize)+1) // +1 to mark truncate

	b, err := io.ReadAll(limited)
	if err != nil {
		return nil, nil, err
	}

	if len(b) > maxSize {
		b = b[:maxSize]
	}

	return b, io.NopCloser(bytes.NewReader(b)), nil
}
