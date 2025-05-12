package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"

	dotenv "github.com/joho/godotenv"
)

var (
	addr string
	dsn  string
)

func init() {
	if err := dotenv.Load(); err != nil {
		panic(err)
	}

	if dsn = os.Getenv("DSN"); dsn == "" {
		panic("must set $DSN")
	}

	if addr = os.Getenv("ADDR"); addr == "" {
		addr = ":8080"
	}
}

func home(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./web/html/pages/index.html")
		if err != nil {
			logger.Error("failed to parse template", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := t.Execute(w, nil); err != nil {
			logger.Error("failed to execute temlpate", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()

	mux.Handle("GET /css/", http.StripPrefix("/css", http.FileServer(http.Dir("./web/css"))))
	mux.Handle("GET /js/", http.StripPrefix("/js", http.FileServer(http.Dir("./web/js"))))
	mux.Handle("GET /img/", http.StripPrefix("/img", http.FileServer(http.Dir("./web/img"))))

	mux.Handle("GET /{$}", home(logger))

	s := http.Server{
		Addr:     addr,
		Handler:  mux,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", s.Addr)
	logger.Error("failed to start server", "err", s.ListenAndServe().Error())
}
