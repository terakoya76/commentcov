package filepath

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattn/go-zglob"
)

// Fileset represents the coverage target file list.
// K: extension, V: Set<filepath>.
type Fileset struct {
	Set map[string]map[string]struct{}
}

// Extensions returns a list of extensions of the coverage target files.
func (s *Fileset) Extensions() []string {
	exts := []string{}

	for k := range s.Set {
		exts = append(exts, k)
	}

	return exts
}

// Files returns a list of filepaths of the coverage target files.
func (s *Fileset) Files(ext string) []string {
	fs := []string{}

	if val, ok := s.Set[ext]; ok {
		for f := range val {
			fs = append(fs, f)
		}
	}

	return fs
}

// ExcludeFileSet represents a list of excluded files as Set.
type ExcludeFileSet = map[string]struct{}

// NewExcludeFileSet returns a excluded File Set from the given excludePaths.
func NewExcludeFileSet(excludePaths []string) (ExcludeFileSet, error) {
	excluded := make(map[string]struct{})
	for _, ep := range excludePaths {
		matches, err := zglob.Glob(ep)
		if err != nil {
			return nil, fmt.Errorf("failed to zglob.Glob: %w", err)
		}

		for _, match := range matches {
			path, err := filepath.Abs(match)
			if err != nil {
				return nil, fmt.Errorf("failed to filepath.Abs: %w", err)
			}

			excluded[path] = struct{}{}
		}
	}

	return excluded, nil
}

// Extract returns a Fileset from the given targetPath and excludePaths.
func Extract(targetPath string, excluded ExcludeFileSet) (*Fileset, error) {
	set := &Fileset{
		Set: make(map[string]map[string]struct{}),
	}

	walkFn := func(pathname string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		if fi.IsDir() {
			return nil
		}

		ext := filepath.Ext(pathname)
		f := filepath.Clean(pathname)

		if _, ok := excluded[f]; ok {
			return nil
		}

		if _, ok := set.Set[ext]; !ok {
			set.Set[ext] = make(map[string]struct{})
		}

		set.Set[ext][f] = struct{}{}
		return nil
	}

	if err := filepath.Walk(targetPath, walkFn); err != nil {
		return set, fmt.Errorf("%w", err)
	}

	return set, nil
}
