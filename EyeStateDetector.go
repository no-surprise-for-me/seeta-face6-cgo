package seeta

// #include <Seetaface6CGO.h>
import "C"
import "unsafe"

type EyeState int

const (
	EYE_CLOSE = iota
	EYE_OPEN
	EYE_RANDOM
	EYE_UNKNOWN
)

type EyeStates struct {
	Left  EyeState
	Right EyeState
}
type EyeStateDetector struct {
	ptr C.int64_t
}

func (esd EyeStateDetector) getPtr() C.int64_t {
	if esd.ptr == 0 {
		panic("EyeStateDetector对象已关闭")
	}
	return esd.ptr
}

//void DeleteEyeStateDetector(int64_t ptr)
func (esd *EyeStateDetector) Close() {
	if esd.ptr != 0 {
		C.DeleteEyeStateDetector(esd.ptr)
		esd.ptr = 0
	}
}

//int64_t NewEyeStateDetector(SeetaModelSetting model)
func NewEyeStateDetector(modelSetting SeetaModelSetting) *EyeStateDetector {
	setting, free := modelSetting.CSeetaModelSetting()
	defer free()
	var ptr C.int64_t = C.NewEyeStateDetector(setting)
	return &EyeStateDetector{
		ptr: ptr,
	}
}

//EyeStates DetectEyeState(int64_t ptr, SeetaImageData image, SeetaPointF* points)
func (esd EyeStateDetector) DetectEyeState(img *CSeetaImageData, points SeetaPointFs) EyeStates {
	var state C.EyeStates = C.DetectEyeState(esd.getPtr(), C.SeetaImageData(*img), (*C.SeetaPointF)(unsafe.Pointer(&points.data[0])))
	return EyeStates{
		Left:  EyeState(state.left),
		Right: EyeState(state.right),
	}
}
