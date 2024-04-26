package command

import (
	"context"
	"github.com/urfave/cli/v3"
)

func RunTestCases(ctx context.Context, cmd *cli.Command) error {
	//fmt.Println("removed task template: ", cmd.Args().First())
	return nil
}

func ValidateConfig(ctx context.Context, cmd *cli.Command) error {
	//fmt.Println("removed task template: ", cmd.Args().First())
	return nil
}
