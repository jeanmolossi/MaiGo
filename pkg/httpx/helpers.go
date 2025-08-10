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

	if maxSize <= 0 {
		return nil, body, nil
	}

	maxSize = min(maxSize, MaxSafeBodyCap)

	buf := make([]byte, maxSize)
	n, rerr := io.ReadFull(body, buf)

	switch {
	case rerr == nil:
	case errors.Is(rerr, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF):
		buf = buf[:n]
	default:
		return nil, nil, rerr
	}

	// our readClose will restore original body with the read bytes and the rest of original body
	rc := &readClose{
		Reader: io.MultiReader(bytes.NewReader(buf), body),
		Closer: body, // closes original body
	}

	return buf, rc, nil
}

type readClose struct {
	io.Reader
	io.Closer
}
