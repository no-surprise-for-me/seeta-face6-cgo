package seeta

// #include <Seetaface6CGO.h>
import "C"
import (
	"unsafe"
)

type FaceRecognizer struct {
	ptr C.int64_t
}

func (fr FaceRecognizer) getPtr() C.int64_t {
	if fr.ptr == 0 {
		panic("FaceRecognizer对象已关闭")
	}
	return fr.ptr
}
func (fr FaceRecognizer) GetExtractFeatureSize() int {
	return int(C.GetExtractFeatureSize(fr.getPtr()))
}

// 耗时47~49ms
func (fr FaceRecognizer) Extract(img *CSeetaImageData, pfs SeetaPointFs) []float32 {
	i := int(C.GetExtractFeatureSize(fr.getPtr()))
	features := make([]float32, i)
	C.Extract(fr.ptr, C.SeetaImageData(*img), (*C.SeetaPointF)(unsafe.Pointer(&pfs.data[0])), (*C.float)(unsafe.Pointer(&features[0])))
	return features
}

func (fr *FaceRecognizer) Close() {
	if fr.ptr != 0 {
		C.DeleteFaceRecognizer(fr.ptr)
		fr.ptr = 0
	}
}

//face_recognizer_mask  口罩人脸识别   先使用MaskDetector判断是否戴口罩，如果戴了口罩，使用这个来识别人脸
//face_recognizer       通用人脸识别
func NewFaceRecognizer(modelSetting SeetaModelSetting) *FaceRecognizer {
	setting, free := modelSetting.CSeetaModelSetting()
	defer free()
	var ptr C.int64_t = C.NewFaceRecognizer(setting)
	return &FaceRecognizer{
		ptr: ptr,
	}
}
