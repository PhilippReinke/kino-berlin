package main

import (
	"flag"
	"log"
	"path/filepath"
	"time"
)

type Config struct {
	Addr         string
	SyncInterval time.Duration
	TemplateDir  string
	StaticDir    string
}

func parseFlags() Config {
	host := flag.String("host", "localhost", "Host")
	port := flag.String("port", "8080", "Port to listen on")
	syncInterval := flag.Duration("sync-interval", 30*time.Minute, "Background sync interval (0 to disable)")
	templateDir := flag.String("templates", "web/templates", "Template directory")
	staticDir := flag.String("static", "web/static", "Static files directory")
	flag.Parse()

	templateDirAbs, err := filepath.Abs(*templateDir)
	if err != nil {
		log.Printf("Warning: could not resolve template dir, using relative: %v", err)
		templateDirAbs = *templateDir
	}

	staticDirAbs, err := filepath.Abs(*staticDir)
	if err != nil {
		log.Printf("Warning: could not resolve static dir, using relative: %v", err)
		staticDirAbs = *staticDir
	}

	return Config{
		Addr:         *host + ":" + *port,
		SyncInterval: *syncInterval,
		TemplateDir:  templateDirAbs,
		StaticDir:    staticDirAbs,
	}
}
