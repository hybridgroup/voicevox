package voicevox

import (
	"errors"
	"unsafe"

	"github.com/hybridgroup/voicevox/pkg/loader"
	"github.com/jupiterrider/ffi"
)

// InitializeOptions are the options for NewSynthesizer.
type InitializeOptions struct {
	AccelerationMode AccelerationMode
	CPUNumThreads    uint16
}

// LoadVoiceModelOptions are the options for LoadVoiceModel.
type LoadVoiceModelOptions struct {
	OnExisting OnExistingVoiceModelID
}

// TTSOptions are the options for TTS and TTSFromKana.
type TTSOptions struct {
	EnableInterrogativeUpspeak uint8 // bool as uint8
}

// SynthesisOptions are the options for Synthesis.
type SynthesisOptions struct {
	EnableInterrogativeUpspeak uint8 // bool as uint8
}

var errInvalidSynthesizer = errors.New("invalid synthesizer")

var (
	ffiTypeInitializeOptions     = ffi.NewType(&ffi.TypeSint32, &ffi.TypeUint16)
	ffiTypeLoadVoiceModelOptions = ffi.NewType(&ffi.TypeSint32)
	ffiTypeTTSOptions            = ffi.NewType(&ffi.TypeUint8)
	ffiTypeSynthesisOptions      = ffi.NewType(&ffi.TypeUint8)

	// VOICEVOX_C_API struct VoicevoxInitializeOptions voicevox_make_default_initialize_options(void);
	makeDefaultInitializeOptionsFunc ffi.Fun

	// VOICEVOX_C_API VoicevoxResultCode voicevox_synthesizer_new(const struct VoicevoxOnnxruntime *onnxruntime,
	//     const struct OpenJtalkRc *open_jtalk, struct VoicevoxInitializeOptions options, struct VoicevoxSynthesizer **out_synthesizer);
	synthesizerNewFunc ffi.Fun

	// VOICEVOX_C_API void voicevox_synthesizer_delete(struct VoicevoxSynthesizer *synthesizer);
	synthesizerDeleteFunc ffi.Fun

	// VOICEVOX_C_API bool voicevox_synthesizer_is_gpu_mode(const struct VoicevoxSynthesizer *synthesizer);
	synthesizerIsGPUModeFunc ffi.Fun

	// VOICEVOX_C_API struct VoicevoxLoadVoiceModelOptions voicevox_make_default_load_voice_model_options(void);
	makeDefaultLoadVoiceModelOptionsFunc ffi.Fun

	// VOICEVOX_C_API VoicevoxResultCode voicevox_synthesizer_load_voice_model(const struct VoicevoxSynthesizer *synthesizer,
	//     const struct VoicevoxVoiceModelFile *model, struct VoicevoxLoadVoiceModelOptions options);
	synthesizerLoadVoiceModelFunc ffi.Fun

	// VOICEVOX_C_API char *voicevox_synthesizer_create_metas_json(const struct VoicevoxSynthesizer *synthesizer);
	synthesizerCreateMetasJSONFunc ffi.Fun

	// VOICEVOX_C_API struct VoicevoxTtsOptions voicevox_make_default_tts_options(void);
	makeDefaultTTSOptionsFunc ffi.Fun

	// VOICEVOX_C_API VoicevoxResultCode voicevox_synthesizer_tts(const struct VoicevoxSynthesizer *synthesizer, const char *text,
	//     VoicevoxStyleId style_id, struct VoicevoxTtsOptions options, uintptr_t *output_wav_length, uint8_t **output_wav);
	synthesizerTTSFunc ffi.Fun

	// VOICEVOX_C_API VoicevoxResultCode voicevox_synthesizer_tts_from_kana(const struct VoicevoxSynthesizer *synthesizer, const char *kana,
	//     VoicevoxStyleId style_id, struct VoicevoxTtsOptions options, uintptr_t *output_wav_length, uint8_t **output_wav);
	synthesizerTTSFromKanaFunc ffi.Fun

	// VOICEVOX_C_API VoicevoxResultCode voicevox_synthesizer_create_audio_query(const struct VoicevoxSynthesizer *synthesizer,
	//     const char *text, VoicevoxStyleId style_id, char **output_audio_query_json);
	synthesizerCreateAudioQueryFunc ffi.Fun

	// VOICEVOX_C_API VoicevoxResultCode voicevox_synthesizer_create_audio_query_from_kana(const struct VoicevoxSynthesizer *synthesizer,
	//     const char *kana, VoicevoxStyleId style_id, char **output_audio_query_json);
	synthesizerCreateAudioQueryFromKanaFunc ffi.Fun

	// VOICEVOX_C_API struct VoicevoxSynthesisOptions voicevox_make_default_synthesis_options(void);
	makeDefaultSynthesisOptionsFunc ffi.Fun

	// VOICEVOX_C_API VoicevoxResultCode voicevox_synthesizer_synthesis(const struct VoicevoxSynthesizer *synthesizer,
	//     const char *audio_query_json, VoicevoxStyleId style_id, struct VoicevoxSynthesisOptions options,
	//     uintptr_t *output_wav_length, uint8_t **output_wav);
	synthesizerSynthesisFunc ffi.Fun
)

func loadSynthesizerFuncs(lib loader.Lib) error {
	var err error

	if makeDefaultInitializeOptionsFunc, err = lib.Prep("voicevox_make_default_initialize_options", &ffiTypeInitializeOptions); err != nil {
		return loadError("voicevox_make_default_initialize_options", err)
	}

	if synthesizerNewFunc, err = lib.Prep("voicevox_synthesizer_new", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffiTypeInitializeOptions, &ffi.TypePointer); err != nil {
		return loadError("voicevox_synthesizer_new", err)
	}

	if synthesizerDeleteFunc, err = lib.Prep("voicevox_synthesizer_delete", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("voicevox_synthesizer_delete", err)
	}

	if synthesizerIsGPUModeFunc, err = lib.Prep("voicevox_synthesizer_is_gpu_mode", &ffi.TypeUint8, &ffi.TypePointer); err != nil {
		return loadError("voicevox_synthesizer_is_gpu_mode", err)
	}

	if makeDefaultLoadVoiceModelOptionsFunc, err = lib.Prep("voicevox_make_default_load_voice_model_options", &ffiTypeLoadVoiceModelOptions); err != nil {
		return loadError("voicevox_make_default_load_voice_model_options", err)
	}

	if synthesizerLoadVoiceModelFunc, err = lib.Prep("voicevox_synthesizer_load_voice_model", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffiTypeLoadVoiceModelOptions); err != nil {
		return loadError("voicevox_synthesizer_load_voice_model", err)
	}

	if synthesizerCreateMetasJSONFunc, err = lib.Prep("voicevox_synthesizer_create_metas_json", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("voicevox_synthesizer_create_metas_json", err)
	}

	if makeDefaultTTSOptionsFunc, err = lib.Prep("voicevox_make_default_tts_options", &ffiTypeTTSOptions); err != nil {
		return loadError("voicevox_make_default_tts_options", err)
	}

	if synthesizerTTSFunc, err = lib.Prep("voicevox_synthesizer_tts", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32, &ffiTypeTTSOptions, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("voicevox_synthesizer_tts", err)
	}

	if synthesizerTTSFromKanaFunc, err = lib.Prep("voicevox_synthesizer_tts_from_kana", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32, &ffiTypeTTSOptions, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("voicevox_synthesizer_tts_from_kana", err)
	}

	if synthesizerCreateAudioQueryFunc, err = lib.Prep("voicevox_synthesizer_create_audio_query", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32, &ffi.TypePointer); err != nil {
		return loadError("voicevox_synthesizer_create_audio_query", err)
	}

	if synthesizerCreateAudioQueryFromKanaFunc, err = lib.Prep("voicevox_synthesizer_create_audio_query_from_kana", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32, &ffi.TypePointer); err != nil {
		return loadError("voicevox_synthesizer_create_audio_query_from_kana", err)
	}

	if makeDefaultSynthesisOptionsFunc, err = lib.Prep("voicevox_make_default_synthesis_options", &ffiTypeSynthesisOptions); err != nil {
		return loadError("voicevox_make_default_synthesis_options", err)
	}

	if synthesizerSynthesisFunc, err = lib.Prep("voicevox_synthesizer_synthesis", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32, &ffiTypeSynthesisOptions, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("voicevox_synthesizer_synthesis", err)
	}

	return nil
}

// DefaultInitializeOptions returns the default options for NewSynthesizer.
func DefaultInitializeOptions() InitializeOptions {
	var opts InitializeOptions
	makeDefaultInitializeOptionsFunc.Call(&opts)
	return opts
}

// DefaultLoadVoiceModelOptions returns the default options for LoadVoiceModel.
func DefaultLoadVoiceModelOptions() LoadVoiceModelOptions {
	var opts LoadVoiceModelOptions
	makeDefaultLoadVoiceModelOptionsFunc.Call(&opts)
	return opts
}

// DefaultTTSOptions returns the default options for TTS.
func DefaultTTSOptions() TTSOptions {
	var opts TTSOptions
	makeDefaultTTSOptionsFunc.Call(&opts)
	return opts
}

// DefaultSynthesisOptions returns the default options for Synthesis.
func DefaultSynthesisOptions() SynthesisOptions {
	var opts SynthesisOptions
	makeDefaultSynthesisOptionsFunc.Call(&opts)
	return opts
}

// NewSynthesizer creates a synthesizer. The OpenJtalk can be deleted once this returns.
func NewSynthesizer(ort Onnxruntime, ojt OpenJtalk, opts InitializeOptions) (Synthesizer, error) {
	var (
		result ffi.Arg
		syn    Synthesizer
	)
	out := unsafe.Pointer(&syn)
	synthesizerNewFunc.Call(&result, &ort, &ojt, &opts, &out)

	if err := resultError(result); err != nil {
		return 0, err
	}
	return syn, nil
}

// Delete frees the synthesizer.
func (s Synthesizer) Delete() {
	if s == 0 {
		return
	}
	synthesizerDeleteFunc.Call(nil, &s)
}

// IsGPUMode reports whether the synthesizer runs on the GPU.
func (s Synthesizer) IsGPUMode() bool {
	if s == 0 {
		return false
	}
	var result ffi.Arg
	synthesizerIsGPUModeFunc.Call(&result, &s)
	return result.Bool()
}

// LoadVoiceModel loads a voice model into the synthesizer. The model file can be deleted once this returns.
func (s Synthesizer) LoadVoiceModel(model VoiceModelFile, opts LoadVoiceModelOptions) error {
	if s == 0 {
		return errInvalidSynthesizer
	}
	if model == 0 {
		return errInvalidModel
	}

	var result ffi.Arg
	synthesizerLoadVoiceModelFunc.Call(&result, &s, &model, &opts)
	return resultError(result)
}

// MetasJSON returns the metadata of all loaded voice models as JSON.
func (s Synthesizer) MetasJSON() string {
	if s == 0 {
		return ""
	}
	var p *byte
	synthesizerCreateMetasJSONFunc.Call(unsafe.Pointer(&p), &s)
	return takeJSON(p)
}

// Metas returns the metadata of all loaded voice models.
func (s Synthesizer) Metas() ([]CharacterMeta, error) {
	if s == 0 {
		return nil, errInvalidSynthesizer
	}
	return ParseMetas(s.MetasJSON())
}

// TTS converts Japanese text to WAV audio.
func (s Synthesizer) TTS(text string, style StyleID, opts TTSOptions) ([]byte, error) {
	return s.wavCall(synthesizerTTSFunc, text, style, &opts)
}

// TTSFromKana converts AquesTalk style kana to WAV audio.
func (s Synthesizer) TTSFromKana(kana string, style StyleID, opts TTSOptions) ([]byte, error) {
	return s.wavCall(synthesizerTTSFromKanaFunc, kana, style, &opts)
}

// CreateAudioQuery creates an AudioQuery JSON for text.
func (s Synthesizer) CreateAudioQuery(text string, style StyleID) (string, error) {
	return s.jsonCall(synthesizerCreateAudioQueryFunc, text, style)
}

// CreateAudioQueryFromKana creates an AudioQuery JSON for AquesTalk style kana.
func (s Synthesizer) CreateAudioQueryFromKana(kana string, style StyleID) (string, error) {
	return s.jsonCall(synthesizerCreateAudioQueryFromKanaFunc, kana, style)
}

// Synthesis converts an AudioQuery JSON to WAV audio.
func (s Synthesizer) Synthesis(audioQuery string, style StyleID, opts SynthesisOptions) ([]byte, error) {
	return s.wavCall(synthesizerSynthesisFunc, audioQuery, style, &opts)
}

func (s Synthesizer) wavCall(fn ffi.Fun, input string, style StyleID, opts any) ([]byte, error) {
	if s == 0 {
		return nil, errInvalidSynthesizer
	}

	var (
		result ffi.Arg
		length uint64
		wav    *byte
	)
	in := cString(input)
	outLen := unsafe.Pointer(&length)
	outWav := unsafe.Pointer(&wav)
	fn.Call(&result, &s, &in, &style, opts, &outLen, &outWav)

	if err := resultError(result); err != nil {
		return nil, err
	}
	return takeWAV(wav, length), nil
}

func (s Synthesizer) jsonCall(fn ffi.Fun, input string, style StyleID) (string, error) {
	if s == 0 {
		return "", errInvalidSynthesizer
	}

	var (
		result ffi.Arg
		query  *byte
	)
	in := cString(input)
	out := unsafe.Pointer(&query)
	fn.Call(&result, &s, &in, &style, &out)

	if err := resultError(result); err != nil {
		return "", err
	}
	return takeJSON(query), nil
}
