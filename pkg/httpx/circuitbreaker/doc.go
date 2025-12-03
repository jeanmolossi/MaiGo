// Package circuitbreaker provides a simple HTTP client circuit breaker.
//
// It exposes a RoundTripper middleware that can be composed with others to
// prevent calls to an upstream service when it is deemed unhealthy. The breaker
// transitions through three states:
//   - Closed: requests flow normally and failures are counted.
//   - Open: requests are short-circuited and fail immediately.
//   - Half-open: after a recovery window a single probe request is allowed. If
//     it succeeds the circuit closes, otherwise it reopens.
//
// Configuration is done through CircuitBreakerConfig:
//   - FailureThreshold: consecutive failures allowed before opening (default 5).
//   - RecoveryWindow: time the circuit stays open before a probe (default 30s).
//   - ShouldTrip: optional predicate to mark responses or errors as failures;
//     if nil, any error or HTTP status >=500 is considered a failure.
package circuitbreaker
