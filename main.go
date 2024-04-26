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
				Name: "test",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "config-file", Aliases: []string{"c"}},
				},
				Usage:  "validate configuration files",
				Action: command.RunTestCases,
			},
			{
				Name: "run",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "config-file", Aliases: []string{"c"}},
					&cli.StringFlag{Name: "env-file", Aliases: []string{"e"}},
				},
				Usage:  "run all test cases",
				Action: command.ValidateConfig,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
