package file_list

import (
	"fmt"
	"path/filepath"
	"sort"
)

// FileList represents a list of filenames with the ability to iterate over them.
type FileList struct {
	filenames    []string
	currentIndex int
}

// NewFileListFromGlob creates a FileList from a glob pattern.
func NewFileListFromGlob(pattern string) (*FileList, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob pattern: %w", err)
	}

	sort.Strings(matches) // Sort filenames like in Rust version
	return &FileList{
		filenames:    matches,
		currentIndex: 0,
	}, nil
}

// NewFileListFromList creates a FileList from an explicit list of filenames.
func NewFileListFromList(filenames []string) *FileList {
	sort.Strings(filenames) // Optional: Sort filenames
	return &FileList{
		filenames:    filenames,
		currentIndex: 0,
	}
}

// CurrentFilename returns the current filename in the list.
func (fl *FileList) CurrentFilename() string {
	if fl.currentIndex < len(fl.filenames) {
		return fl.filenames[fl.currentIndex]
	}
	return ""
}

// CurrentIndex returns the current index in the list.
func (fl *FileList) CurrentIndex() int {
	return fl.currentIndex
}

// Len returns the total number of files in the list.
func (fl *FileList) Len() int {
	return len(fl.filenames)
}

// NextFile advances to the next file in the list and returns its filename.
func (fl *FileList) NextFile() string {
	if fl.currentIndex+1 < len(fl.filenames) {
		fl.currentIndex++
	}
	return fl.CurrentFilename()
}
