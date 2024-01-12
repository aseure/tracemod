package main

import (
	"fmt"
	"regexp"
)

func buildModuleMatchingFunc(moduleFilter string, isExactMatch bool) (ModuleMatchingFunc, error) {
	if isExactMatch {
		return func(m ModuleName) bool {
			return moduleFilter == string(m)
		}, nil
	}

	re, err := regexp.Compile(moduleFilter)
	if err != nil {
		return nil, fmt.Errorf("invalid regex filter %q: %w", moduleFilter, err)
	}

	return func(m ModuleName) bool {
		return re.MatchString(string(m))
	}, nil
}
