package actions

import (
	"github.com/spf13/cobra"
	"go-cyclic/tool"
	"log"
)

func RunCyclic(cmd *cobra.Command, args []string) {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		log.Fatalf("get dir params failed: %v\n", err)
	}
	if err = tool.CheckCycleDepend(dir); err != nil {
		log.Fatalf("run failed: %v\n", err)
	}
}

func init() {
	cmd := &cobra.Command{
		Use: "gocyclic",
		Run: RunCyclic,
	}
	rootCmd.AddCommand(cmd)
}
