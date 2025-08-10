package logger

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/jeanmolossi/MaiGo/pkg/httpx"
)

const maxLogBodySize = 65536 // 64KiB

type LoggerHooks struct {
	Logger Logger

	OnStart func(ctx context.Context, r *http.Request)
	OnEnd   func(ctx context.Context, r *http.Request, res *http.Response, err error)

	ReqBodyTransformerFn func(ctx context.Context) func([]byte) []byte
	ResBodyTransformerFn func(ctx context.Context) func([]byte) []byte

	StartMessage string
	EndMessage   string

	LogStart      bool
	LogEnd        bool
	SupressErrors bool
}

func (h LoggerHooks) reqTx(ctx context.Context) func([]byte) []byte {
	if h.ReqBodyTransformerFn != nil {
		return h.ReqBodyTransformerFn(ctx)
	}

	// If has not a hook. We do not return with body. It's for security purposes.
	return func(b []byte) []byte { return make([]byte, 0) }
}

func (h LoggerHooks) resTx(ctx context.Context) func([]byte) []byte {
	if h.ResBodyTransformerFn != nil {
		return h.ResBodyTransformerFn(ctx)
	}

	// If has not a hook. We do not return with body. It's for security purposes.
	return func(b []byte) []byte { return make([]byte, 0) }
}

func LoggerRoundTripper(h LoggerHooks) httpx.ChainedRoundTripper {
	if h.Logger == nil {
		h.Logger = NewConsole()
	}

	if h.StartMessage == "" {
		const startMsg = "http.client.start"
		h.StartMessage = startMsg
	}

	if h.EndMessage == "" {
		const endMsg = "http.client.end"
		h.EndMessage = endMsg
	}

	return func(next http.RoundTripper) http.RoundTripper {
		return httpx.RoundTripperFn(func(r *http.Request) (*http.Response, error) {
			req := httpx.CloneRequest(r)
			ctx := req.Context()
			start := time.Now()

			if h.OnStart != nil {
				h.OnStart(ctx, req)
			}

			var rawReqBody []byte

			if req.Body != nil {
				raw, restored, err := httpx.ReadAndRestoreBody(req.Body, maxLogBodySize)
				if err == nil {
					req.Body = restored
					rawReqBody = h.reqTx(ctx)(raw)
				} else {
					if !h.SupressErrors {
						h.Logger.Error(ctx, err, "http.client.read_request_body_error")
					}
				}
			}

			attrs := make([]any, 0, 8)
			attrs = append(attrs,
				slog.String("method", req.Method),
				slog.String("url", req.URL.String()),
			)

			if h.LogStart {
				h.Logger.Info(ctx, h.StartMessage, attrs...)
			}

			resp, err := next.RoundTrip(req)
			elapsed := time.Since(start)

			var rawResBody []byte

			if err == nil && resp != nil && resp.Body != nil {
				raw, restored, rerr := httpx.ReadAndRestoreBody(resp.Body, maxLogBodySize)
				if rerr == nil {
					resp.Body = restored
					rawResBody = h.resTx(ctx)(raw)
				} else {
					if !h.SupressErrors {
						h.Logger.Error(ctx, rerr, "http.client.read_response_body_error")
					}
				}
			}

			var status, cl int
			if resp != nil {
				status = resp.StatusCode
				clHeader := resp.Header.Get("Content-Length")

				var serr error
				cl, serr = strconv.Atoi(clHeader)

				if serr != nil && !h.SupressErrors {
					h.Logger.Error(ctx, serr, "http.client.read_content_length_header_error")
				}
			}

			if h.OnEnd != nil {
				h.OnEnd(ctx, req, resp, err)
			}

			attrs = append(attrs,
				slog.Int("status", status),
				slog.Int("content_length", cl),
				slog.String("elapsed_ms", strconv.FormatInt(elapsed.Milliseconds(), 10)),
			)

			if len(rawReqBody) > 0 {
				attrs = append(attrs,
					slog.String("req_body", string(rawReqBody)),
					slog.String("req_body_size", strconv.Itoa(len(rawReqBody))+"B"),
				)
			}

			if len(rawResBody) > 0 {
				attrs = append(attrs,
					slog.String("res_body", string(rawResBody)),
					slog.String("res_body_size", strconv.Itoa(len(rawResBody))+"B"),
				)
			}

			if h.LogEnd {
				if err != nil {
					h.Logger.Error(ctx, err, h.EndMessage, attrs...)
				} else {
					h.Logger.Info(ctx, h.EndMessage, attrs...)
				}
			}

			return resp, err
		})
	}
}
