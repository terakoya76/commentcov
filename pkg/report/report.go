package report

import (
	"fmt"

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
	sc := make(scopedCounter)

	for _, item := range items {
		scope := item.File

		if _, ok := sc[scope]; !ok {
			sc[scope] = NewCounter()
		}

		sc[scope].Profile(item)
	}

	for scope, counter := range sc {
		percent := counter.CalcRate()
		fmt.Printf("%v: %v\n", scope, percent)
	}
}

// byScope reports coverage by node type.
func byScope(items []*proto.CoverageItem) {
	sc := make(scopedCounter)

	unknownScope := proto.CoverageItem_UNKNOWN.String()
	for _, item := range items {
		scope := item.Scope.String()

		if scope == unknownScope {
			continue
		}

		if _, ok := sc[scope]; !ok {
			sc[scope] = NewCounter()
		}

		sc[scope].Profile(item)
	}

	for scope, counter := range sc {
		percent := counter.CalcRate()
		fmt.Printf("%v: %v\n", scope, percent)
	}
}

// byFileScope reports coverage by node type per file.
func byFileScope(items []*proto.CoverageItem) {
	cc := make(map[string]scopedCounter)

	unknownScope := proto.CoverageItem_UNKNOWN.String()
	for _, item := range items {
		scope := item.Scope.String()

		if scope == unknownScope {
			continue
		}

		if _, ok := cc[item.File]; !ok {
			cc[item.File] = make(scopedCounter)
		}

		if _, ok := cc[item.File][scope]; !ok {
			cc[item.File][scope] = NewCounter()
		}

		cc[item.File][scope].Profile(item)
	}

	for file, sc := range cc {
		for scope, counter := range sc {
			percent := counter.CalcRate()
			fmt.Printf("%v,%v: %v\n", file, scope, percent)
		}
	}
}
