package rulemancer

import (
	"fmt"
	"log"
	"os"
)

func (e *Engine) webTemplates(baseDir, templateDir string) (map[string]string, error) {
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
				subTemplateMap, err := e.webTemplates(baseDir, templateDir+"/"+file.Name())
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

func (e *Engine) shellTemplates(templateDir string) (map[string]string, error) {
	templateMap := make(map[string]string)
	// Load the shell templates
	if _, err := os.Stat(templateDir + "/shell"); os.IsNotExist(err) {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
			l.Printf("Shell template directory does not exist: %s", templateDir+"/shell")
		}
		return nil, fmt.Errorf("template location does not exist: %s", templateDir+"/shell")
	}
	if templateFiles, err := os.ReadDir(templateDir + "/shell"); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
			l.Printf("Error reading shell template directory: %v", err)
		}
		return nil, fmt.Errorf("failed to read rules location: %w", err)
	} else {
		// Process each template file
		for _, file := range templateFiles {
			if !file.IsDir() {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/BuildEngineExtras]")+" ", 0)
					l.Printf("Processing shell template file: %s", file.Name())
				}
				content, err := os.ReadFile(templateDir + "/shell/" + file.Name())
				if err != nil {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
						l.Printf("Error reading template file %s: %v", file.Name(), err)
					}
					return nil, fmt.Errorf("failed to read template file %s: %w", file.Name(), err)
				}
				templateMap[file.Name()] = string(content)
			}
		}
	}

	return templateMap, nil
}
