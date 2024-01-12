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

	for module, deps := range goModGraph {
		for _, dep := range deps {
			if isModuleMatching(dep) {
				wg.Add(1)
				go computeDependencyChain(ctx, wg, goModGraph, rootModule, module, NewDependencyChain(dep), chains)
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
	goModGraph map[ModuleName][]ModuleName,
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

	for parent, deps := range goModGraph {
		if slices.Contains(deps, module) {
			wg.Add(1)
			go computeDependencyChain(ctx, wg, goModGraph, rootModule, parent, chain, chains)
		}
	}
}
