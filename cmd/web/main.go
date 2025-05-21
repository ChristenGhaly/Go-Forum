package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"html/template"
	"time"

	"forum.christen.net/internal/models"

	"github.com/alexedwards/scs/v2"
	_ "modernc.org/sqlite"
)

type application struct {
	logger *slog.Logger
	users *models.UserModel
	threads *models.ThreadModel
	msgs *models.MessageModel
	templateCache map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dbPath := flag.String("db", "./forum.db", "Path to SQLite database")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dbPath)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	tempCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger: logger,
		users: &models.UserModel{DB: db},
		threads: &models.ThreadModel{DB: db},
		msgs: &models.MessageModel{DB: db},
		templateCache: tempCache,
		sessionManager: sessionManager,
	}

	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
