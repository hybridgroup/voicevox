package voicevox

import (
	"os"
	"path/filepath"
	"testing"
)

func testSetup(t *testing.T) {
	if os.Getenv("VOICEVOX_LIB") == "" {
		t.Fatal("no VOICEVOX_LIB set for tests")
	}
	if err := Load(os.Getenv("VOICEVOX_LIB")); err != nil {
		t.Fatal("unable to load library", err.Error())
	}
}

func testDataDir(t *testing.T) string {
	if os.Getenv("VOICEVOX_DIR") == "" {
		t.Skip("no VOICEVOX_DIR skipping test")
	}
	return os.Getenv("VOICEVOX_DIR")
}

func testOnnxruntime(t *testing.T) Onnxruntime {
	testSetup(t)
	dir := testDataDir(t)

	ort, err := LoadOnnxruntime(filepath.Join(dir, "onnxruntime", "lib"))
	if err != nil {
		t.Fatal("unable to load onnxruntime", err)
	}
	return ort
}

// testSynthesizer loads the first voice model and returns its first talk style.
func testSynthesizer(t *testing.T) (Synthesizer, StyleID) {
	testSetup(t)
	return testSynthesizerWithOptions(t, DefaultInitializeOptions())
}

func testSynthesizerWithOptions(t *testing.T, opts InitializeOptions) (Synthesizer, StyleID) {
	ort := testOnnxruntime(t)
	dir := testDataDir(t)

	ojt, err := NewOpenJtalk(filepath.Join(dir, "dict", "open_jtalk_dic_utf_8-1.11"))
	if err != nil {
		t.Fatal("unable to load open jtalk dictionary", err)
	}
	defer ojt.Delete()

	syn, err := NewSynthesizer(ort, ojt, opts)
	if err != nil {
		t.Fatal("unable to create synthesizer", err)
	}
	t.Cleanup(syn.Delete)

	files, _ := filepath.Glob(filepath.Join(dir, "models", "vvms", "*.vvm"))
	if len(files) == 0 {
		t.Fatal("no voice models found")
	}

	model, err := OpenVoiceModelFile(files[0])
	if err != nil {
		t.Fatal("unable to open voice model", err)
	}
	defer model.Delete()

	metas, err := model.Metas()
	if err != nil {
		t.Fatal("unable to read voice model metas", err)
	}

	if err := syn.LoadVoiceModel(model, DefaultLoadVoiceModelOptions()); err != nil {
		t.Fatal("unable to load voice model", err)
	}

	for _, c := range metas {
		for _, s := range c.Styles {
			if s.Type == "talk" {
				return syn, s.ID
			}
		}
	}

	t.Fatal("no talk style in", files[0])
	return 0, 0
}
