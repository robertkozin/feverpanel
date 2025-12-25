package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/kong"
	_ "github.com/robertkozin/feverpanel/backend/lib/tr"
	"go.opentelemetry.io/otel"
)

var (
	tracer = otel.Tracer("cli")
)

type CLI struct {
	Run RunCmd `cmd:"" help:"Runs Feverpanel"`
}

type CmdEnv struct {
	ctx    context.Context
	cwd    string
	stdout io.Writer
}

func main() {
	ctx := context.Background()
	cwd, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("getting current working directory: %w", err)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(exitCodeFromError(err))
	}

	if err = run(ctx, cwd, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(exitCodeFromError(err))
	}
}

func run(ctx context.Context, cwd string, stdout io.Writer, args []string) error {
	var cli CLI
	kongParser, err := kong.New(&cli,
		kong.Name("fever"),
		kong.Writers(stdout, stdout),
		kong.Bind(CmdEnv{ctx, cwd, stdout}),
	)
	if err != nil {
		return fmt.Errorf("constructing cli parser: %w", err)
	}
	cmd, err := kongParser.Parse(args[1:])
	if err != nil {
		// todo: replace with AsType in 1.26
		var parseErr *kong.ParseError
		if errors.As(err, &parseErr) {
			parseErr.Context.PrintUsage(true)
			fmt.Fprintln(stdout)
		}
		return err
	}

	return cmd.Run()
}

// https://github.com/square/exit?tab=readme-ov-file#about
func exitCodeFromError(err error) int {
	if err == nil {
		return 0
	} else if exitCodeErr, ok := err.(interface{ ExitCode() int }); ok {
		return exitCodeErr.ExitCode()
	}
	return 1
}
