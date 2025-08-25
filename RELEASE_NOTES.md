# Release Notes

## v1.2.7

### BREAKING CHANGES
- Added `Len() int` method to `contracts.Cookies` interface

### Features
- Harden cookies with nil, RFC-compliant name validation, bounds checks, and deep copy semantics (Add and Get clone cookies and normalize stored Names)
- Add unit tests and benchmarks for cookie helpers
- Fix HTTP initialism in cookie helper name
- Introduce Len alias and deprecate Count in cookie interface

## v1.2.6

- ensure `Context.Set` safely ignores nil, avoids redundant assignments, and compares pointer identities to prevent panics with uncomparable contexts
- document the context wrapper's nil-safe semantics and lack of concurrent-safety, clarifying `Unwrap`'s use of `context.Background`
- expand tests and benchmarks for `Context` to cover nil handling, idempotent `Unwrap`, and uncomparable contexts

## v1.2.5

- fix data race in UnbufferedBody by using write locks
- add tests for body read/write helpers and concurrency safety
- lint-only changes
