package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jupiterrider/ffi"
)

// Lib is a handle to a shared library.
type Lib struct {
	lib ffi.Lib
}

// Prep gets the address of a function and describes its signature.
func (l Lib) Prep(name string, ret *ffi.Type, args ...*ffi.Type) (ffi.Fun, error) {
	return l.lib.Prep(name, ret, args...)
}

// LoadLibrary loads a shared library by short name from path, or from VOICEVOX_LIB if path is empty.
func LoadLibrary(path, lib string) (Lib, error) {
	if path == "" {
		path = os.Getenv("VOICEVOX_LIB")
	}

	if path == "" {
		return Lib{}, fmt.Errorf("library path not specified and VOICEVOX_LIB env variable not set")
	}

	l, err := ffi.Load(GetLibraryFilename(path, lib))
	if err != nil {
		return Lib{}, err
	}

	return Lib{lib: l}, nil
}

// GetLibraryFilename returns the full path to the library file for the current OS.
func GetLibraryFilename(path, lib string) string {
	switch runtime.GOOS {
	case "linux", "freebsd":
		return filepath.Join(path, fmt.Sprintf("lib%s.so", lib))
	case "windows":
		return filepath.Join(path, fmt.Sprintf("%s.dll", lib))
	case "darwin":
		return filepath.Join(path, fmt.Sprintf("lib%s.dylib", lib))
	default:
		return filepath.Join(path, lib)
	}
}
