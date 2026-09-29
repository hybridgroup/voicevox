package voicevox

import (
	"encoding/json"
	"errors"
	"unsafe"

	"github.com/hybridgroup/voicevox/pkg/loader"
	"github.com/jupiterrider/ffi"
)

// CharacterMeta describes a character and its styles.
type CharacterMeta struct {
	Name        string      `json:"name"`
	Styles      []StyleMeta `json:"styles"`
	Version     string      `json:"version"`
	SpeakerUUID string      `json:"speaker_uuid"`
	Order       *int        `json:"order"`
}

// StyleMeta describes a single voice style.
type StyleMeta struct {
	ID    StyleID `json:"id"`
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Order *int    `json:"order"`
}

var errInvalidModel = errors.New("invalid voice model file")

var (
	// VOICEVOX_C_API VoicevoxResultCode voicevox_voice_model_file_open(const char *path,
	//                                                                 struct VoicevoxVoiceModelFile **out_model);
	voiceModelFileOpenFunc ffi.Fun

	// VOICEVOX_C_API char *voicevox_voice_model_file_create_metas_json(const struct VoicevoxVoiceModelFile *model);
	voiceModelFileCreateMetasJSONFunc ffi.Fun

	// VOICEVOX_C_API void voicevox_voice_model_file_delete(struct VoicevoxVoiceModelFile *model);
	voiceModelFileDeleteFunc ffi.Fun
)

func loadModelFuncs(lib loader.Lib) error {
	var err error

	if voiceModelFileOpenFunc, err = lib.Prep("voicevox_voice_model_file_open", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("voicevox_voice_model_file_open", err)
	}

	if voiceModelFileCreateMetasJSONFunc, err = lib.Prep("voicevox_voice_model_file_create_metas_json", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("voicevox_voice_model_file_create_metas_json", err)
	}

	if voiceModelFileDeleteFunc, err = lib.Prep("voicevox_voice_model_file_delete", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("voicevox_voice_model_file_delete", err)
	}

	return nil
}

// OpenVoiceModelFile opens the VVM file at path.
func OpenVoiceModelFile(path string) (VoiceModelFile, error) {
	var (
		result ffi.Arg
		model  VoiceModelFile
	)
	p := cString(path)
	out := unsafe.Pointer(&model)
	voiceModelFileOpenFunc.Call(&result, &p, &out)

	if err := resultError(result); err != nil {
		return 0, err
	}
	return model, nil
}

// MetasJSON returns the character metadata in the model as JSON.
func (m VoiceModelFile) MetasJSON() string {
	if m == 0 {
		return ""
	}
	var p *byte
	voiceModelFileCreateMetasJSONFunc.Call(unsafe.Pointer(&p), &m)
	return takeJSON(p)
}

// Metas returns the character metadata in the model.
func (m VoiceModelFile) Metas() ([]CharacterMeta, error) {
	if m == 0 {
		return nil, errInvalidModel
	}
	return ParseMetas(m.MetasJSON())
}

// Delete closes the model file.
func (m VoiceModelFile) Delete() {
	if m == 0 {
		return
	}
	voiceModelFileDeleteFunc.Call(nil, &m)
}

// ParseMetas parses character metadata JSON.
func ParseMetas(s string) ([]CharacterMeta, error) {
	var metas []CharacterMeta
	if err := json.Unmarshal([]byte(s), &metas); err != nil {
		return nil, err
	}
	return metas, nil
}
