package rulemancer

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	chi "github.com/go-chi/chi/v5"
)

func (e *Engine) handlerGenerator(path string, templateMap map[string]string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(templateMap[path]))
	}
}

func (e *Engine) webClientRoutes(r chi.Router) {

	webPages, err := e.webPagesReady()
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/webClientRoutes]")+" ", 0)
			l.Printf("Failed to prepare web pages: %v", err)
		}
		return
	}
	fmt.Printf("Web pages ready: %v\n", webPages)
	r.Route("/", func(r chi.Router) {
		for path := range webPages {
			r.Get(path, e.handlerGenerator(path, webPages))
		}
	})
}

func (e *Engine) webPagesReady() (map[string]string, error) {
	templateDir := "pkg/rulemancer/templates/webclient"

	templateMap, err := e.webTemplates(templateDir, "")
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/webPagesReady]")+" ", 0)
			l.Printf("Failed to generate templates: %v", err)
		}
		return nil, fmt.Errorf("failed to generate templates: %w", err)
	}

	// Load each game relations information about slots and multislots
	gamesInterfaces, err := e.gameInterfaces()
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/webPagesReady]")+" ", 0)
			l.Printf("Error loading game interfaces: %v", err)
		}
		return nil, fmt.Errorf("failed to load game interfaces: %w", err)
	}

	webPages := make(map[string]string)

	for gameName, pd := range gamesInterfaces {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/webPagesReady]")+" ", 0)
			l.Printf("Generating shell files for game: %s", gameName)
		}

		// Execute each template
		for templateName, templateContent := range templateMap {
			tmpl, err := template.New(templateName).Funcs(pd.funcMap).Parse(templateContent)
			if err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/webPagesReady]")+" ", 0)
					l.Printf("Error parsing template %s for game %s: %v", templateName, gameName, err)
				}
				return nil, fmt.Errorf("failed to parse template %s for game %s: %w", templateName, gameName, err)
			}

			switch templateName {
			case "/client/game":
				var pageContent bytes.Buffer

				if err := tmpl.Execute(&pageContent, pd); err != nil {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/webPagesReady]")+" ", 0)
						l.Printf("Error executing template %s for game %s: %v", templateName, gameName, err)
					}
					return nil, fmt.Errorf("failed to execute template %s for game %s: %w", templateName, gameName, err)
				}
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/webPagesReady]")+" ", 0)
					l.Printf("Generated web page for template %s of game %s", templateName, gameName)
				}
				webPages["/client/"+gameName] = pageContent.String()
			default:
				if _, ok := webPages[templateName]; !ok {
					var pageContent bytes.Buffer

					if err := tmpl.Execute(&pageContent, pd); err != nil {
						if e.Debug {
							l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/webPagesReady]")+" ", 0)
							l.Printf("Error executing template %s for game %s: %v", templateName, gameName, err)
						}
						return nil, fmt.Errorf("failed to execute template %s for game %s: %w", templateName, gameName, err)
					}
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/webPagesReady]")+" ", 0)
						l.Printf("Generated web page for template %s of game %s", templateName, gameName)
					}
					webPages[templateName] = pageContent.String()
				}
			}
		}
	}
	return webPages, nil
}
