package cmd

import (
	"log"
	"strings"

	"github.com/elza2/go-cyclic/errors"
	"github.com/elza2/go-cyclic/tool"
	"github.com/spf13/cobra"
)

const (
	GoSuffix = ".go"
)

func RunCyclic(cmd *cobra.Command, args []string) {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		log.Fatalf("get dir params failed: %v\n", err)
	}
	filters := make([]string, 0)
	filter, err := cmd.Flags().GetString("filter")
	if filter != "" {
		filters, err = HandleFilters(filter)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	if err = tool.CheckCycleDepend(&tool.Params{
		Dir:     dir,
		Filters: filters,
	}); err != nil {
		log.Fatalf("run failed. %v\n", err)
	}
}

func HandleFilters(filter string) (filters []string, err error) {
	if strings.Contains(filter, "，") {
		return nil, errors.NotSupportCNComma()
	}
	filters = strings.Split(filter, ",")
	for i, f := range filters {
		if strings.Contains(f, GoSuffix) {
			continue
		}
		filters[i] += GoSuffix
	}
	return filters, nil
}

func init() {
	cmd := &cobra.Command{
		Use: "gocyclic",
		Run: RunCyclic,
	}
	rootCmd.AddCommand(cmd)
}
