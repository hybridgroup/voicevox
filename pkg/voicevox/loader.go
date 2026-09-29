package voicevox

import (
	"fmt"

	"github.com/hybridgroup/voicevox/pkg/loader"
)

// Load loads the voicevox_core library from path, or from VOICEVOX_LIB if path is empty.
func Load(path string) error {
	lib, err := loader.LoadLibrary(path, "voicevox_core")
	if err != nil {
		return err
	}

	if err := loadVoicevoxFuncs(lib); err != nil {
		return err
	}
	if err := loadMemoryFuncs(lib); err != nil {
		return err
	}
	if err := loadOnnxruntimeFuncs(lib); err != nil {
		return err
	}
	if err := loadOpenJtalkFuncs(lib); err != nil {
		return err
	}
	if err := loadModelFuncs(lib); err != nil {
		return err
	}
	if err := loadSynthesizerFuncs(lib); err != nil {
		return err
	}

	return nil
}

func loadError(name string, err error) error {
	return fmt.Errorf("could not load %q: %w", name, err)
}
