package voicevox

import (
	"unsafe"

	"github.com/hybridgroup/voicevox/pkg/loader"
	"github.com/jupiterrider/ffi"
)

var (
	// VOICEVOX_C_API VoicevoxResultCode voicevox_open_jtalk_rc_new(const char *open_jtalk_dic_dir,
	//                                                             struct OpenJtalkRc **out_open_jtalk);
	openJtalkRcNewFunc ffi.Fun

	// VOICEVOX_C_API void voicevox_open_jtalk_rc_delete(struct OpenJtalkRc *open_jtalk);
	openJtalkRcDeleteFunc ffi.Fun
)

func loadOpenJtalkFuncs(lib loader.Lib) error {
	var err error

	if openJtalkRcNewFunc, err = lib.Prep("voicevox_open_jtalk_rc_new", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("voicevox_open_jtalk_rc_new", err)
	}

	if openJtalkRcDeleteFunc, err = lib.Prep("voicevox_open_jtalk_rc_delete", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("voicevox_open_jtalk_rc_delete", err)
	}

	return nil
}

// NewOpenJtalk creates an Open JTalk text analyzer using the dictionary in dictDir.
func NewOpenJtalk(dictDir string) (OpenJtalk, error) {
	var (
		result ffi.Arg
		ojt    OpenJtalk
	)
	dir := cString(dictDir)
	out := unsafe.Pointer(&ojt)
	openJtalkRcNewFunc.Call(&result, &dir, &out)

	if err := resultError(result); err != nil {
		return 0, err
	}
	return ojt, nil
}

// Delete frees the Open JTalk text analyzer.
func (o OpenJtalk) Delete() {
	if o == 0 {
		return
	}
	openJtalkRcDeleteFunc.Call(nil, &o)
}
