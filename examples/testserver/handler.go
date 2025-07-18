package testserver

import (
	"encoding/binary"
	"encoding/json"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

type (
	Middleware func(http.Handler) http.Handler

	HelthcheckMessage struct {
		ServerID  string `json:"server_id"`
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
	}

	ErrorMessage struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
)

func NewMux(config *Server, repository *Provider) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Healthcheck(&config.ID))

	mux.HandleFunc("GET /users", GetAll(repository.User))
	mux.HandleFunc("GET /users/{id}", GetByID(repository.User))
	mux.HandleFunc("POST /users", Create[*User](repository.User))

	mux.HandleFunc("GET /resources", GetAll(repository.Resource))
	mux.HandleFunc("GET /resources/{id}", GetByID(repository.Resource))
	mux.HandleFunc("POST /resources", Create[*Resource](repository.Resource))

	middlewares := []Middleware{
		HeaderDebugMiddleware(config),
		BusyMiddleware(config),
		SleepMiddleware(config),
	}

	return handleWithMiddlewares(middlewares, mux)
}

func HeaderDebugMiddleware(config *Server) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.EnableHeaderDebug {
				slog.Info("request headers:", "headers", r.Header)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func BusyMiddleware(config *Server) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.EnableBusy && shouldSimulateServerError() {
				slog.Error("server is busy!", "serverID", config.ID)
				writeErrorResponse(w, ErrorMessage{
					Status:  http.StatusInternalServerError,
					Message: http.StatusText(http.StatusInternalServerError),
				})

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func SleepMiddleware(config *Server) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.Interval > time.Duration(0) {
				time.Sleep(calculateRandomDelay(config))
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Healthcheck(serverID *int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(HelthcheckMessage{
			ServerID:  strconv.Itoa(*serverID),
			Status:    "UP",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}
}

func GetAll(repository Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(repository.GetAll())
	}
}

func GetByID(repository Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 0 {
			writeErrorResponse(w, ErrorMessage{
				Status:  http.StatusBadRequest,
				Message: "invalid id",
			})

			return
		}

		resource, found := repository.GetByID(uint(id))
		if !found {
			writeErrorResponse(w, ErrorMessage{
				Status:  http.StatusNotFound,
				Message: "not found",
			})

			return
		}

		_ = json.NewEncoder(w).Encode(resource)
	}
}

func Create[T Model](repository Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var raw T
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			writeErrorResponse(w, ErrorMessage{
				Status:  http.StatusUnprocessableEntity,
				Message: "invalid request body",
			})

			return
		}

		newResource := repository.Create(raw)

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(newResource)
	}
}

func handleWithMiddlewares(middlewares []Middleware, mux http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		mux = middlewares[i](mux)
	}

	return mux
}

func writeErrorResponse(w http.ResponseWriter, errorMessage ErrorMessage) {
	w.WriteHeader(errorMessage.Status)
	_ = json.NewEncoder(w).Encode(errorMessage)
}

func shouldSimulateServerError() bool {
	return rand.IntN(100) < 50
}

func calculateRandomDelay(config *Server) time.Duration {
	interval := config.Interval
	rate := config.BackoffRate

	delay := float64(interval) * rate
	delay = secureFloat64() * delay

	return time.Duration(delay)
}

func secureFloat64() float64 {
	var seed [32]byte
	for i := range seed {
		seed[i] = byte(i)
	}

	c := rand.NewChaCha8(seed)

	bits := make([]byte, 8)

	// read randomic 8 bytes from secure source
	_, err := c.Read(bits)
	if err != nil {
		panic(err)
	}

	// parse bytes to 64 bits int
	n := binary.LittleEndian.Uint64(bits[:])

	// normalize to interval [0, 1]
	//nolint:mnd // 64 shift
	return float64(n) / (1 << 64)
}
