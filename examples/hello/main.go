package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hybridgroup/voicevox/pkg/voicevox"
)

var (
	dataDir = flag.String("data", "", "path to the voicevox_core directory from the downloader")
	model   = flag.String("model", "0.vvm", "voice model file in models/vvms")
	style   = flag.Uint("style", 0, "style id")
	text    = flag.String("text", "こんにちは、世界", "text to speak")
	output  = flag.String("o", "out.wav", "output file")
	gpu     = flag.Bool("gpu", false, "use the GPU")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	if *dataDir == "" {
		return fmt.Errorf("missing -data flag")
	}

	if err := voicevox.Load(filepath.Join(*dataDir, "c_api", "lib")); err != nil {
		return err
	}

	if *gpu {
		if err := voicevox.LoadAdditionalLibraries(filepath.Join(*dataDir, "additional_libraries")); err != nil {
			return err
		}
	}

	ort, err := voicevox.LoadOnnxruntime(filepath.Join(*dataDir, "onnxruntime", "lib"))
	if err != nil {
		return err
	}

	ojt, err := voicevox.NewOpenJtalk(filepath.Join(*dataDir, "dict", "open_jtalk_dic_utf_8-1.11"))
	if err != nil {
		return err
	}

	opts := voicevox.DefaultInitializeOptions()
	if *gpu {
		opts.AccelerationMode = voicevox.AccelerationModeGPU
	}

	syn, err := voicevox.NewSynthesizer(ort, ojt, opts)
	ojt.Delete()
	if err != nil {
		return err
	}
	defer syn.Delete()

	vvm, err := voicevox.OpenVoiceModelFile(filepath.Join(*dataDir, "models", "vvms", *model))
	if err != nil {
		return err
	}
	err = syn.LoadVoiceModel(vvm, voicevox.DefaultLoadVoiceModelOptions())
	vvm.Delete()
	if err != nil {
		return err
	}

	wav, err := syn.TTS(*text, voicevox.StyleID(*style), voicevox.DefaultTTSOptions())
	if err != nil {
		return err
	}

	fmt.Println("voicevox_core", voicevox.Version(), "gpu:", syn.IsGPUMode())
	return os.WriteFile(*output, wav, 0644)
}
