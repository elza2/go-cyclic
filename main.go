package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"

	"github.com/elza2/go-cyclic/action"
)

var (
	version     = "v1.1.0"
	serviceName = "go-cyclic"
	commands    = []*cli.Command{
		{
			Name:   "run",
			Usage:  "run tool for check circular references",
			Action: action.CheckCyclic,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "dir",
					Usage: "destination file address, must be in the same directory as the go.mod file",
				},
				&cli.StringFlag{
					Name:  "filter",
					Usage: "filter the specified file, the wildcard mode is supported, eg: *_test.go",
				},
			},
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Name = serviceName
	app.Usage = "a tool to check the go application for circular references"
	app.Version = fmt.Sprintf("%s %s/%s", version, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("go-cyclic: %+v\n", err)
	}
}
