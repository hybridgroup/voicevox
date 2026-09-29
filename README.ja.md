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

## 音声

ダウンローダーで取得できる 0.16.0 の音声モデルに含まれるトーク用スタイルの一覧です。ID をスタイルとして `TTS` や `CreateAudioQuery` に渡し、そのスタイルを含むモデルファイルを読み込んでください。実行時にスタイルを一覧するには `VoiceModelFile.Metas()` または `Synthesizer.Metas()` を呼び出します。

キャラクターごとに利用規約があります。詳細は [VOICEVOX 公式サイト](https://voicevox.hiroshiba.jp/) を確認してください。

<details>
<summary>43 キャラクター、127 スタイル</summary>

| キャラクター | スタイル (ID と名前) | モデルファイル |
| --- | --- | --- |
| 四国めたん | 2 ノーマル、0 あまあま、4 セクシー、6 ツンツン、36 ささやき、37 ヒソヒソ | 0.vvm, 5.vvm |
| ずんだもん | 3 ノーマル、1 あまあま、5 セクシー、7 ツンツン、22 ささやき、38 ヒソヒソ、75 ヘロヘロ、76 なみだめ | 0.vvm, 5.vvm, 15.vvm |
| 春日部つむぎ | 8 ノーマル | 0.vvm |
| 波音リツ | 9 ノーマル、65 クイーン | 3.vvm |
| 雨晴はう | 10 ノーマル | 0.vvm |
| 玄野武宏 | 11 ノーマル、39 喜び、40 ツンギレ、41 悲しみ | 4.vvm, 10.vvm |
| 白上虎太郎 | 12 ふつう、32 わーい、33 びくびく、34 おこ、35 びえーん | 9.vvm |
| 青山龍星 | 13 ノーマル、81 熱血、82 不機嫌、83 喜び、84 しっとり、85 かなしみ、86 囁き | 15.vvm |
| 冥鳴ひまり | 14 ノーマル | 1.vvm |
| 九州そら | 16 ノーマル、15 あまあま、17 セクシー、18 ツンツン、19 ささやき | 2.vvm, 5.vvm |
| もち子さん | 20 ノーマル、66 セクシー／あん子、77 泣き、78 怒り、79 喜び、80 のんびり | 15.vvm |
| 剣崎雌雄 | 21 ノーマル | 4.vvm |
| WhiteCUL | 23 ノーマル、24 たのしい、25 かなしい、26 びえーん | 8.vvm |
| 後鬼 | 27 人間ver.、28 ぬいぐるみver.、87 人間（怒り）ver.、88 鬼ver. | 7.vvm, 16.vvm |
| No.7 | 29 ノーマル、30 アナウンス、31 読み聞かせ | 6.vvm |
| ちび式じい | 42 ノーマル | 10.vvm |
| 櫻歌ミコ | 43 ノーマル、44 第二形態、45 ロリ | 11.vvm |
| 小夜/SAYO | 46 ノーマル | 15.vvm |
| ナースロボ＿タイプＴ | 47 ノーマル、48 楽々、49 恐怖、50 内緒話 | 11.vvm |
| †聖騎士 紅桜† | 51 ノーマル | 12.vvm |
| 雀松朱司 | 52 ノーマル | 12.vvm |
| 麒ヶ島宗麟 | 53 ノーマル | 12.vvm |
| 春歌ナナ | 54 ノーマル | 13.vvm |
| 猫使アル | 55 ノーマル、56 おちつき、57 うきうき、110 つよつよ、111 へろへろ | 13.vvm, 21.vvm |
| 猫使ビィ | 58 ノーマル、59 おちつき、60 人見知り、112 つよつよ | 13.vvm, 21.vvm |
| 中国うさぎ | 61 ノーマル、62 おどろき、63 こわがり、64 へろへろ | 3.vvm |
| 栗田まろん | 67 ノーマル | 14.vvm |
| あいえるたん | 68 ノーマル | 14.vvm |
| 満別花丸 | 69 ノーマル、70 元気、71 ささやき、72 ぶりっ子、73 ボーイ | 14.vvm |
| 琴詠ニア | 74 ノーマル | 14.vvm |
| Voidoll | 89 ノーマル | 17.vvm |
| ぞん子 | 90 ノーマル、91 低血圧、92 覚醒、93 実況風 | 18.vvm |
| 中部つるぎ | 94 ノーマル、95 怒り、96 ヒソヒソ、97 おどおど、98 絶望と敗北 | 18.vvm |
| 離途 | 99 ノーマル、101 シリアス | 19.vvm |
| 黒沢冴白 | 100 ノーマル | 19.vvm |
| ユーレイちゃん | 102 ノーマル、103 甘々、104 哀しみ、105 ささやき、106 ツクモちゃん | 20.vvm |
| 東北ずん子 | 107 ノーマル | 21.vvm |
| 東北きりたん | 108 ノーマル | 21.vvm |
| 東北イタコ | 109 ノーマル | 21.vvm |
| あんこもん | 113 ノーマル、114 つよつよ、115 よわよわ、116 けだるげ、117 ささやき | 22.vvm, 23.vvm |
| 夜語トバリ | 118 ノーマル、119 明るい、120 哀しみ、121 呆れ | 24.vvm |
| 暁記ミタマ | 122 ノーマル、123 怒り、124 哀しみ、125 ささやき | 24.vvm |
| 里石ユカ | 126 つぼみ | 24.vvm |

</details>

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
