# voicevox

English | [日本語](./README.ja.md)

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
  onnxruntime/lib/libvoicevox_onnxruntime.so.<version>
  dict/open_jtalk_dic_utf_8-1.11/
  models/vvms/*.vvm
```

`LoadOnnxruntime` uses the recommended ONNX Runtime version if it is in the directory. Otherwise it uses the newest version it finds there.

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

## Voices

These are the talk styles in the 0.16.0 voice models from the downloader. Pass the id as the style to `TTS` or `CreateAudioQuery`, and load the model file that contains it. To list the styles at runtime, call `VoiceModelFile.Metas()` or `Synthesizer.Metas()`.

Each character has its own terms. See the [VOICEVOX website](https://voicevox.hiroshiba.jp/) for details.

<details>
<summary>43 characters, 127 styles</summary>

| Character | Styles (id, name) | Model files |
| --- | --- | --- |
| 四国めたん | 2 ノーマル (normal), 0 あまあま (sweet), 4 セクシー (sexy), 6 ツンツン (prickly), 36 ささやき (whisper), 37 ヒソヒソ (hushed) | 0.vvm, 5.vvm |
| ずんだもん | 3 ノーマル (normal), 1 あまあま (sweet), 5 セクシー (sexy), 7 ツンツン (prickly), 22 ささやき (whisper), 38 ヒソヒソ (hushed), 75 ヘロヘロ (exhausted), 76 なみだめ (teary) | 0.vvm, 5.vvm, 15.vvm |
| 春日部つむぎ | 8 ノーマル (normal) | 0.vvm |
| 波音リツ | 9 ノーマル (normal), 65 クイーン (queen) | 3.vvm |
| 雨晴はう | 10 ノーマル (normal) | 0.vvm |
| 玄野武宏 | 11 ノーマル (normal), 39 喜び (joy), 40 ツンギレ (snappy), 41 悲しみ (sadness) | 4.vvm, 10.vvm |
| 白上虎太郎 | 12 ふつう (normal), 32 わーい (yay), 33 びくびく (nervous), 34 おこ (angry), 35 びえーん (crying) | 9.vvm |
| 青山龍星 | 13 ノーマル (normal), 81 熱血 (passionate), 82 不機嫌 (grumpy), 83 喜び (joy), 84 しっとり (calm), 85 かなしみ (sadness), 86 囁き (whisper) | 15.vvm |
| 冥鳴ひまり | 14 ノーマル (normal) | 1.vvm |
| 九州そら | 16 ノーマル (normal), 15 あまあま (sweet), 17 セクシー (sexy), 18 ツンツン (prickly), 19 ささやき (whisper) | 2.vvm, 5.vvm |
| もち子さん | 20 ノーマル (normal), 66 セクシー／あん子 (sexy/Anko), 77 泣き (crying), 78 怒り (angry), 79 喜び (joy), 80 のんびり (relaxed) | 15.vvm |
| 剣崎雌雄 | 21 ノーマル (normal) | 4.vvm |
| WhiteCUL | 23 ノーマル (normal), 24 たのしい (happy), 25 かなしい (sad), 26 びえーん (crying) | 8.vvm |
| 後鬼 | 27 人間ver. (human), 28 ぬいぐるみver. (plush), 87 人間（怒り）ver. (human, angry), 88 鬼ver. (demon) | 7.vvm, 16.vvm |
| No.7 | 29 ノーマル (normal), 30 アナウンス (announcer), 31 読み聞かせ (storytelling) | 6.vvm |
| ちび式じい | 42 ノーマル (normal) | 10.vvm |
| 櫻歌ミコ | 43 ノーマル (normal), 44 第二形態 (second form), 45 ロリ (childlike) | 11.vvm |
| 小夜/SAYO | 46 ノーマル (normal) | 15.vvm |
| ナースロボ＿タイプＴ | 47 ノーマル (normal), 48 楽々 (relaxed), 49 恐怖 (fear), 50 内緒話 (secret talk) | 11.vvm |
| †聖騎士 紅桜† | 51 ノーマル (normal) | 12.vvm |
| 雀松朱司 | 52 ノーマル (normal) | 12.vvm |
| 麒ヶ島宗麟 | 53 ノーマル (normal) | 12.vvm |
| 春歌ナナ | 54 ノーマル (normal) | 13.vvm |
| 猫使アル | 55 ノーマル (normal), 56 おちつき (calm), 57 うきうき (cheerful), 110 つよつよ (strong), 111 へろへろ (exhausted) | 13.vvm, 21.vvm |
| 猫使ビィ | 58 ノーマル (normal), 59 おちつき (calm), 60 人見知り (shy), 112 つよつよ (strong) | 13.vvm, 21.vvm |
| 中国うさぎ | 61 ノーマル (normal), 62 おどろき (surprised), 63 こわがり (scared), 64 へろへろ (exhausted) | 3.vvm |
| 栗田まろん | 67 ノーマル (normal) | 14.vvm |
| あいえるたん | 68 ノーマル (normal) | 14.vvm |
| 満別花丸 | 69 ノーマル (normal), 70 元気 (energetic), 71 ささやき (whisper), 72 ぶりっ子 (cutesy), 73 ボーイ (boy) | 14.vvm |
| 琴詠ニア | 74 ノーマル (normal) | 14.vvm |
| Voidoll | 89 ノーマル (normal) | 17.vvm |
| ぞん子 | 90 ノーマル (normal), 91 低血圧 (groggy), 92 覚醒 (awakened), 93 実況風 (commentator) | 18.vvm |
| 中部つるぎ | 94 ノーマル (normal), 95 怒り (angry), 96 ヒソヒソ (hushed), 97 おどおど (timid), 98 絶望と敗北 (despair and defeat) | 18.vvm |
| 離途 | 99 ノーマル (normal), 101 シリアス (serious) | 19.vvm |
| 黒沢冴白 | 100 ノーマル (normal) | 19.vvm |
| ユーレイちゃん | 102 ノーマル (normal), 103 甘々 (sweet), 104 哀しみ (sorrow), 105 ささやき (whisper), 106 ツクモちゃん (Tsukumo-chan) | 20.vvm |
| 東北ずん子 | 107 ノーマル (normal) | 21.vvm |
| 東北きりたん | 108 ノーマル (normal) | 21.vvm |
| 東北イタコ | 109 ノーマル (normal) | 21.vvm |
| あんこもん | 113 ノーマル (normal), 114 つよつよ (strong), 115 よわよわ (weak), 116 けだるげ (languid), 117 ささやき (whisper) | 22.vvm, 23.vvm |
| 夜語トバリ | 118 ノーマル (normal), 119 明るい (bright), 120 哀しみ (sorrow), 121 呆れ (exasperated) | 24.vvm |
| 暁記ミタマ | 122 ノーマル (normal), 123 怒り (angry), 124 哀しみ (sorrow), 125 ささやき (whisper) | 24.vvm |
| 里石ユカ | 126 つぼみ (Tsubomi) | 24.vvm |

</details>

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
