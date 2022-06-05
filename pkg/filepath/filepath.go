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

// Files returns a list of filepathes of the coverage target files.
func (s *Fileset) Files(ext string) []string {
	fs := []string{}

	if val, ok := s.Set[ext]; ok {
		for f := range val {
			fs = append(fs, f)
		}
	}

	return fs
}

// Extract returns a Fileset from the given targetPath and excludePathes.
func Extract(targetPath string, excludePathes []string) (*Fileset, error) {
	excluded := make(map[string]struct{})
	for _, ep := range excludePathes {
		matches, err := zglob.Glob(ep)
		if err != nil {
			return nil, fmt.Errorf("failed to zglob.Glob: %w", err)
		}

		for _, match := range matches {
			excluded[match] = struct{}{}
		}
	}

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
