package seeta

// #include <Seetaface6CGO.h>
import "C"
import "unsafe"

type FaceLandmarker struct {
	ptr C.int64_t
}

func (fl FaceLandmarker) getPtr() C.int64_t {
	if fl.ptr == 0 {
		panic("FaceLandmarker对象已关闭")
	}
	return fl.ptr
}
func (fl FaceLandmarker) GetMarkPointNumber() int {
	return int(C.GetMarkPointNumber(fl.getPtr()))
}
func (fl FaceLandmarker) Mark(img *CSeetaImageData, rect SeetaRect) SeetaPointFs {
	size := fl.GetMarkPointNumber()
	points := make([]C.SeetaPointF, size) // 用来接收标记的人脸特征点信息
	C.mark(fl.getPtr(), C.SeetaImageData(*img), rect.CSeetaRect(), (*C.SeetaPointF)(unsafe.Pointer(&points[0])))
	return SeetaPointFs{
		size: size,
		data: points,
	}
}

// 耗时4~6ms
// 标记人脸特征点同时判断人脸是否被遮挡
// void markWithMask(int64_t ptr, SeetaImageData image, SeetaRect sr, SeetaPointF* points, int32_t* mask);
func (fl FaceLandmarker) MarkWithMask(img *CSeetaImageData, rect SeetaRect) (SeetaPointFs, QualityResult) {
	size := fl.GetMarkPointNumber()
	points := make([]C.SeetaPointF, size) // 用来接收标记的人脸特征点信息
	mask := make([]C.int32_t, size)       // 用来接收标记的人脸特征点信息
	C.markWithMask(fl.getPtr(), C.SeetaImageData(*img), rect.CSeetaRect(), (*C.SeetaPointF)(unsafe.Pointer(&points[0])), (*C.int32_t)(unsafe.Pointer(&mask[0])))
	maskCount := 0
	for _, m := range mask {
		if m != 0 { // 不等于0 true 有遮挡
			maskCount++
		}
	}
	var level QualityLevel
	if maskCount > 0 {
		level = LOW //  有任意一个被遮挡，认定图片质量为LOW
	} else {
		level = HIGH
	}
	return SeetaPointFs{
			size: size,
			data: points,
		}, QualityResult{
			Level: level,
			Score: 1 - float32(maskCount)/float32(size), //  一般是5点标记，所以这个分数只能是 0 ， 0.20  0.4  0.6  0.8 1
			Desc:  "遮挡评估",
		}
}

func (fl *FaceLandmarker) Close() {
	if fl.ptr != 0 {
		C.DeleteFaceLandmarker(fl.ptr)
		fl.ptr = 0
	}
}
func NewFaceLandmarker(modelSetting SeetaModelSetting) *FaceLandmarker {
	setting, free := modelSetting.CSeetaModelSetting()
	defer free()
	var ptr C.int64_t = C.NewFaceLandmarker(setting)
	return &FaceLandmarker{
		ptr: ptr,
	}
}
