/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package cmd

import (
	"log"

	"github.com/mmirko/rulemancer/pkg/rulemancer"
	"github.com/spf13/cobra"
)

var webBaseURL string

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server with a web interface for the client",
	Long:  `Start the web server with a web interface for the client.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Initialize the engine with the secret
		w := rulemancer.NewWebClient("")

		if cfgFile != "" {
			err := w.LoadConfig(cfgFile)
			if err != nil {
				log.Fatalf("Error loading config file: %v", err)
			}
		}
		if err := w.SpawnWebClient(webBaseURL); err != nil {
			log.Fatalf("Error starting web client: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().StringVarP(&webBaseURL, "baseurl", "b", "https://localhost:3000/api/v1", "Base URL for the engine API")
}
