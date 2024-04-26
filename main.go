package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"github.com/vearne/autotest/internal/command"
	"log"
	"os"
)

const (
	version = "v0.0.1"
)

func main() {
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Printf("version=%s\n", cmd.Root().Version)
	}

	cmd := &cli.Command{
		Name:      "autotest",
		Version:   version,
		Usage:     "automate test",
		Authors:   []any{"vearne"},
		Copyright: "2024 vearne",
		Commands: []*cli.Command{
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "validate configuration files",
				Action:  command.RunTestCases,
			},
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "run all test cases",
				Action:  command.ValidateConfig,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
