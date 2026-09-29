# voicevox

[English](./README.md) | 日本語

[![Go Reference](https://pkg.go.dev/badge/github.com/hybridgroup/voicevox.svg)](https://pkg.go.dev/github.com/hybridgroup/voicevox)

[VOICEVOX コア](https://github.com/VOICEVOX/voicevox_core) の Go バインディングです。ローカルで動く日本語の音声合成を Go から使えます。

- [`purego`](https://github.com/ebitengine/purego) と [`ffi`](https://github.com/JupiterRider/ffi) を使うので、CGo は不要です。
- voicevox_core 0.17.0 に対応しています。

## インストール

### libffi

Linux では libffi をインストールしてください。

```shell
sudo apt install libffi8
```

macOS と Windows x64 では `ffi` パッケージに libffi が同梱されているので、インストールは不要です。

### VOICEVOX ランタイム

[voicevox_core のリリース](https://github.com/VOICEVOX/voicevox_core/releases/tag/0.17.0) から公式のダウンローダーを取得して実行します。

```shell
curl -sSfL https://github.com/VOICEVOX/voicevox_core/releases/download/0.17.0/download-linux-x64 -o download
chmod +x download
./download -o ~/voicevox_core
```

ダウンローダーの実行中に「VOICEVOX 音声モデル 利用規約」と「VOICEVOX ONNX Runtime 利用規約」への同意を求められます。トーク用のモデルだけが必要な場合は `--models-pattern '[0-9]*.vvm'` を追加してください。

次の構成でファイルが作成されます。

```
voicevox_core/
  c_api/lib/libvoicevox_core.so
  onnxruntime/lib/libvoicevox_onnxruntime.so.<version>
  dict/open_jtalk_dic_utf_8-1.11/
  models/vvms/*.vvm
```

`LoadOnnxruntime` は推奨バージョンの ONNX Runtime があればそれを使い、なければディレクトリ内で最も新しいバージョンを使います。

## 使用例

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

出力は 24 kHz、16 ビット、モノラルの WAV ファイルです。

話速、音高、抑揚、音量を変えるには、次の手順で行います。
1. `CreateAudioQuery` を呼び出します。
2. 返された JSON の `speedScale`、`pitchScale`、`intonationScale`、`volumeScale` を編集します。
3. その JSON を `Synthesis` に渡します。

完全な例は [examples/hello](./examples/hello) にあります。

```shell
go run ./examples/hello -data ~/voicevox_core -style 0 -text "こんにちは"
```

## GPU

CUDA を使う場合は、CUDA 版の ONNX Runtime と CUDA ライブラリを同じディレクトリにダウンロードします。

```shell
./download -o ~/voicevox_core --devices cuda --only onnxruntime additional-libraries
```

ONNX Runtime の CUDA プロバイダーは、`additional_libraries` にある CUDA と cuDNN のライブラリを読み込めなければなりません。`LoadOnnxruntime` の前に `LoadAdditionalLibraries` を呼び出し、シンセサイザーを作成するときにアクセラレーションモードを指定してください。

```go
voicevox.LoadAdditionalLibraries("voicevox_core/additional_libraries")
ort, _ := voicevox.LoadOnnxruntime("voicevox_core/onnxruntime/lib")

opts := voicevox.DefaultInitializeOptions()
opts.AccelerationMode = voicevox.AccelerationModeGPU
syn, err := voicevox.NewSynthesizer(ort, ojt, opts)
```

GPU が使えない場合、`NewSynthesizer` はエラーを返します。現在のモードは `syn.IsGPUMode()` で確認できます。ONNX Runtime のビルドが対応しているデバイスは `ort.SupportedDevices()` で確認できます。

```shell
go run ./examples/hello -data ~/voicevox_core -gpu
```

## テスト

```shell
VOICEVOX_LIB=~/voicevox_core/c_api/lib VOICEVOX_DIR=~/voicevox_core go test ./...
```

`VOICEVOX_LIB` は必須です。`VOICEVOX_DIR` が設定されていない場合、音声モデルが必要なテストはスキップされます。GPU のテストも実行するには `VOICEVOX_TEST_GPU=1` を設定してください。

## クレジット

VOICEVOX 音声モデルの利用規約では、使用するキャラクターごとにクレジットの記載が必要です。例えば「VOICEVOX:ずんだもん」のように記載します。各キャラクターの利用規約を確認してください。

## ライセンス

MIT ライセンスです。VOICEVOX コア、音声モデル、ONNX Runtime には、それぞれ別のライセンスと利用規約があります。
