## v1.2.7

- harden Cookies with nil and bounds checks and copy semantics
- add unit tests and benchmarks for cookie helpers
- fix HTTP initialism in cookie helper name

## v1.2.6

- ensure `Context.Set` safely ignores nil, avoids redundant assignments, and compares pointer identities to prevent panics with uncomparable contexts
- document the context wrapper's nil-safe semantics and lack of concurrent-safety, clarifying `Unwrap`'s use of `context.Background`
- expand tests and benchmarks for `Context` to cover nil handling, idempotent `Unwrap`, and uncomparable contexts

## v1.2.5

- fix data race in UnbufferedBody by using write locks
- add tests for body read/write helpers and concurrency safety
- lint-only changes
