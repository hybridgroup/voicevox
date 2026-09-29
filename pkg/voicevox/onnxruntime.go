package voicevox

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"unsafe"

	"github.com/hybridgroup/voicevox/pkg/loader"
	"github.com/jupiterrider/ffi"
	"golang.org/x/sys/unix"
)

// LoadOnnxruntimeOptions are the options for LoadOnnxruntime.
type LoadOnnxruntimeOptions struct {
	Filename *byte
}

// SupportedDevices lists the devices the loaded ONNX Runtime was built for.
// A device can be true without the hardware, per voicevox_core.h 0.17.0 voicevox_onnxruntime_create_supported_devices_json.
type SupportedDevices struct {
	CPU  bool `json:"cpu"`
	CUDA bool `json:"cuda"`
	DML  bool `json:"dml"`
}

var errInvalidOnnxruntime = errors.New("invalid onnxruntime")

var (
	ffiTypeLoadOnnxruntimeOptions = ffi.NewType(&ffi.TypePointer)

	// VOICEVOX_C_API struct VoicevoxLoadOnnxruntimeOptions voicevox_make_default_load_onnxruntime_options(void);
	makeDefaultLoadOnnxruntimeOptionsFunc ffi.Fun

	// VOICEVOX_C_API VoicevoxResultCode voicevox_onnxruntime_load_once(struct VoicevoxLoadOnnxruntimeOptions options,
	//                                                                 const struct VoicevoxOnnxruntime **out_onnxruntime);
	onnxruntimeLoadOnceFunc ffi.Fun

	// VOICEVOX_C_API const char *voicevox_get_onnxruntime_lib_recommended_versioned_filename(void);
	getOnnxruntimeLibVersionedFilenameFunc ffi.Fun

	// VOICEVOX_C_API const char *voicevox_get_onnxruntime_lib_recommended_unversioned_filename(void);
	getOnnxruntimeLibUnversionedFilenameFunc ffi.Fun

	// VOICEVOX_C_API VoicevoxResultCode voicevox_onnxruntime_create_supported_devices_json(const struct VoicevoxOnnxruntime *onnxruntime,
	//                                                                                     char **output_supported_devices_json);
	onnxruntimeCreateSupportedDevicesJSONFunc ffi.Fun
)

func loadOnnxruntimeFuncs(lib loader.Lib) error {
	var err error

	if makeDefaultLoadOnnxruntimeOptionsFunc, err = lib.Prep("voicevox_make_default_load_onnxruntime_options", &ffiTypeLoadOnnxruntimeOptions); err != nil {
		return loadError("voicevox_make_default_load_onnxruntime_options", err)
	}

	if onnxruntimeLoadOnceFunc, err = lib.Prep("voicevox_onnxruntime_load_once", &ffi.TypeSint32, &ffiTypeLoadOnnxruntimeOptions, &ffi.TypePointer); err != nil {
		return loadError("voicevox_onnxruntime_load_once", err)
	}

	if getOnnxruntimeLibVersionedFilenameFunc, err = lib.Prep("voicevox_get_onnxruntime_lib_recommended_versioned_filename", &ffi.TypePointer); err != nil {
		return loadError("voicevox_get_onnxruntime_lib_recommended_versioned_filename", err)
	}

	if getOnnxruntimeLibUnversionedFilenameFunc, err = lib.Prep("voicevox_get_onnxruntime_lib_recommended_unversioned_filename", &ffi.TypePointer); err != nil {
		return loadError("voicevox_get_onnxruntime_lib_recommended_unversioned_filename", err)
	}

	if onnxruntimeCreateSupportedDevicesJSONFunc, err = lib.Prep("voicevox_onnxruntime_create_supported_devices_json", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("voicevox_onnxruntime_create_supported_devices_json", err)
	}

	return nil
}

// DefaultLoadOnnxruntimeOptions returns the default options for LoadOnnxruntime.
func DefaultLoadOnnxruntimeOptions() LoadOnnxruntimeOptions {
	var opts LoadOnnxruntimeOptions
	makeDefaultLoadOnnxruntimeOptionsFunc.Call(&opts)
	return opts
}

// OnnxruntimeLibFilename returns the recommended versioned ONNX Runtime library filename.
func OnnxruntimeLibFilename() string {
	var name *byte
	getOnnxruntimeLibVersionedFilenameFunc.Call(unsafe.Pointer(&name))
	return unix.BytePtrToString(name)
}

// OnnxruntimeLibUnversionedFilename returns the recommended ONNX Runtime library filename without a version.
func OnnxruntimeLibUnversionedFilename() string {
	var name *byte
	getOnnxruntimeLibUnversionedFilenameFunc.Call(unsafe.Pointer(&name))
	return unix.BytePtrToString(name)
}

// LoadOnnxruntime loads the ONNX Runtime library from dir. It only loads once per process.
// If the recommended version is not in dir, it uses the newest other version found there.
func LoadOnnxruntime(dir string) (Onnxruntime, error) {
	opts := DefaultLoadOnnxruntimeOptions()
	if dir != "" {
		opts.Filename = cString(findOnnxruntimeLib(dir))
	}

	return LoadOnnxruntimeWithOptions(opts)
}

func findOnnxruntimeLib(dir string) string {
	recommended := filepath.Join(dir, OnnxruntimeLibFilename())
	if _, err := os.Stat(recommended); err == nil {
		return recommended
	}

	name := OnnxruntimeLibUnversionedFilename()
	matches, _ := filepath.Glob(filepath.Join(dir, strings.TrimSuffix(name, filepath.Ext(name))+".*"))
	if len(matches) == 0 {
		return recommended
	}

	sort.Strings(matches)
	return matches[len(matches)-1]
}

// LoadOnnxruntimeWithOptions loads the ONNX Runtime library using opts.
func LoadOnnxruntimeWithOptions(opts LoadOnnxruntimeOptions) (Onnxruntime, error) {
	var (
		result ffi.Arg
		ort    Onnxruntime
	)
	out := unsafe.Pointer(&ort)
	onnxruntimeLoadOnceFunc.Call(&result, &opts, &out)
	runtime.KeepAlive(opts.Filename)

	if err := resultError(result); err != nil {
		return 0, err
	}
	return ort, nil
}

// SupportedDevicesJSON returns the devices supported by the ONNX Runtime as JSON.
func (o Onnxruntime) SupportedDevicesJSON() (string, error) {
	if o == 0 {
		return "", errInvalidOnnxruntime
	}

	var (
		result ffi.Arg
		p      *byte
	)
	out := unsafe.Pointer(&p)
	onnxruntimeCreateSupportedDevicesJSONFunc.Call(&result, &o, &out)

	if err := resultError(result); err != nil {
		return "", err
	}
	return takeJSON(p), nil
}

// SupportedDevices returns the devices supported by the ONNX Runtime.
func (o Onnxruntime) SupportedDevices() (SupportedDevices, error) {
	var d SupportedDevices
	s, err := o.SupportedDevicesJSON()
	if err != nil {
		return d, err
	}
	err = json.Unmarshal([]byte(s), &d)
	return d, err
}
