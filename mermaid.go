package main

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

func generateMermaidDiagram(
	ctxCancel func(),
	rootModule ModuleName,
	isModuleMatching ModuleMatchingFunc,
	direction Direction,
	chains <-chan DependencyChain,
	maxTraces uint,
) string {
	var (
		b strings.Builder
	)

	write := func(format string, a ...interface{}) {
		_, _ = fmt.Fprintf(&b, format+"\n", a...)
	}

	mermaidConfig := map[string]any{
		"flowchart": map[string]any{
			"useMaxWidth": 0,
			"curve":       "linear",
			"htmlLabels":  false,
		},
		"theme": "base",
		"themeVariables": map[string]any{
			"primaryColor":       "#3e4042",
			"primaryBorderColor": "#6e7072",
			"primaryTextColor":   "#f0f1f2",
			"lineColor":          "#50b8e0",
			"secondaryColor":     "#6e7072",
			"tertiaryColor":      "#6e7072",
		},
	}

	links, maxDepth := generateMermaidLinks(ctxCancel, chains, maxTraces)
	ids := generateMermaidIDs(links)

	data, _ := json.MarshalIndent(mermaidConfig, "", "  ")
	write("%%%%{\ninit")
	write(string(data))
	write("}%%%%")

	if direction == "" {
		direction = computeBestDirection(maxDepth)
	}
	write("graph %s", direction)

	for module, id := range ids {
		write("%d(%q)", id, module)
	}

	for module, deps := range links {
		for _, dep := range deps {
			moduleID, depID := ids[module], ids[dep]
			write("%d --> %d", moduleID, depID)
		}
	}

	write("style %d stroke:#50b8e0,color:#50b8e0;", ids[rootModule])
	for module, id := range ids {
		if isModuleMatching(module) {
			write("style %d fill:#fceea5,stroke:#807956,color:#202224;", id)
		}
	}

	return b.String()
}

func generateMermaidLinks(
	ctxCancel func(),
	chains <-chan DependencyChain,
	maxTraces uint,
) (map[ModuleName][]ModuleName, int) {
	var (
		links         = make(map[ModuleName][]ModuleName)
		maxDepth      = 0
		nbTraces uint = 0
	)

	for chain := range chains {
		modules := chain.Deps()

		if len(modules) > maxDepth {
			maxDepth = len(modules)
		}

		for i := 0; i < len(modules)-1; i++ {
			module, dep := modules[i], modules[i+1]
			deps := links[module]
			if !slices.Contains(deps, dep) {
				links[module] = append(deps, dep)
			}
		}

		nbTraces++
		if nbTraces == maxTraces {
			ctxCancel()
			break
		}
	}

	return links, maxDepth
}

func generateMermaidIDs(links map[ModuleName][]ModuleName) map[ModuleName]int {
	uniqueModules := make(map[ModuleName]bool)
	for module, deps := range links {
		uniqueModules[module] = true
		for _, dep := range deps {
			uniqueModules[dep] = true
		}
	}

	labels := make(map[ModuleName]int)
	id := 0
	for module := range uniqueModules {
		labels[module] = id
		id++
	}

	return labels
}

func computeBestDirection(maxDepth int) Direction {
	if maxDepth > 5 {
		return "LR"
	}
	return "TB"
}
