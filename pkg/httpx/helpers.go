package httpx

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

var ErrBodyTruncated = errors.New("body truncated to maxSize")

const (
	MaxSafeBodyCap = 1 << 30 // 1 MiB
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

	if maxSize <= 0 {
		return nil, http.NoBody, nil
	}

	maxSize = min(maxSize, MaxSafeBodyCap)

	buf := make([]byte, maxSize+1)
	n, readErr := io.ReadAtLeast(body, buf[:0], 0)

	if n == 0 && readErr == io.ErrUnexpectedEOF {
		readErr = nil
	}

	if readErr != nil && !errors.Is(readErr, io.EOF) {
		return nil, nil, readErr
	}

	n = min(n, maxSize)

	b := buf[:n]

	rc := io.NopCloser(bytes.NewReader(b))

	if n > maxSize {
		return b, rc, ErrBodyTruncated
	}

	return b, rc, nil
}
