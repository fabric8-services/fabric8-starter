package cmd

import (
	"os"

	"github.com/fabric8-services/fabric8-starter/bootstrap"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Short:  "initializes a new fabric8-service project",
	Use:    "init",
	Args:   cobra.ExactArgs(1), // the name of the project to init
	PreRun: setLogLevel,
	Run: func(cmd *cobra.Command, args []string) {
		err := bootstrap.NewProject(args[0])
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}
