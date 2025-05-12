package main

import (
	"database/sql"
	"errors"
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

		if err := t.ExecuteTemplate(w, "index", nil); err != nil {
			logger.Error("failed to execute temlpate", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

type ConversionItem struct {
	Name                  string
	BomID                 string
	ConversionFactor      float32
	DefaultQuantity       int8
	DestinationItemNumber string
	SourceItemNumber      string
	DestinationUnit       string
}

func getItem(logger *slog.Logger, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		row := db.QueryRow("select * from item_convert where source_item_number=$1", r.URL.Query().Get("item"))

		var item ConversionItem

		err := row.Scan(
			&item.Name,
			&item.BomID,
			&item.ConversionFactor,
			&item.DefaultQuantity,
			&item.DestinationItemNumber,
			&item.SourceItemNumber,
			&item.DestinationUnit,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				logger.Error("failed to scan row", "item", r.URL.Query().Get("item"), "err", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		t, err := template.ParseFiles("./web/html/pages/index.html")
		if err != nil {
			logger.Error("failed to parse template", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var pd struct {
			Form  map[string]string
			Error map[string]string
			Item  ConversionItem
		}
		pd.Form["item"] = item.SourceItemNumber
		pd.Form["quantity"] = "1"
		pd.Item = item

		if err := t.ExecuteTemplate(w, "form", pd); err != nil {
			logger.Error("failed to execute template", "err", err)
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
