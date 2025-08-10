package retry

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type BodyReplayStrategy int

const (
	ReplayIfSmallElseNoRetry BodyReplayStrategy = iota // secure default
	ReplayIfSmallElseSpillToDisk
	NoReplay
)

func ensureReopenableBody(r *http.Request, capMem int64, strat BodyReplayStrategy) (cleanup func(), ok bool, err error) {
	if r.Body == nil || r.GetBody != nil {
		return func() {}, true, nil
	}

	if strat == NoReplay {
		return func() {}, false, nil
	}

	if capMem <= 0 {
		capMem = 1 << 20 // 1MiB default
	}

	// try in memo
	var buf bytes.Buffer

	n, err := io.CopyN(&buf, r.Body, capMem+1) // read +1 to detect overflow
	if cerr := r.Body.Close(); err == nil {
		err = cerr
	}

	if err != nil && !errors.Is(err, io.EOF) {
		return nil, false, err
	}

	if n <= capMem {
		_ = r.Body.Close()
		b := buf.Bytes()
		r.Body = io.NopCloser(bytes.NewReader(b))
		r.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(b)), nil }
		r.ContentLength = int64(len(b))

		return func() {}, true, nil
	}

	// If reach here. Size was exceeded.
	// Decide by policy
	if strat != ReplayIfSmallElseSpillToDisk {
		// No spill, caller should avoid retry when have body
		_ = r.Body.Close()
		return func() {}, false, nil
	}

	// Big: spill to file

	f, ferr := os.CreateTemp(".spill", "retry-body-*")
	if ferr != nil {
		_ = r.Body.Close()
		return nil, false, ferr
	}

	cleanup = func() { _ = os.Remove(f.Name()) }

	//nolint:wsl
	if _, err = f.Write(buf.Bytes()); err != nil {
		_ = f.Close()
		cleanup()
		_ = r.Body.Close()

		return nil, false, err
	}

	//nolint:wsl
	if _, err = io.Copy(f, r.Body); err != nil {
		_ = f.Close()
		cleanup()
		_ = r.Body.Close()

		return nil, false, err
	}

	if cerr := r.Body.Close(); cerr != nil {
		err = cerr
	}

	//nolint:wsl
	if err != nil {
		_ = f.Close()
		cleanup()

		return nil, false, err
	}

	// restart at start and plug in as initial body
	//nolint:wsl
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		_ = f.Close()
		cleanup()

		return nil, false, err
	}

	r.Body = f

	name := f.Name()
	r.GetBody = func() (io.ReadCloser, error) {
		return os.Open(name)
	}

	if stat, serr := os.Stat(name); serr == nil {
		r.ContentLength = stat.Size()
	}

	return cleanup, true, nil
}

func parseRetryAfter(h string) (time.Duration, bool) {
	if h == "" {
		return 0, false
	}

	if secs, err := strconv.Atoi(h); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second, true
	}

	if t, err := http.ParseTime(h); err == nil {
		if d := time.Until(t); d > 0 {
			return d, true
		}
	}

	return 0, false
}
