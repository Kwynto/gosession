package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
)

// Creating an `application` structure to store the dependencies of the entire web application.
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
}

// Structure for configuration
type Config struct {
	Addr      string
	StaticDir string
}

func main() {
	// Reading flags from the application launch bar.
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	// Creating loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize the template cache.
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize the structure with the application dependencies.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: templateCache,
	}

	// Server structure with address, logger and routing
	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Starting the server
	app.infoLog.Printf("Start server: %s", cfg.Addr)
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
