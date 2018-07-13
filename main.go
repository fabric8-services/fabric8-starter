package main

import (
	"fmt"
	"os"

	"github.com/fabric8-services/fabric8-starter/cmd"
)

func main() {

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
