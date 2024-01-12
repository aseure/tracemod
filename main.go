package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	var (
		direction       Direction
		maxTraces       uint
		useFixedStrings bool
		timeout         time.Duration
	)

	rootCmd := &cobra.Command{
		Use:   "tracemod",
		Short: "Trace a module dependency from a Go project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rootModule, err := detectRootModule()
			if err != nil {
				return err
			}

			goModGraph, err := parseGoModGraph()
			if err != nil {
				return err
			}

			moduleFilter := args[0]

			isModuleMatching, err := buildModuleMatchingFunc(moduleFilter, useFixedStrings)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			chains, err := computeDependencyChains(ctx, goModGraph, rootModule, isModuleMatching)
			if err != nil {
				return err
			}

			diagram := generateMermaidDiagram(cancel, rootModule, isModuleMatching, direction, chains, maxTraces)

			filepath, err := generateHTML(diagram)
			if err != nil {
				return err
			}

			return open(filepath)
		},
	}

	rootCmd.PersistentFlags().UintVarP(&maxTraces, "max-traces", "m", 0, "Limit the number of maximum traces to detect")
	rootCmd.PersistentFlags().VarP(&direction, "direction", "d", `Direction of the dependency tree, defaults to "LR"`)
	rootCmd.PersistentFlags().BoolVarP(&useFixedStrings, "fixed-strings", "F", false, "Treat all patterns as literals instead of as regular expressions.")
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 30*time.Second, "Timeout duration")

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to execute: %v", err)
		os.Exit(1)
	}
}
