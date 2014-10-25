package wbem

import (
	"io"
	"os"
	"path/filepath"
)

// Provider specified the interface types must implement to be used as a
// debugging sink. Having multiple such sink implementations allows it to be
// changed externally (for example when running tests).
type DebugProvider interface {
	NewFile(s string) io.WriteCloser
	Flush()
}

var currentDebugProvider DebugProvider = nil

func SetDebugProvider(p DebugProvider) {
	if currentDebugProvider != nil {
		currentDebugProvider.Flush()
	}
	currentDebugProvider = p
}

// Enabled returns whether debugging is enabled or not.
func DebugEnabled() bool {
	return currentDebugProvider != nil
}

// NewFile dispatches to the current provider's NewFile function.
func DebugNewFile(s string) io.WriteCloser {
	return currentDebugProvider.NewFile(s)
}

// Flush dispatches to the current provider's Flush function.
func DebugFlush() {
	currentDebugProvider.Flush()
}

// FileProvider implements a debugging provider that creates a real file for
// every call to NewFile. It maintains a list of all files that it creates,
// such that it can close them when its Flush function is called.
type FileDebugProvider struct {
	Path string

	files []*os.File
}

func (fp *FileDebugProvider) NewFile(p string) io.WriteCloser {
	f, err := os.Create(filepath.Join(fp.Path, p))
	if err != nil {
		panic(err)
	}

	fp.files = append(fp.files, f)

	return f
}

func (fp *FileDebugProvider) Flush() {
	for _, f := range fp.files {
		f.Close()
	}
}
