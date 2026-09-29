# voicevox

[![Go Reference](https://pkg.go.dev/badge/github.com/hybridgroup/voicevox.svg)](https://pkg.go.dev/github.com/hybridgroup/voicevox)

Go bindings for [VOICEVOX core](https://github.com/VOICEVOX/voicevox_core), for Japanese text to speech that runs locally.

- Uses [`purego`](https://github.com/ebitengine/purego) and [`ffi`](https://github.com/JupiterRider/ffi), so CGo is not needed.
- Targets voicevox_core 0.17.0.

## Installation

### libffi

On Linux, install libffi.

```shell
sudo apt install libffi8
```

On macOS and Windows x64, the `ffi` package ships libffi, so you don't need to install it.

### VOICEVOX runtime

Get the official downloader from the [voicevox_core releases](https://github.com/VOICEVOX/voicevox_core/releases/tag/0.17.0), then run it:

```shell
curl -sSfL https://github.com/VOICEVOX/voicevox_core/releases/download/0.17.0/download-linux-x64 -o download
chmod +x download
./download -o ~/voicevox_core
```

The downloader asks you to accept the VOICEVOX voice model terms and the VOICEVOX ONNX Runtime terms. To get only the talk models, add `--models-pattern '[0-9]*.vvm'`.

This creates the following layout:

```
voicevox_core/
  c_api/lib/libvoicevox_core.so
  onnxruntime/lib/libvoicevox_onnxruntime.so.1.23.2
  dict/open_jtalk_dic_utf_8-1.11/
  models/vvms/*.vvm
```

## Example

```go
package main

import (
	"os"

	"github.com/hybridgroup/voicevox/pkg/voicevox"
)

func main() {
	voicevox.Load("voicevox_core/c_api/lib")

	ort, _ := voicevox.LoadOnnxruntime("voicevox_core/onnxruntime/lib")
	ojt, _ := voicevox.NewOpenJtalk("voicevox_core/dict/open_jtalk_dic_utf_8-1.11")

	syn, _ := voicevox.NewSynthesizer(ort, ojt, voicevox.DefaultInitializeOptions())
	ojt.Delete()
	defer syn.Delete()

	model, _ := voicevox.OpenVoiceModelFile("voicevox_core/models/vvms/0.vvm")
	syn.LoadVoiceModel(model, voicevox.DefaultLoadVoiceModelOptions())
	model.Delete()

	wav, _ := syn.TTS("こんにちは", 0, voicevox.DefaultTTSOptions())
	os.WriteFile("out.wav", wav, 0644)
}
```

The output is a 24 kHz, 16-bit, mono WAV file.

To change speed, pitch, intonation or volume:
1. Call `CreateAudioQuery`.
2. Edit `speedScale`, `pitchScale`, `intonationScale` or `volumeScale` in the returned JSON.
3. Pass the JSON to `Synthesis`.

A full example is in [examples/hello](./examples/hello):

```shell
go run ./examples/hello -data ~/voicevox_core -style 0 -text "こんにちは"
```

## GPU

To use CUDA, download the CUDA build of ONNX Runtime and the CUDA libraries into the same directory:

```shell
./download -o ~/voicevox_core --devices cuda --only onnxruntime additional-libraries
```

The ONNX Runtime CUDA provider needs to find the CUDA and cuDNN libraries in `additional_libraries`. Call `LoadAdditionalLibraries` before `LoadOnnxruntime`, then set the acceleration mode when you create the synthesizer:

```go
voicevox.LoadAdditionalLibraries("voicevox_core/additional_libraries")
ort, _ := voicevox.LoadOnnxruntime("voicevox_core/onnxruntime/lib")

opts := voicevox.DefaultInitializeOptions()
opts.AccelerationMode = voicevox.AccelerationModeGPU
syn, err := voicevox.NewSynthesizer(ort, ojt, opts)
```

`NewSynthesizer` returns an error if the GPU can't be used. Call `syn.IsGPUMode()` to see which mode is active, and `ort.SupportedDevices()` to see which devices the ONNX Runtime build supports.

```shell
go run ./examples/hello -data ~/voicevox_core -gpu
```

## Tests

```shell
VOICEVOX_LIB=~/voicevox_core/c_api/lib VOICEVOX_DIR=~/voicevox_core go test ./...
```

`VOICEVOX_LIB` is required. Tests that need voice models are skipped when `VOICEVOX_DIR` is not set. To also run the GPU test, set `VOICEVOX_TEST_GPU=1`.

## Credits

The VOICEVOX voice model terms require a credit for each character you use, for example "VOICEVOX:ずんだもん". Check the terms for each character.

## License

MIT. VOICEVOX core, its voice models and its ONNX Runtime each have their own licenses and terms.
