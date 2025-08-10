// Package logger provides middleware for logging HTTP client requests and
// responses. It exposes a Logger interface and a round tripper factory that can
// be composed with other middlewares to trace requests, responses and elapsed
// time. Hooks allow transforming or redacting bodies before they are logged.
package logger
