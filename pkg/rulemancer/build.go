package rulemancer

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"text/template"
)

type ProtocolData struct {
	*Engine
	GameName            string              // Name of the game
	CurrentAssert       string              // The name of the current assert (used in assertions)
	CurrentAssertParams []string            // The union of the sets of assertions parameters (used in assertions)
	CurrentQuery        string              // The name of the current query (used in queries)
	Assertables         map[string][]string // Assertable of the game
	Responses           map[string][]string // Responses of the game
	Queryables          map[string][]string // Queryable of the game
	Relations           []string            // Relations of the game
	Slots               map[string][]string // Slots of the game
	Multislots          map[string][]string // Multislots of the game
	funcMap             template.FuncMap
}

func (e *Engine) newProtocolData(withFuncMap bool) *ProtocolData {
	slots := make(map[string][]string)
	multislots := make(map[string][]string)
	for _, game := range e.games {
		for _, relations := range game.assertable {
			for _, rel := range relations {
				slots[rel] = make([]string, 0)
				multislots[rel] = make([]string, 0)
			}
		}
		for _, relations := range game.responses {
			for _, rel := range relations {
				slots[rel] = make([]string, 0)
				multislots[rel] = make([]string, 0)
			}
		}
		for _, relations := range game.queryable {
			for _, rel := range relations {
				slots[rel] = make([]string, 0)
				multislots[rel] = make([]string, 0)
			}
		}
	}
	if !withFuncMap {
		return &ProtocolData{Engine: e, Slots: slots, Multislots: multislots, funcMap: nil}
	}

	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"dec": func(i int) int {
			return i - 1
		},
		"n": func(start, end int) []int {
			var result []int
			for i := start; i < end; i++ {
				result = append(result, i)
			}
			return result
		},
		"ns": func(start, end, step int) []int {
			var result []int
			for i := start; i < end; i += step {
				result = append(result, i)
			}
			return result
		},
		"sum": func(a, b int) int {
			return a + b
		},
		"div": func(a, b int) int {
			return a / b
		},
		"mult": func(a, b int) int {
			return a * b
		},
		"pow": func(a, b int) int {
			return int(math.Pow(float64(a), float64(b)))
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}
	return &ProtocolData{Engine: e, Slots: slots, Multislots: multislots, funcMap: funcMap}
}

// Merge merges another ProtocolData into this one.
func (pd *ProtocolData) Merge(other *ProtocolData) {
	for rel, slots := range other.Slots {
		if _, exists := pd.Slots[rel]; !exists {
			pd.Slots[rel] = slots
		} else {
			existingSlots := pd.Slots[rel]
			slotSet := make(map[string]struct{})
			for _, slot := range existingSlots {
				slotSet[slot] = struct{}{}
			}
			for _, slot := range slots {
				if _, found := slotSet[slot]; !found {
					existingSlots = append(existingSlots, slot)
				}
			}
			pd.Slots[rel] = existingSlots
		}
	}
	for rel, multislots := range other.Multislots {
		if _, exists := pd.Multislots[rel]; !exists {
			pd.Multislots[rel] = multislots
		} else {
			existingMultislots := pd.Multislots[rel]
			multislotSet := make(map[string]struct{})
			for _, multislot := range existingMultislots {
				multislotSet[multislot] = struct{}{}
			}
			for _, multislot := range multislots {
				if _, found := multislotSet[multislot]; !found {
					existingMultislots = append(existingMultislots, multislot)
				}
			}
			pd.Multislots[rel] = existingMultislots
		}
	}
}

func (e *Engine) BuildEngineExtras(shellOutdir string) error {
	// The rebuild engine reads the rules games directories and the assertables,results and querables from there.
	// Then uses re2c to write the various artifacts needed to interact with the engine.

	e.loadGames()

	// Define the templates directory
	templateDir := "pkg/rulemancer/templates"

	templateMap := make(map[string]string)

	// Load the shell templates
	if _, err := os.Stat(templateDir + "/shell"); os.IsNotExist(err) {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
			l.Printf("Shell template directory does not exist: %s", templateDir+"/shell")
		}
		return fmt.Errorf("template location does not exist: %s", templateDir+"/shell")
	}
	if templateFiles, err := os.ReadDir(templateDir + "/shell"); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
			l.Printf("Error reading shell template directory: %v", err)
		}
		return fmt.Errorf("failed to read rules location: %w", err)
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
					return fmt.Errorf("failed to read template file %s: %w", file.Name(), err)
				}
				templateMap[file.Name()] = string(content)
			}
		}
	}

	// Create the output directory if it doesn't exist
	if _, err := os.Stat(shellOutdir); os.IsNotExist(err) {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/BuildEngineExtras]")+" ", 0)
			l.Printf("Output directory does not exist, creating: %s", shellOutdir)
		}
		if err := os.MkdirAll(shellOutdir, 0755); err != nil {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
				l.Printf("Error creating output directory %s: %v", shellOutdir, err)
			}
			return fmt.Errorf("failed to create output directory %s: %w", shellOutdir, err)
		}
	}

	// Load each game relations information about slots and multislots
	gamesInterfaces := make(map[string]*ProtocolData)

	for _, game := range e.games {

		// Create a new ProtocolData for the game, it will hold the slots and multislots information from
		// the rules files. It also has the funcMap for template execution.
		rf := e.newProtocolData(true)

		gameName := game.name
		rulesLocation := game.rulesLocation
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/BuildEngineExtras]")+" ", 0)
			l.Printf("Building engine extras for game: %s from rules location: %s", gameName, rulesLocation)
		}
		// Load a game from the specified rules location
		if _, err := os.Stat(rulesLocation); os.IsNotExist(err) {
			return fmt.Errorf("rules location does not exist: %s", rulesLocation)
		}
		if rulesFiles, err := os.ReadDir(rulesLocation); err != nil {
			return fmt.Errorf("failed to read rules location: %w", err)
		} else {
			// Load each rule file into CLIPS
			for _, file := range rulesFiles {
				// load the file
				if !file.IsDir() {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/BuildEngineExtras]")+" ", 0)
						l.Printf("Processing rule file: %s", file.Name())
					}
					fileContent, err := os.ReadFile(rulesLocation + "/" + file.Name())
					if err != nil {
						if e.Debug {
							l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
							l.Printf("Error reading rule file %s: %v", file.Name(), err)
						}
						return fmt.Errorf("failed to read rule file %s: %w", file.Name(), err)
					}

					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/BuildEngineExtras]")+" ", 0)
						l.Printf("Executing parser on rule file: %s", file.Name())
					}

					pd := e.newProtocolData(false)

					if err := pd.Compile(string(fileContent) + "\x00"); err != nil {
						if e.Debug {
							l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
							l.Printf("Error compiling rule file %s: %v", file.Name(), err)
						}
						return fmt.Errorf("failed to compile rule file %s: %w", file.Name(), err)
					}
					rf.Merge(pd)

					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/BuildEngineExtras]")+" ", 0)
						l.Printf("Successfully processed rule file: %s", file.Name())
					}
				}
			}
		}

		rf.Assertables = game.assertable
		rf.Responses = game.responses
		rf.Queryables = game.queryable
		rf.GameName = gameName

		gamesInterfaces[gameName] = rf
	}

	// Now, we have all the games interfaces loaded with their slots and multislots. We also have the templates loaded.
	// Execute the templates for each game with the corresponding ProtocolData.

	for gameName, pd := range gamesInterfaces {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/BuildEngineExtras]")+" ", 0)
			l.Printf("Generating shell files for game: %s", gameName)
		}

		// Create the game directory inside the output directory if it doesn't exist
		gameOutdir := fmt.Sprintf("%s/%s", shellOutdir, gameName)
		if _, err := os.Stat(gameOutdir); os.IsNotExist(err) {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/BuildEngineExtras]")+" ", 0)
				l.Printf("Game output directory does not exist, creating: %s", gameOutdir)
			}
			if err := os.MkdirAll(gameOutdir, 0755); err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
					l.Printf("Error creating game output directory %s: %v", gameOutdir, err)
				}
				return fmt.Errorf("failed to create game output directory %s: %w", gameOutdir, err)
			}
		}

		// Execute each template
		for templateName, templateContent := range templateMap {
			tmpl, err := template.New(templateName).Funcs(pd.funcMap).Parse(templateContent)
			if err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
					l.Printf("Error parsing template %s for game %s: %v", templateName, gameName, err)
				}
				return fmt.Errorf("failed to parse template %s for game %s: %w", templateName, gameName, err)
			}

			// Process the template name to determine the type of file to generate.
			// For example, if the template name is "assert.sh", we might want to generate an assert shell script for every
			// assertable relation in the game. Same for "query.sh" and queryable relations. For other templates,
			// we might want to generate a single file with the template executed with the whole ProtocolData.
			switch templateName {
			case "assert.sh":
				for ass := range pd.Assertables {
					pd.CurrentAssert = ass
					pd.CurrentAssertParams = make([]string, 0)
					for _, rel := range pd.Assertables[ass] {
						if _, ok := pd.Slots[rel]; ok {
							for _, slot := range pd.Slots[rel] {
								if !isInSlice(pd.CurrentAssertParams, slot) {
									pd.CurrentAssertParams = append(pd.CurrentAssertParams, slot)
								}
							}
						}
					}
					for _, rel := range pd.Assertables[ass] {
						if _, ok := pd.Multislots[rel]; ok {
							for _, multislot := range pd.Multislots[rel] {
								if !isInSlice(pd.CurrentAssertParams, multislot) {
									pd.CurrentAssertParams = append(pd.CurrentAssertParams, multislot)
								}
							}
						}
					}
					outputFilePath := fmt.Sprintf("%s/%s/assert-%s.sh", shellOutdir, gameName, ass)
					if err := e.commitTemplate(tmpl, templateContent, outputFilePath, pd, gameName, templateName); err != nil {
						return err
					}
				}
			case "query.sh":
				for query := range pd.Queryables {
					pd.CurrentQuery = query
					outputFilePath := fmt.Sprintf("%s/%s/query-%s.sh", shellOutdir, gameName, query)
					if err := e.commitTemplate(tmpl, templateContent, outputFilePath, pd, gameName, templateName); err != nil {
						return err
					}
				}
			default:
				outputFilePath := fmt.Sprintf("%s/%s/%s", shellOutdir, gameName, templateName)
				if err := e.commitTemplate(tmpl, templateContent, outputFilePath, pd, gameName, templateName); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (e *Engine) commitTemplate(tmpl *template.Template, templateContent string, outputFilePath string, pd *ProtocolData, gameName string, templateName string) error {
	outputFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0700)
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
			l.Printf("Error creating output file %s: %v", outputFilePath, err)
		}
		return fmt.Errorf("failed to create output file %s: %w", outputFilePath, err)
	}
	defer outputFile.Close()
	if err := tmpl.Execute(outputFile, pd); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
			l.Printf("Error executing template %s for game %s: %v", templateName, gameName, err)
		}
		return fmt.Errorf("failed to execute template %s for game %s: %w", templateName, gameName, err)
	}
	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/BuildEngineExtras]")+" ", 0)
		l.Printf("Generated shell file: %s", outputFilePath)
	}
	return nil
}
