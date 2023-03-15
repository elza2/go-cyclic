package action

import (
	"github.com/urfave/cli/v2"

	"github.com/elza2/go-cyclic/core"
)

func CheckCyclic(ctx *cli.Context) error {
	dir := ctx.String("dir")
	filter := ctx.String("filter")

	return core.Do(dir, filter)
}
