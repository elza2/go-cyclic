package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Run: RunCyclic,
	}

	rootCmd.PersistentFlags().String("dir", os.Getenv("dir"), "dir. eg: directory address of the go.mod file")
	rootCmd.PersistentFlags().String("filter", os.Getenv("filter"), "dir. eg: filters out files of the specified type. separate multiple types with commas (,)")

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
