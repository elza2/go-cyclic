package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

const Service = "github.com/elza2/go-cyclic"

var (
	rootCmd = &cobra.Command{Use: Service}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("dir", os.Getenv("dir"), "dir. eg: full path of the go.mod file")
}
