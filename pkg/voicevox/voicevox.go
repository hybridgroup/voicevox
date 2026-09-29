package voicevox

import (
	"unsafe"

	"github.com/hybridgroup/voicevox/pkg/loader"
	"github.com/jupiterrider/ffi"
	"golang.org/x/sys/unix"
)

// Onnxruntime is a handle to the loaded ONNX Runtime.
type Onnxruntime uintptr

// OpenJtalk is a handle to an Open JTalk text analyzer.
type OpenJtalk uintptr

// Synthesizer is a handle to a speech synthesizer.
type Synthesizer uintptr

// VoiceModelFile is a handle to an opened VVM file.
type VoiceModelFile uintptr

// StyleID identifies a voice style.
type StyleID uint32

// AccelerationMode selects the hardware used for inference.
type AccelerationMode int32

const (
	AccelerationModeAuto AccelerationMode = iota
	AccelerationModeCPU
	AccelerationModeGPU
)

// OnExistingVoiceModelID sets what happens when a model with the same ID is already loaded.
type OnExistingVoiceModelID int32

const (
	OnExistingVoiceModelIDError OnExistingVoiceModelID = iota
	OnExistingVoiceModelIDReload
	OnExistingVoiceModelIDSkip
)

// ResultCode is a VOICEVOX result code. Any code other than ResultOK is an error.
type ResultCode int32

const (
	ResultOK                     ResultCode = 0
	ResultNotLoadedOpenJtalkDict ResultCode = 1
	ResultGetSupportedDevices    ResultCode = 3
	ResultGPUSupport             ResultCode = 4
	ResultStyleNotFound          ResultCode = 6
	ResultModelNotFound          ResultCode = 7
	ResultRunModel               ResultCode = 8
	ResultAnalyzeText            ResultCode = 11
	ResultInvalidUTF8Input       ResultCode = 12
	ResultParseKana              ResultCode = 13
	ResultInvalidAudioQuery      ResultCode = 14
	ResultInvalidAccentPhrase    ResultCode = 15
	ResultOpenZipFile            ResultCode = 16
	ResultReadZipEntry           ResultCode = 17
	ResultModelAlreadyLoaded     ResultCode = 18
	ResultLoadUserDict           ResultCode = 20
	ResultSaveUserDict           ResultCode = 21
	ResultUserDictWordNotFound   ResultCode = 22
	ResultUseUserDict            ResultCode = 23
	ResultInvalidUserDictWord    ResultCode = 24
	ResultInvalidUUID            ResultCode = 25
	ResultStyleAlreadyLoaded     ResultCode = 26
	ResultInvalidModelData       ResultCode = 27
	ResultInvalidModelFormat     ResultCode = 28
	ResultInitInferenceRuntime   ResultCode = 29
	ResultInvalidMora            ResultCode = 30
	ResultInvalidScore           ResultCode = 31
	ResultInvalidNote            ResultCode = 32
	ResultInvalidFrameAudioQuery ResultCode = 33
	ResultInvalidFramePhoneme    ResultCode = 34
	ResultIncompatibleQueries    ResultCode = 35
)

var (
	// VOICEVOX_C_API const char *voicevox_error_result_to_message(VoicevoxResultCode result_code);
	errorResultToMessageFunc ffi.Fun

	// VOICEVOX_C_API const char *voicevox_get_version(void);
	getVersionFunc ffi.Fun
)

func loadVoicevoxFuncs(lib loader.Lib) error {
	var err error

	if errorResultToMessageFunc, err = lib.Prep("voicevox_error_result_to_message", &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("voicevox_error_result_to_message", err)
	}

	if getVersionFunc, err = lib.Prep("voicevox_get_version", &ffi.TypePointer); err != nil {
		return loadError("voicevox_get_version", err)
	}

	return nil
}

// Error returns the message for the result code.
func (r ResultCode) Error() string {
	var msg *byte
	errorResultToMessageFunc.Call(unsafe.Pointer(&msg), &r)
	return unix.BytePtrToString(msg)
}

// Version returns the voicevox_core version.
func Version() string {
	var v *byte
	getVersionFunc.Call(unsafe.Pointer(&v))
	return unix.BytePtrToString(v)
}

func resultError(result ffi.Arg) error {
	if code := ResultCode(int32(result)); code != ResultOK {
		return code
	}
	return nil
}

func cString(s string) *byte {
	return &[]byte(s + "\x00")[0]
}
