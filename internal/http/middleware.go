package http

import (
	"fmt"
	"github.com/tomwright/finance-planner/internal/errs"
	"net/http"
	"os"
	"runtime/debug"
)

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err := errs.New().
					WithCode(errs.ErrUnknown).
					WithStatusCode(http.StatusInternalServerError).
					WithMessage("Internal Server Error")

				_, _ = fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
				debug.PrintStack()

				sendError(err, w)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[HTTP] %s: %s\n", r.Method, r.RequestURI)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func CORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH, PUT, DELETE, OPTIONS, TRACE, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "X-Platform-UUID,X-User-UUID,X-Tenant-UUID,Content-Type,Content-Length,Accept,Authorization,If-Match,Accept-Encoding")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "ETag,Etag")
		w.Header().Set("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "-1")

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Options(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(nil)
			return
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
