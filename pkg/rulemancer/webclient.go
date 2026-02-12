package rulemancer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebClient struct {
	*Config
	router   chi.Router
	stopChan chan os.Signal
}

func NewWebClient(secret string) *WebClient {
	return &WebClient{
		Config:   NewConfig(),
		router:   chi.NewRouter(),
		stopChan: make(chan os.Signal, 1),
	}
}

func (e *WebClient) templateGenerator(baseDir, templateDir string) (map[string]string, error) {
	templateMap := make(map[string]string)

	if templateFiles, err := os.ReadDir(baseDir + "/" + templateDir); err == nil {
		// Process each template file
		for _, file := range templateFiles {
			if !file.IsDir() {
				content, err := os.ReadFile(baseDir + "/" + templateDir + "/" + file.Name())
				if err != nil {
					return nil, err
				}
				templateMap[templateDir+"/"+file.Name()] = string(content)
			} else {
				subTemplateMap, err := e.templateGenerator(baseDir, templateDir+"/"+file.Name())
				if err != nil {
					return nil, err
				}
				for k, v := range subTemplateMap {
					templateMap[k] = v
				}
			}
		}
	}

	return templateMap, nil
}

func (e *WebClient) handlerGenerator(path string, templateMap map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(templateMap[path]))
	}
}

func (e *WebClient) SpawnWebClient(baseURL string) error {

	// e.loadGames()

	r := e.router
	c := e.Config

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	if e.Debug {
		r.Use(middleware.Logger)
	}

	templateDir := "pkg/rulemancer/templates/webclient"

	templateMap, err := e.templateGenerator(templateDir, "")
	if err != nil {
		return fmt.Errorf("failed to generate templates: %w", err)
	}

	fmt.Println(templateMap)

	r.Route("/", func(r chi.Router) {
		for path := range templateMap {
			r.Get(path, e.handlerGenerator(path, templateMap))
		}
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServeTLS(c.TLSCertFile, c.TLSKeyFile); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v\n", err)
		}
	}()

	signal.Notify(e.stopChan, os.Interrupt, syscall.SIGTERM)
	<-e.stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")

	return nil
}
