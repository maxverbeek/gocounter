package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/maxverbeek/gocounter/pkg/lexer"
)

type CountResult struct {
	Count int
	Error error
}

func run(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("Expected 1 argument. Usage: %s <file-or-directory>", args[0])
	}

	files, err := lexer.ListGoFiles(args[1])

	if err != nil {
		return err
	}

	sum := 0
	var error error

	for _, f := range files {
		count, err := lexer.CountTokens(f)

		if err != nil {
			error = errors.Join(fmt.Errorf("failed to count tokens: %w", err), error)
			continue
		}

		sum += count
	}

	if error != nil {
		return error
	}

	fmt.Printf("%d\n", sum)

	return nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelWarn,
	})))

	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
	}
}
