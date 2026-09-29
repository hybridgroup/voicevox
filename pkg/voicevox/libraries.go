package voicevox

import (
	"errors"
	"path/filepath"
	"runtime"

	"github.com/jupiterrider/ffi"
)

// LoadAdditionalLibraries loads the CUDA and cuDNN libraries in dir, so the ONNX Runtime
// CUDA provider can find them without LD_LIBRARY_PATH. Call it before LoadOnnxruntime.
func LoadAdditionalLibraries(dir string) error {
	pattern := "*.so*"
	switch runtime.GOOS {
	case "windows":
		pattern = "*.dll"
	case "darwin":
		pattern = "*.dylib"
	}

	pending, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return err
	}
	if len(pending) == 0 {
		return errors.New("no additional libraries found in " + dir)
	}

	// Some libraries depend on others in dir, so retry until no more can be loaded.
	for len(pending) > 0 {
		var failed []string
		var lastErr error
		for _, f := range pending {
			if _, err := ffi.Load(f); err != nil {
				failed = append(failed, f)
				lastErr = err
			}
		}

		if len(failed) == len(pending) {
			return lastErr
		}
		pending = failed
	}

	return nil
}
