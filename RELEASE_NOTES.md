# Release Notes

## v1.2.17

- Added `RequestBuilder.Unwrap` to emit fully configured `*http.Request` values for direct `net/http` compatibility without verb-specific helpers.
- Retained compatibility helpers for unwrapping configured `*http.Client` instances while keeping existing builder APIs intact.
- Documented standard library interop with updated examples showcasing `DefaultClientCompat` and request unwrapping.

## v1.2.16

- Added OpenTelemetry tracing round tripper for outbound HTTP requests with context propagation and span lifecycle management.
- Documented tracing requirements, minimal tool versions, and provided an example of OTEL integration.
- Updated examples to leverage the new tracing middleware and refreshed linting prerequisites.

## v1.2.15

- Added HTTP client circuit breaker middleware with configurable failure thresholds, recovery windows, and trip predicate support.
- Documented circuit breaker behaviour, configuration options, and state transitions.
- Added tests covering circuit opening, half-open probing, and recovery paths.

## v1.2.14

- add a metrics-enabled HTTP round tripper that records per-method/status counters and histograms for client requests
- multple lint issues fixed

## v1.2.7

### BREAKING CHANGES

- Added `Len() int` method to `contracts.Cookies` interface

### Features

- Harden cookies with nil, RFC-compliant name validation, bounds checks, and deep copy semantics (Add and Get clone cookies and normalize stored Names)
- Add unit tests and benchmarks for cookie helpers
- Fix HTTP initialism in cookie helper name
- Introduce Len alias and deprecate Count in cookie interface
- Optimize cookie name validation with a lookup table to reduce branching

## v1.2.6

- ensure `Context.Set` safely ignores nil, avoids redundant assignments, and compares pointer identities to prevent panics with uncomparable contexts
- document the context wrapper's nil-safe semantics and lack of concurrent-safety, clarifying `Unwrap`'s use of `context.Background`
- expand tests and benchmarks for `Context` to cover nil handling, idempotent `Unwrap`, and uncomparable contexts

## v1.2.5

- fix data race in UnbufferedBody by using write locks
- add tests for body read/write helpers and concurrency safety
- lint-only changes
