package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var verbose bool

// Execute runs the root command
func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "fabric8-starter",
		Short: "fabric8-starter is a CLI tool to initialize a fabric8-services project from scratch",
	}
	rootCmd.SetHelpCommand(helpCommand)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "set the level of logs verbose.") // flag for the log level
	rootCmd.AddCommand(initCmd)
	return rootCmd.Execute()
}

func setLogLevel(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
		log.Info("log level set to 'verbose'")
	}
}
