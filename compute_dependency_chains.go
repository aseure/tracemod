package main

import (
	"context"
	"slices"
	"sync"
)

func computeDependencyChains(
	ctx context.Context,
	goModGraph map[ModuleName][]ModuleName,
	rootModule ModuleName,
	isModuleMatching ModuleMatchingFunc,
) (<-chan DependencyChain, error) {
	var (
		chains = make(chan DependencyChain)
		wg     = new(sync.WaitGroup)
	)

	moduleToParents := make(map[ModuleName][]ModuleName)
	for mod, deps := range goModGraph {
		for _, dep := range deps {
			moduleToParents[dep] = append(moduleToParents[dep], mod)
		}
	}

	for module, parents := range moduleToParents {
		if isModuleMatching(module) {
			for _, parent := range parents {
				wg.Add(1)
				go computeDependencyChain(ctx, wg, moduleToParents, rootModule, parent, NewDependencyChain(module), chains)
			}
		}
	}

	go func() {
		wg.Wait()
		close(chains)
	}()

	return chains, nil
}

func computeDependencyChain(
	ctx context.Context,
	wg *sync.WaitGroup,
	moduleToParents map[ModuleName][]ModuleName,
	rootModule ModuleName,
	module ModuleName,
	chain DependencyChain,
	chains chan<- DependencyChain,
) {
	defer wg.Done()

	if isContextDone(ctx) {
		return
	}

	// Prevent circular dependencies
	if chain.Has(module) {
		return
	}

	chain = chain.Add(module)

	if module == rootModule {
		slices.Reverse(chain.chain)
		chains <- chain
		return
	}

	for _, parent := range moduleToParents[module] {
		wg.Add(1)
		go computeDependencyChain(ctx, wg, moduleToParents, rootModule, parent, chain, chains)
	}
}
