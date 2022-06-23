package report

import (
	"fmt"
	"sort"

	"github.com/commentcov/commentcov/proto"
)

// Mode controls at what granularity coverage reporting is performed.
type Mode string

const (
	// ModeInvalid is invalid mode.
	ModeInvalid Mode = ""
	// ModeFile reports coverage per file.
	ModeFile Mode = "file"
	// ModeFile reports coverage by node type.
	ModeScope Mode = "scope"
	// ModeFileScope reports coverage by node type per file.
	ModeFileScope Mode = "file_scope"
)

// StringToMode classifies from string to Mode.
func StringToMode(str string) (Mode, error) {
	switch str {
	case "file":
		return ModeFile, nil
	case "scope":
		return ModeScope, nil
	case "file_scope":
		return ModeFileScope, nil
	default:
		return ModeInvalid, fmt.Errorf("invalid string for Mode: %s", str)
	}
}

// Report will report coverage according to mode values.
func Report(mode Mode, items []*proto.CoverageItem) {
	switch mode {
	case ModeFile:
		byFile(items)
	case ModeScope:
		byScope(items)
	case ModeFileScope:
		byFileScope(items)
	default:
		// do nothing
	}
}

// byFile reports coverage per file.
func byFile(items []*proto.CoverageItem) {
	cc, files, scopes := Profile(items)

	sc := make(ScopedCounter)
	for _, file := range files {
		sc[file] = NewCounter()

		for _, scope := range scopes {
			if counter, ok := cc[file][scope]; ok {
				sc[file].Merge(counter)
			}
		}
	}

	for _, file := range files {
		percent := sc[file].CalcRate()
		fmt.Printf("%v,,%v\n", file, percent)
	}
}

// byScope reports coverage by node type.
func byScope(items []*proto.CoverageItem) {
	cc, files, scopes := Profile(items)

	sc := make(ScopedCounter)
	for _, scope := range scopes {
		sc[scope] = NewCounter()
	}

	for _, file := range files {
		for _, scope := range scopes {
			if counter, ok := cc[file][scope]; ok {
				sc[scope].Merge(counter)
			}
		}
	}

	for _, scope := range scopes {
		percent := sc[scope].CalcRate()
		fmt.Printf(",%v,%v\n", scope, percent)
	}
}

// byFileScope reports coverage by node type per file.
func byFileScope(items []*proto.CoverageItem) {
	cc, files, scopes := Profile(items)

	for _, file := range files {
		for _, scope := range scopes {
			if counter, ok := cc[file][scope]; ok {
				percent := counter.CalcRate()
				fmt.Printf("%v,%v,%v\n", file, scope, percent)
			}
		}
	}
}

// Profile converts a list of items into the map of filename:scope:couter.
func Profile(items []*proto.CoverageItem) (map[string]ScopedCounter, []string, []string) {
	cc := make(map[string]ScopedCounter)
	files := []string{}
	scopes := []string{}

	scopeMemo := make(map[string]struct{})
	unknownScope := proto.CoverageItem_UNKNOWN.String()
	for _, item := range items {
		scope := item.Scope.String()

		if scope == unknownScope {
			continue
		}

		if _, ok := scopeMemo[scope]; !ok {
			scopeMemo[scope] = struct{}{}
			scopes = append(scopes, scope)
		}

		if _, ok := cc[item.File]; !ok {
			cc[item.File] = make(ScopedCounter)
			files = append(files, item.File)
		}

		if _, ok := cc[item.File][scope]; !ok {
			cc[item.File][scope] = NewCounter()
		}

		cc[item.File][scope].Add(item)
	}

	sort.Strings(files)
	sort.Strings(scopes)

	return cc, files, scopes
}
