package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func detectRootModule() (ModuleName, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("could not read go.mod file: %w", err)
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return ModuleName(strings.TrimPrefix(line, "module ")), nil
		}
	}

	return "", errors.New("missing module name in go.mod file")
}

func parseGoModGraph() (map[ModuleName][]ModuleName, error) {
	output, err := exec.Command("go", "mod", "graph").Output()
	if err != nil {
		return nil, fmt.Errorf("could not run `go mod graph` command: %w", err)
	}

	modules := make(map[ModuleName][]ModuleName)

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		splits := strings.Split(line, " ")
		if len(splits) == 2 {
			module, dep := ModuleName(splits[0]), ModuleName(splits[1])
			modules[module] = append(modules[module], dep)
			modules[dep] = modules[dep]
		}
	}

	return modules, nil
}
