package voicevox

import (
	"unsafe"

	"github.com/hybridgroup/voicevox/pkg/loader"
	"github.com/jupiterrider/ffi"
	"golang.org/x/sys/unix"
)

var (
	// VOICEVOX_C_API void voicevox_json_free(char *json);
	jsonFreeFunc ffi.Fun

	// VOICEVOX_C_API void voicevox_wav_free(uint8_t *wav);
	wavFreeFunc ffi.Fun
)

func loadMemoryFuncs(lib loader.Lib) error {
	var err error

	if jsonFreeFunc, err = lib.Prep("voicevox_json_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("voicevox_json_free", err)
	}

	if wavFreeFunc, err = lib.Prep("voicevox_wav_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("voicevox_wav_free", err)
	}

	return nil
}

func takeJSON(p *byte) string {
	if p == nil {
		return ""
	}
	s := unix.BytePtrToString(p)
	jsonFreeFunc.Call(nil, unsafe.Pointer(&p))
	return s
}

func takeWAV(p *byte, n uint64) []byte {
	if p == nil {
		return nil
	}
	wav := make([]byte, n)
	copy(wav, unsafe.Slice(p, n))
	wavFreeFunc.Call(nil, unsafe.Pointer(&p))
	return wav
}
