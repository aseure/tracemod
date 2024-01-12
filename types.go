package main

type Direction string

func (d Direction) String() string {
	return string(d)
}

func (d *Direction) Set(v string) error {
	switch v {
	case "TB", "BT", "LR", "RL":
		*d = Direction(v)
		return nil
	default:
		*d = ""
	}
	return nil
}

func (d Direction) Type() string {
	return `"TB"|"BT"|"LR"|"RL"`
}

type ModuleName string

type ModuleMatchingFunc func(m ModuleName) bool

type DependencyChain struct {
	chain []ModuleName
}

func NewDependencyChain(startModule ModuleName) DependencyChain {
	return DependencyChain{
		chain: []ModuleName{startModule},
	}
}

func (c DependencyChain) Add(dep ModuleName) DependencyChain {
	newChain := make([]ModuleName, len(c.chain))
	copy(newChain, c.chain)
	return DependencyChain{
		chain: append(newChain, dep),
	}
}

func (c DependencyChain) Has(dep ModuleName) bool {
	for _, d := range c.chain {
		if dep == d {
			return true
		}
	}
	return false
}

func (c DependencyChain) Deps() []ModuleName {
	return c.chain
}
