package voicevox

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/sys/unix"
)

func TestVersion(t *testing.T) {
	testSetup(t)

	if v := Version(); v == "" {
		t.Fatal("empty version")
	}
}

func TestResultCodeError(t *testing.T) {
	testSetup(t)

	if msg := ResultStyleNotFound.Error(); msg == "" {
		t.Fatal("empty error message")
	}
	if ResultStyleNotFound.Error() == ResultModelNotFound.Error() {
		t.Fatal("different codes have the same message")
	}
}

func TestDefaultOptions(t *testing.T) {
	testSetup(t)

	initOpts := DefaultInitializeOptions()
	if initOpts.AccelerationMode < AccelerationModeAuto || initOpts.AccelerationMode > AccelerationModeGPU {
		t.Fatalf("invalid acceleration mode %d", initOpts.AccelerationMode)
	}

	if opts := DefaultLoadVoiceModelOptions(); opts.OnExisting != OnExistingVoiceModelIDError {
		t.Fatalf("invalid on existing %d", opts.OnExisting)
	}

	if opts := DefaultTTSOptions(); opts.EnableInterrogativeUpspeak > 1 {
		t.Fatalf("invalid tts bool %d", opts.EnableInterrogativeUpspeak)
	}

	if opts := DefaultSynthesisOptions(); opts.EnableInterrogativeUpspeak > 1 {
		t.Fatalf("invalid synthesis bool %d", opts.EnableInterrogativeUpspeak)
	}

	opts := DefaultLoadOnnxruntimeOptions()
	if opts.Filename == nil {
		t.Fatal("nil onnxruntime filename")
	}
	if got, want := unix.BytePtrToString(opts.Filename), OnnxruntimeLibFilename(); got != want {
		t.Fatalf("onnxruntime filename %q, want %q", got, want)
	}
}

func TestOnnxruntimeLibFilename(t *testing.T) {
	testSetup(t)

	if name := OnnxruntimeLibFilename(); !strings.Contains(name, "onnxruntime") {
		t.Fatalf("unexpected filename %q", name)
	}
}

func TestOpenVoiceModelFileMissing(t *testing.T) {
	testSetup(t)

	if _, err := OpenVoiceModelFile("does-not-exist.vvm"); err == nil {
		t.Fatal("expected error")
	}
}

func TestTTS(t *testing.T) {
	syn, style := testSynthesizer(t)

	wav, err := syn.TTS("こんにちは", style, DefaultTTSOptions())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(wav, []byte("RIFF")) {
		t.Fatal("output is not a WAV file")
	}
}

func TestTTSBadStyle(t *testing.T) {
	syn, _ := testSynthesizer(t)

	_, err := syn.TTS("こんにちは", 99999, DefaultTTSOptions())
	if !errors.Is(err, ResultStyleNotFound) {
		t.Fatalf("got %v, want %v", err, ResultStyleNotFound)
	}
}

func TestAudioQuerySynthesis(t *testing.T) {
	syn, style := testSynthesizer(t)

	query, err := syn.CreateAudioQuery("こんにちは", style)
	if err != nil {
		t.Fatal(err)
	}

	normal, err := syn.Synthesis(query, style, DefaultSynthesisOptions())
	if err != nil {
		t.Fatal(err)
	}

	var q map[string]any
	if err := json.Unmarshal([]byte(query), &q); err != nil {
		t.Fatal(err)
	}
	q["speedScale"] = 2.0
	fast, _ := json.Marshal(q)

	faster, err := syn.Synthesis(string(fast), style, DefaultSynthesisOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(faster) >= len(normal) {
		t.Fatalf("speedScale 2.0 gave %d bytes, normal gave %d", len(faster), len(normal))
	}
}

func TestSynthesizerMetas(t *testing.T) {
	syn, style := testSynthesizer(t)

	metas, err := syn.Metas()
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range metas {
		for _, s := range c.Styles {
			if s.ID == style {
				return
			}
		}
	}
	t.Fatalf("style %d not in synthesizer metas", style)
}

func TestSupportedDevices(t *testing.T) {
	ort := testOnnxruntime(t)

	d, err := ort.SupportedDevices()
	if err != nil {
		t.Fatal(err)
	}
	if !d.CPU {
		t.Fatal("cpu not supported")
	}
}

func TestCPUMode(t *testing.T) {
	testSetup(t)
	opts := DefaultInitializeOptions()
	opts.AccelerationMode = AccelerationModeCPU
	syn, _ := testSynthesizerWithOptions(t, opts)

	if syn.IsGPUMode() {
		t.Fatal("cpu mode synthesizer reports gpu mode")
	}
}

func TestGPUMode(t *testing.T) {
	if os.Getenv("VOICEVOX_TEST_GPU") == "" {
		t.Skip("no VOICEVOX_TEST_GPU skipping test")
	}

	testSetup(t)
	if err := LoadAdditionalLibraries(filepath.Join(testDataDir(t), "additional_libraries")); err != nil {
		t.Fatal("unable to load additional libraries", err)
	}

	opts := DefaultInitializeOptions()
	opts.AccelerationMode = AccelerationModeGPU
	syn, style := testSynthesizerWithOptions(t, opts)

	if !syn.IsGPUMode() {
		t.Fatal("gpu mode synthesizer reports cpu mode")
	}
	if _, err := syn.TTS("こんにちは", style, DefaultTTSOptions()); err != nil {
		t.Fatal(err)
	}
}
